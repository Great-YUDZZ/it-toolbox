package youtube

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/kkdai/youtube/v2"
)

const (
	// AndroidClientUserAgent bypasses YouTube 403 Forbidden throttling on direct streams.
	AndroidClientUserAgent = "com.google.android.youtube/20.10.38 (Linux; U; Android 11) gzip"
)

// ResolutionOption represents an available resolution or audio track for a video.
type ResolutionOption struct {
	Itag          int
	QualityLabel  string // e.g. "1080p60", "720p", "480p", "360p", "144p", "Audio M4A"
	Resolution    int    // numeric height: 2160, 1080, 720, etc. (0 for audio)
	FPS           int
	IsAudio       bool
	IsMuxed       bool   // contains both video and audio
	MimeType      string
	Container     string // "mp4", "m4a", "webm"
	ContentLength int64
	SizeMB        float64
	DisplayLabel  string // string displayed in UI select dropdown
	format        *youtube.Format
}

// VideoDetails holds metadata and available resolutions for a YouTube video.
type VideoDetails struct {
	ID           string
	Title        string
	Author       string
	Duration     time.Duration
	DurationStr  string
	ThumbnailURL string
	Options      []ResolutionOption
	BestAudio    *ResolutionOption
	RawVideo     *youtube.Video
}

// ProgressCallback reports live download progress.
type ProgressCallback func(downloaded int64, total int64, percent float64, speedMBps float64, statusText string)

// Downloader manages fetching video details and downloading streams.
type Downloader struct {
	client *youtube.Client
}

// NewDownloader returns a new Downloader instance.
func NewDownloader() *Downloader {
	return &Downloader{
		client: &youtube.Client{},
	}
}

// HasFFmpeg checks if ffmpeg is installed and accessible.
func HasFFmpeg() bool {
	_, err := exec.LookPath("ffmpeg")
	return err == nil
}

// ExtractVideoID parses YouTube video ID from various URL formats or raw IDs.
func ExtractVideoID(input string) string {
	input = strings.TrimSpace(input)
	if input == "" {
		return ""
	}
	if strings.Contains(input, "://") {
		if !strings.Contains(input, "youtube.com") && !strings.Contains(input, "youtu.be") {
			return ""
		}
	} else if strings.Contains(input, "/") || strings.Contains(input, ".") {
		if !strings.Contains(input, "youtube.com") && !strings.Contains(input, "youtu.be") {
			return ""
		}
	} else {
		if len(input) == 11 {
			return input
		}
		return ""
	}

	if id, err := youtube.ExtractVideoID(input); err == nil && id != "" {
		return id
	}
	// Fallback regex for shorts or unusual URLs
	re := regexp.MustCompile(`(?:v=|/v/|youtu\.be/|/embed/|/shorts/)([a-zA-Z0-9_-]{11})`)
	matches := re.FindStringSubmatch(input)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

// SanitizeFilename cleans the title into a safe filename.
func SanitizeFilename(name string) string {
	// Replace illegal chars along with surrounding spaces with single underscore
	reg := regexp.MustCompile(`\s*[<>:"/\\|?*\x00-\x1F]+\s*`)
	clean := reg.ReplaceAllString(name, "_")
	clean = regexp.MustCompile(`_+`).ReplaceAllString(clean, "_")
	clean = regexp.MustCompile(`\s+`).ReplaceAllString(clean, " ")
	clean = strings.TrimSpace(clean)
	if len(clean) > 90 {
		clean = clean[:90]
	}
	clean = strings.Trim(clean, " ._")
	if clean == "" {
		clean = "youtube_video"
	}
	return clean
}

// FormatDuration formats duration into mm:ss or hh:mm:ss.
func FormatDuration(d time.Duration) string {
	totalSec := int(d.Seconds())
	h := totalSec / 3600
	m := (totalSec % 3600) / 60
	s := totalSec % 60
	if h > 0 {
		return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
	}
	return fmt.Sprintf("%02d:%02d", m, s)
}

// parseResolutionInt extracts numeric height from quality label like "1080p60", "720p".
func parseResolutionInt(label string) int {
	re := regexp.MustCompile(`(\d+)p`)
	matches := re.FindStringSubmatch(label)
	if len(matches) > 1 {
		val, _ := strconv.Atoi(matches[1])
		return val
	}
	return 0
}

// FetchVideoDetails fetches video metadata and organizes all distinct resolution options.
func (d *Downloader) FetchVideoDetails(rawURL string) (*VideoDetails, error) {
	vidID := ExtractVideoID(rawURL)
	if vidID == "" {
		return nil, fmt.Errorf("URL atau ID YouTube tidak valid: %s", rawURL)
	}

	video, err := d.client.GetVideo(vidID)
	if err != nil {
		return nil, fmt.Errorf("gagal memuat data video: %w", err)
	}

	details := &VideoDetails{
		ID:          video.ID,
		Title:       video.Title,
		Author:      video.Author,
		Duration:    video.Duration,
		DurationStr: FormatDuration(video.Duration),
		RawVideo:    video,
	}

	if len(video.Thumbnails) > 0 {
		details.ThumbnailURL = video.Thumbnails[len(video.Thumbnails)-1].URL
	}

	// 1. Separate video formats and audio formats
	type bestResEntry struct {
		format *youtube.Format
		option ResolutionOption
	}
	resMap := make(map[int]bestResEntry) // mapped by height (e.g. 1080, 720)
	var audioOptions []ResolutionOption
	var maxAudioBitrate int

	for i := range video.Formats {
		f := &video.Formats[i]
		isAudio := strings.Contains(f.MimeType, "audio")
		isVideo := strings.Contains(f.MimeType, "video")

		if isAudio {
			container := "m4a"
			if strings.Contains(f.MimeType, "webm") {
				container = "webm"
			}
			sizeMB := float64(f.ContentLength) / (1024 * 1024)

			tier := "Sedang"
			if strings.Contains(f.AudioQuality, "HIGH") || f.AverageBitrate >= 200000 {
				tier = "Tinggi"
			} else if strings.Contains(f.AudioQuality, "LOW") || (f.AverageBitrate > 0 && f.AverageBitrate < 96000) {
				tier = "Hemat"
			}

			bitrateKbps := f.AverageBitrate / 1000
			if bitrateKbps == 0 {
				bitrateKbps = f.Bitrate / 1000
			}

			label := fmt.Sprintf("Audio %s (%s", strings.ToUpper(container), tier)
			if bitrateKbps > 0 {
				label = fmt.Sprintf("%s, ~%d kbps)", label, bitrateKbps)
			} else {
				label = label + ")"
			}

			disp := fmt.Sprintf("%s — %.2f MB (Audio Only)", label, sizeMB)
			opt := ResolutionOption{
				Itag:          f.ItagNo,
				QualityLabel:  label,
				Resolution:    0,
				IsAudio:       true,
				IsMuxed:       false,
				MimeType:      f.MimeType,
				Container:     container,
				ContentLength: f.ContentLength,
				SizeMB:        sizeMB,
				DisplayLabel:  disp,
				format:        f,
			}
			audioOptions = append(audioOptions, opt)

			if strings.Contains(f.MimeType, "audio/mp4") && f.AverageBitrate > maxAudioBitrate {
				maxAudioBitrate = f.AverageBitrate
				bestOpt := opt
				details.BestAudio = &bestOpt
			}
			continue
		}

		if isVideo {
			qLabel := f.QualityLabel
			if qLabel == "" {
				qLabel = f.Quality
			}
			height := parseResolutionInt(qLabel)
			if height == 0 && f.Height > 0 {
				height = f.Height
				qLabel = fmt.Sprintf("%dp", height)
			}
			if height == 0 {
				continue
			}

			isMuxed := f.AudioChannels > 0
			container := "mp4"
			if strings.Contains(f.MimeType, "webm") {
				container = "webm"
			}

			sizeMB := float64(f.ContentLength) / (1024 * 1024)
			opt := ResolutionOption{
				Itag:          f.ItagNo,
				QualityLabel:  qLabel,
				Resolution:    height,
				FPS:           f.FPS,
				IsAudio:       false,
				IsMuxed:       isMuxed,
				MimeType:      f.MimeType,
				Container:     container,
				ContentLength: f.ContentLength,
				SizeMB:        sizeMB,
				format:        f,
			}

			// Deduplication: choose best format for this resolution
			// Preference rules:
			// 1. MP4 over WebM
			// 2. Muxed format if available (e.g. itag 18)
			// 3. Higher FPS (60fps over 30fps)
			// 4. Higher content length / bitrate
			existing, exists := resMap[height]
			if !exists {
				resMap[height] = bestResEntry{format: f, option: opt}
			} else {
				replace := false
				if !existing.option.IsMuxed && isMuxed {
					replace = true
				} else if existing.option.Container != "mp4" && container == "mp4" {
					replace = true
				} else if opt.FPS > existing.option.FPS {
					replace = true
				} else if opt.ContentLength > existing.option.ContentLength && container == existing.option.Container {
					replace = true
				}
				if replace {
					resMap[height] = bestResEntry{format: f, option: opt}
				}
			}
		}
	}

	// If no M4A audio was marked best, fallback to any audio format
	if details.BestAudio == nil && len(audioOptions) > 0 {
		details.BestAudio = &audioOptions[0]
	}

	// 2. Build sorted video options (highest resolution first)
	var heights []int
	for h := range resMap {
		heights = append(heights, h)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(heights)))

	ffmpegInstalled := HasFFmpeg()
	for _, h := range heights {
		entry := resMap[h]
		opt := entry.option

		var extraInfo string
		if opt.IsMuxed {
			extraInfo = "Video + Audio"
		} else if ffmpegInstalled && details.BestAudio != nil {
			extraInfo = "Video + Audio (Auto-Mux)"
		} else {
			extraInfo = "Video Saja"
		}

		sizeStr := fmt.Sprintf("%.2f MB", opt.SizeMB)
		if opt.SizeMB <= 0.01 {
			sizeStr = "Otomatis"
		}

		opt.DisplayLabel = fmt.Sprintf("%s (%s) — %s [%s]",
			opt.QualityLabel, strings.ToUpper(opt.Container), sizeStr, extraInfo)

		details.Options = append(details.Options, opt)
	}

	// 3. Append distinct audio options (prefer M4A, unique tier)
	seenAudioTier := make(map[string]bool)
	for _, a := range audioOptions {
		if a.Container == "m4a" {
			// Extract tier keyword e.g. "Tinggi", "Sedang", "Hemat"
			tierKey := "Sedang"
			if strings.Contains(a.QualityLabel, "Tinggi") {
				tierKey = "Tinggi"
			} else if strings.Contains(a.QualityLabel, "Hemat") {
				tierKey = "Hemat"
			}
			fullKey := "m4a_" + tierKey
			if !seenAudioTier[fullKey] {
				seenAudioTier[fullKey] = true
				details.Options = append(details.Options, a)
			}
		}
	}

	if len(details.Options) == 0 {
		return nil, fmt.Errorf("tidak ada format video/audio yang dapat diunduh untuk video ini")
	}

	return details, nil
}

// Download stream to target destination with live progress callback.
func (d *Downloader) Download(ctx context.Context, details *VideoDetails, opt ResolutionOption, destDir string, onProgress ProgressCallback) (string, error) {
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return "", fmt.Errorf("gagal membuat direktori tujuan: %w", err)
	}

	cleanTitle := SanitizeFilename(details.Title)
	var finalFileName string
	if opt.IsAudio {
		finalFileName = fmt.Sprintf("%s_[%s].%s", cleanTitle, SanitizeFilename(opt.QualityLabel), opt.Container)
	} else {
		finalFileName = fmt.Sprintf("%s_[%s].%s", cleanTitle, opt.QualityLabel, opt.Container)
	}
	finalPath := filepath.Join(destDir, finalFileName)

	// Determine whether we can / should mux video + audio
	needMux := !opt.IsAudio && !opt.IsMuxed && HasFFmpeg() && details.BestAudio != nil

	if !needMux {
		// Single stream download
		err := d.downloadSingleStream(ctx, details.RawVideo, opt.format, finalPath, 0.0, 100.0, onProgress)
		if err != nil {
			return "", err
		}
		return finalPath, nil
	}

	// Two-pass download + FFmpeg muxing
	tempVideoPath := filepath.Join(destDir, fmt.Sprintf(".tmp_%d_video.%s", time.Now().UnixNano(), opt.Container))
	tempAudioPath := filepath.Join(destDir, fmt.Sprintf(".tmp_%d_audio.%s", time.Now().UnixNano(), details.BestAudio.Container))
	defer func() {
		os.Remove(tempVideoPath)
		os.Remove(tempAudioPath)
	}()

	// Pass 1: Download Video (0% - 70%)
	if onProgress != nil {
		onProgress(0, opt.ContentLength, 0, 0, fmt.Sprintf("Mengunduh video %s...", opt.QualityLabel))
	}
	if err := d.downloadSingleStream(ctx, details.RawVideo, opt.format, tempVideoPath, 0.0, 70.0, onProgress); err != nil {
		return "", fmt.Errorf("gagal mengunduh stream video: %w", err)
	}

	// Pass 2: Download Audio (70% - 90%)
	if onProgress != nil {
		onProgress(0, details.BestAudio.ContentLength, 70.0, 0, "Mengunduh audio berkualitas tinggi...")
	}
	if err := d.downloadSingleStream(ctx, details.RawVideo, details.BestAudio.format, tempAudioPath, 70.0, 90.0, onProgress); err != nil {
		return "", fmt.Errorf("gagal mengunduh stream audio: %w", err)
	}

	// Pass 3: Mux with FFmpeg (90% - 100%)
	if onProgress != nil {
		onProgress(0, 0, 92.0, 0, "Menggabungkan video dan audio (FFmpeg Muxing)...")
	}

	cmd := exec.CommandContext(ctx, "ffmpeg", "-y",
		"-i", tempVideoPath,
		"-i", tempAudioPath,
		"-c:v", "copy",
		"-c:a", "aac",
		"-movflags", "+faststart",
		finalPath,
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("proses FFmpeg gagal: %v (%s)", err, strings.TrimSpace(string(output)))
	}

	if onProgress != nil {
		onProgress(opt.ContentLength, opt.ContentLength, 100.0, 0, "Unduhan & penggabungan selesai!")
	}

	return finalPath, nil
}

// downloadSingleStream handles downloading a single format stream to a file.
func (d *Downloader) downloadSingleStream(ctx context.Context, video *youtube.Video, format *youtube.Format, targetPath string, startPct, endPct float64, onProgress ProgressCallback) error {
	streamURL, err := d.client.GetStreamURL(video, format)
	if err != nil {
		return fmt.Errorf("gagal mendapatkan URL stream: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", streamURL, nil)
	if err != nil {
		return fmt.Errorf("gagal membuat request: %w", err)
	}
	req.Header.Set("User-Agent", AndroidClientUserAgent)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("gagal menghubungkan ke YouTube stream: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusPartialContent {
		return fmt.Errorf("YouTube stream mengembalikan status HTTP %d", resp.StatusCode)
	}

	totalBytes := resp.ContentLength
	if totalBytes <= 0 {
		totalBytes = format.ContentLength
	}

	partPath := targetPath + ".part"
	outFile, err := os.Create(partPath)
	if err != nil {
		return fmt.Errorf("gagal membuat berkas sementara: %w", err)
	}

	defer func() {
		outFile.Close()
		// If partPath still exists on failure, remove it
		if _, err := os.Stat(partPath); err == nil {
			os.Remove(partPath)
		}
	}()

	buf := make([]byte, 64*1024)
	var downloaded int64
	startTime := time.Now()
	lastUpdate := time.Now()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		n, rErr := resp.Body.Read(buf)
		if n > 0 {
			_, wErr := outFile.Write(buf[:n])
			if wErr != nil {
				return fmt.Errorf("gagal menulis ke berkas: %w", wErr)
			}
			downloaded += int64(n)

			now := time.Now()
			if now.Sub(lastUpdate) >= 120*time.Millisecond || rErr != nil {
				lastUpdate = now
				elapsedSec := now.Sub(startTime).Seconds()
				var speedMBps float64
				if elapsedSec > 0 {
					speedMBps = (float64(downloaded) / (1024 * 1024)) / elapsedSec
				}

				var fraction float64
				if totalBytes > 0 {
					fraction = float64(downloaded) / float64(totalBytes)
					if fraction > 1.0 {
						fraction = 1.0
					}
				}
				currentPct := startPct + fraction*(endPct-startPct)

				if onProgress != nil {
					status := fmt.Sprintf("%.1f / %.1f MB (%.1f MB/s)",
						float64(downloaded)/(1024*1024), float64(totalBytes)/(1024*1024), speedMBps)
					onProgress(downloaded, totalBytes, currentPct, speedMBps, status)
				}
			}
		}

		if rErr != nil {
			if rErr == io.EOF {
				break
			}
			return fmt.Errorf("koneksi terputus saat mengunduh: %w", rErr)
		}
	}

	outFile.Close()
	if err := os.Rename(partPath, targetPath); err != nil {
		return fmt.Errorf("gagal menyimpan berkas akhir: %w", err)
	}

	return nil
}

package qrcode

import (
	"context"
	"fmt"
	"html/template"
	"mime"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// FileShareServer provides a lightweight HTTP server for instant LAN/Wi-Fi file sharing via QR code
type FileShareServer struct {
	mu            sync.RWMutex
	filePath      string
	fileName      string
	fileSize      int64
	mimeType      string
	port          int
	ipAddress     string
	shareURL      string
	server        *http.Server
	listener      net.Listener
	running       bool
	downloadCount int64
}

// NewFileShareServer creates an inactive file share server instance
func NewFileShareServer() *FileShareServer {
	return &FileShareServer{}
}

// FormatFileSize returns human-readable file size string (B, KB, MB, GB)
func FormatFileSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// GetLocalIPv4Addresses retrieves all non-loopback active IPv4 addresses from network interfaces
func GetLocalIPv4Addresses() ([]string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	var ips []string
	for _, iface := range ifaces {
		// Ignore interfaces that are down or loopback
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip == nil || ip.IsLoopback() {
				continue
			}

			ip = ip.To4()
			if ip == nil {
				continue // skip IPv6 for simple local Wi-Fi pairing
			}

			ips = append(ips, ip.String())
		}
	}

	if len(ips) == 0 {
		ips = append(ips, "127.0.0.1")
	}
	return ips, nil
}

const mobileLandingTemplate = `<!DOCTYPE html>
<html lang="id">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.FileName}} - IT Toolbox File Share</title>
    <style>
        :root {
            --bg: #09090b;
            --card-bg: #18181b;
            --border: #27272a;
            --accent: #0284c7;
            --text-main: #f4f4f5;
            --text-muted: #a1a1aa;
            --btn-hover: #0369a1;
        }
        * { box-sizing: border-box; margin: 0; padding: 0; }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif;
            background-color: var(--bg);
            color: var(--text-main);
            min-height: 100vh;
            display: flex;
            align-items: center;
            justify-content: center;
            padding: 16px;
        }
        .container {
            width: 100%;
            max-width: 480px;
            background: var(--card-bg);
            border: 2px solid var(--border);
            border-radius: 12px;
            box-shadow: 0 10px 25px rgba(0, 0, 0, 0.5);
            overflow: hidden;
        }
        .header {
            padding: 20px;
            border-bottom: 2px solid var(--border);
            text-align: center;
            background: rgba(2, 132, 199, 0.08);
        }
        .brand {
            font-size: 13px;
            font-weight: 700;
            letter-spacing: 1.5px;
            color: var(--accent);
            text-transform: uppercase;
            margin-bottom: 4px;
        }
        .title {
            font-size: 18px;
            font-weight: 700;
            color: var(--text-main);
            word-break: break-all;
        }
        .preview-box {
            padding: 16px;
            text-align: center;
            background: #000;
            border-bottom: 1px solid var(--border);
        }
        .preview-img {
            max-width: 100%;
            max-height: 280px;
            border-radius: 6px;
            object-fit: contain;
        }
        .preview-video {
            width: 100%;
            max-height: 280px;
            border-radius: 6px;
            background: #000;
        }
        .content {
            padding: 20px;
        }
        .meta-list {
            margin-bottom: 24px;
            display: grid;
            grid-template-columns: 1fr 1fr;
            gap: 12px;
        }
        .meta-card {
            background: #121215;
            border: 1px solid var(--border);
            border-radius: 8px;
            padding: 10px 14px;
        }
        .meta-label {
            font-size: 11px;
            color: var(--text-muted);
            text-transform: uppercase;
            margin-bottom: 2px;
        }
        .meta-val {
            font-size: 14px;
            font-weight: 600;
            color: var(--text-main);
        }
        .actions {
            display: flex;
            flex-direction: column;
            gap: 10px;
        }
        .btn {
            display: block;
            text-align: center;
            padding: 14px 20px;
            font-size: 15px;
            font-weight: 700;
            text-decoration: none;
            border-radius: 8px;
            transition: all 0.15s ease;
        }
        .btn-primary {
            background: var(--accent);
            color: #fff;
        }
        .btn-primary:hover {
            background: var(--btn-hover);
        }
        .btn-secondary {
            background: #27272a;
            color: var(--text-main);
            border: 1px solid #3f3f46;
        }
        .btn-secondary:hover {
            background: #3f3f46;
        }
        .footer {
            padding: 14px;
            text-align: center;
            font-size: 12px;
            color: var(--text-muted);
            border-top: 1px solid var(--border);
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <div class="brand">IT Toolbox Wi-Fi File Share</div>
            <div class="title">{{.FileName}}</div>
        </div>

        {{if .IsImage}}
        <div class="preview-box">
            <img class="preview-img" src="/view" alt="{{.FileName}}">
        </div>
        {{else if .IsVideo}}
        <div class="preview-box">
            <video class="preview-video" controls preload="metadata">
                <source src="/view" type="{{.MimeType}}">
                Browser Anda tidak mendukung pemutar video HTML5.
            </video>
        </div>
        {{end}}

        <div class="content">
            <div class="meta-list">
                <div class="meta-card">
                    <div class="meta-label">Ukuran Berkas</div>
                    <div class="meta-val">{{.FormattedSize}}</div>
                </div>
                <div class="meta-card">
                    <div class="meta-label">Tipe MIME</div>
                    <div class="meta-val">{{.ShortMime}}</div>
                </div>
            </div>

            <div class="actions">
                <a href="/download" class="btn btn-primary">Unduh Berkas ({{.FormattedSize}})</a>
                <a href="/view" target="_blank" class="btn btn-secondary">Buka / Stream di Browser</a>
            </div>
        </div>

        <div class="footer">
            Transfer langsung cepat & aman via Wi-Fi lokal
        </div>
    </div>
</body>
</html>`

// Start activates the HTTP file server on the specified port and interface
func (s *FileShareServer) Start(filePath string, port int, targetIP string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.running {
		_ = s.stopLocked()
	}

	info, err := os.Stat(filePath)
	if err != nil {
		return "", fmt.Errorf("berkas tidak ditemukan atau tidak dapat dibaca: %w", err)
	}
	if info.IsDir() {
		return "", fmt.Errorf("jalur yang dipilih adalah direktori, silakan pilih berkas tunggal atau arsipkan (.zip)")
	}

	s.filePath = filePath
	s.fileName = filepath.Base(filePath)
	s.fileSize = info.Size()

	// Detect MIME type
	ext := strings.ToLower(filepath.Ext(filePath))
	s.mimeType = mime.TypeByExtension(ext)
	if s.mimeType == "" {
		s.mimeType = "application/octet-stream"
	}

	// Pick default port
	if port <= 0 {
		port = 8765
	}

	// Find available listener port
	var ln net.Listener
	var actualPort int
	for tryPort := port; tryPort < port+20; tryPort++ {
		l, err := net.Listen("tcp", fmt.Sprintf(":%d", tryPort))
		if err == nil {
			ln = l
			actualPort = tryPort
			break
		}
	}
	if ln == nil {
		return "", fmt.Errorf("gagal mengikat port jaringan (port %d-%d sedang digunakan aplikasi lain)", port, port+20)
	}

	// Determine IP
	if targetIP == "" || targetIP == "0.0.0.0" {
		ips, err := GetLocalIPv4Addresses()
		if err == nil && len(ips) > 0 {
			targetIP = ips[0]
		} else {
			targetIP = "127.0.0.1"
		}
	}

	s.port = actualPort
	s.ipAddress = targetIP
	s.listener = ln
	s.shareURL = fmt.Sprintf("http://%s:%d/", targetIP, actualPort)
	s.downloadCount = 0

	mux := http.NewServeMux()

	// Landing Page
	tmpl, err := template.New("landing").Parse(mobileLandingTemplate)
	if err != nil {
		ln.Close()
		return "", fmt.Errorf("gagal mengompilasi template halaman: %w", err)
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		s.mu.RLock()
		name := s.fileName
		sz := s.fileSize
		mType := s.mimeType
		s.mu.RUnlock()

		isImg := strings.HasPrefix(mType, "image/")
		isVid := strings.HasPrefix(mType, "video/")
		shortM := mType
		if parts := strings.Split(mType, ";"); len(parts) > 0 {
			shortM = parts[0]
		}

		data := struct {
			FileName      string
			FormattedSize string
			MimeType      string
			ShortMime     string
			IsImage       bool
			IsVideo       bool
		}{
			FileName:      name,
			FormattedSize: FormatFileSize(sz),
			MimeType:      mType,
			ShortMime:     shortM,
			IsImage:       isImg,
			IsVideo:       isVid,
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_ = tmpl.Execute(w, data)
	})

	// Direct Download Handler (Attachment)
	mux.HandleFunc("/download", func(w http.ResponseWriter, r *http.Request) {
		s.mu.RLock()
		fPath := s.filePath
		fName := s.fileName
		s.mu.RUnlock()

		file, err := os.Open(fPath)
		if err != nil {
			http.Error(w, "Gagal membuka berkas", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		fInfo, err := file.Stat()
		if err != nil {
			http.Error(w, "Gagal membaca metadata berkas", http.StatusInternalServerError)
			return
		}

		atomic.AddInt64(&s.downloadCount, 1)

		w.Header().Set("Content-Disposition", "attachment; filename=" + strconv.Quote(fName))
		w.Header().Set("Content-Type", "application/octet-stream")
		http.ServeContent(w, r, fName, fInfo.ModTime(), file)
	})

	// Inline View/Stream Handler (Partial Range Support)
	mux.HandleFunc("/view", func(w http.ResponseWriter, r *http.Request) {
		s.mu.RLock()
		fPath := s.filePath
		fName := s.fileName
		mType := s.mimeType
		s.mu.RUnlock()

		file, err := os.Open(fPath)
		if err != nil {
			http.Error(w, "Gagal membuka berkas", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		fInfo, err := file.Stat()
		if err != nil {
			http.Error(w, "Gagal membaca metadata berkas", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Disposition", "inline; filename=" + strconv.Quote(fName))
		w.Header().Set("Content-Type", mType)
		http.ServeContent(w, r, fName, fInfo.ModTime(), file)
	})

	s.server = &http.Server{
		Handler:      mux,
		ReadTimeout:  30 * time.Minute,
		WriteTimeout: 30 * time.Minute,
	}

	s.running = true

	go func() {
		_ = s.server.Serve(s.listener)
	}()

	return s.shareURL, nil
}

func (s *FileShareServer) stopLocked() error {
	if !s.running {
		return nil
	}
	s.running = false
	if s.server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_ = s.server.Shutdown(ctx)
	}
	if s.listener != nil {
		_ = s.listener.Close()
	}
	return nil
}

// Stop cleanly terminates the HTTP server and releases network resources
func (s *FileShareServer) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.stopLocked()
}

// IsRunning reports whether the server is currently hosting a file
func (s *FileShareServer) IsRunning() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.running
}

// GetShareURL returns the current active sharing URL
func (s *FileShareServer) GetShareURL() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.shareURL
}

// GetFileInfo returns metadata of the currently shared file
func (s *FileShareServer) GetFileInfo() (name string, size int64, mime string, count int64) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.fileName, s.fileSize, s.mimeType, atomic.LoadInt64(&s.downloadCount)
}

package updater

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	// DefaultGitHubRepo is the upstream repository
	DefaultGitHubRepo = "Great-YUDZZ/it-toolbox"
	// UserAgent used for GitHub API requests
	UserAgent = "IT-Toolbox-App"
	// RequestTimeout is the HTTP timeout for checking updates
	RequestTimeout = 5 * time.Second
)

// ReleaseAsset represents a downloadable file from a release
type ReleaseAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
	Size               int64  `json:"size"`
}

// GitHubRelease represents the release object returned by GitHub API
type GitHubRelease struct {
	TagName     string         `json:"tag_name"`
	Name        string         `json:"name"`
	Body        string         `json:"body"`
	HTMLURL     string         `json:"html_url"`
	PublishedAt string         `json:"published_at"`
	Assets      []ReleaseAsset `json:"assets"`
}

// UpdateCheckResult holds the result of an update check
type UpdateCheckResult struct {
	CurrentVersion string
	LatestVersion  string
	HasUpdate      bool
	Release        *GitHubRelease
	CheckedAt      time.Time
}

// CheckLatestRelease queries the GitHub API for the latest release and compares it with currentVersion
func CheckLatestRelease(currentVersion string) (*UpdateCheckResult, error) {
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", DefaultGitHubRepo)

	client := &http.Client{
		Timeout: RequestTimeout,
	}

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("gagal membuat request update: %w", err)
	}

	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("koneksi ke GitHub gagal (periksa koneksi internet Anda): %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server GitHub mengembalikan status: %s", resp.Status)
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("gagal memproses data rilis dari GitHub: %w", err)
	}

	latestVer := strings.TrimPrefix(strings.TrimSpace(release.TagName), "v")
	currVer := strings.TrimPrefix(strings.TrimSpace(currentVersion), "v")

	hasUpdate := IsNewerVersion(latestVer, currVer)

	return &UpdateCheckResult{
		CurrentVersion: currentVersion,
		LatestVersion:  release.TagName,
		HasUpdate:      hasUpdate,
		Release:        &release,
		CheckedAt:      time.Now(),
	}, nil
}

// IsNewerVersion returns true if candidate is strictly newer than baseline (SemVer comparison)
func IsNewerVersion(candidate, baseline string) bool {
	candParts := parseSemVer(candidate)
	baseParts := parseSemVer(baseline)

	for i := 0; i < 3; i++ {
		if candParts[i] > baseParts[i] {
			return true
		}
		if candParts[i] < baseParts[i] {
			return false
		}
	}
	return false
}

// parseSemVer parses "1.3.1" or "v1.3.1" into [3]int {1, 3, 1}
func parseSemVer(ver string) [3]int {
	clean := strings.TrimSpace(ver)
	clean = strings.TrimPrefix(clean, "v")
	clean = strings.TrimPrefix(clean, "V")

	// remove any build metadata or pre-release suffixes (e.g. -beta)
	if idx := strings.IndexAny(clean, "-+"); idx != -1 {
		clean = clean[:idx]
	}

	parts := strings.Split(clean, ".")
	var res [3]int
	for i := 0; i < len(parts) && i < 3; i++ {
		if n, err := strconv.Atoi(parts[i]); err == nil {
			res[i] = n
		}
	}
	return res
}

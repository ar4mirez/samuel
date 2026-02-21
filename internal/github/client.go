package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	// GitHubAPI is the base URL for GitHub API
	GitHubAPI = "https://api.github.com"

	// ArchiveURLTemplate is the template for downloading release archives
	// Format: https://github.com/{owner}/{repo}/archive/refs/tags/{tag}.tar.gz
	ArchiveURLTemplate = "https://github.com/%s/%s/archive/refs/tags/v%s.tar.gz"

	// BranchArchiveURLTemplate is the template for downloading branch archives
	// Format: https://github.com/{owner}/{repo}/archive/refs/heads/{branch}.tar.gz
	BranchArchiveURLTemplate = "https://github.com/%s/%s/archive/refs/heads/%s.tar.gz"

	// LatestReleaseURLTemplate is the template for fetching latest release info
	LatestReleaseURLTemplate = "https://api.github.com/repos/%s/%s/releases/latest"

	// TagsURLTemplate is the template for fetching tags
	TagsURLTemplate = "https://api.github.com/repos/%s/%s/tags"

	// DefaultBranch is the fallback when no releases exist
	DefaultBranch = "main"

	// DevVersion is returned when using main branch
	DevVersion = "dev"
)

// MaxDownloadFileSize is the maximum allowed size for a single file
// download (10 MB). Prevents out-of-memory from unbounded reads.
var MaxDownloadFileSize int64 = 10 * 1024 * 1024

// Client provides GitHub API operations
type Client struct {
	httpClient *http.Client
	owner      string
	repo       string
}

// NewClient creates a new GitHub client
func NewClient(owner, repo string) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		owner: owner,
		repo:  repo,
	}
}

// Release represents a GitHub release
type Release struct {
	TagName     string    `json:"tag_name"`
	Name        string    `json:"name"`
	Body        string    `json:"body"`
	PublishedAt time.Time `json:"published_at"`
	TarballURL  string    `json:"tarball_url"`
}

// Tag represents a GitHub tag
type Tag struct {
	Name string `json:"name"`
}

// GetLatestRelease fetches the latest release information
// Returns nil without error if no releases exist (use GetLatestVersionOrBranch instead)
func (c *Client) GetLatestRelease() (*Release, error) {
	url := fmt.Sprintf(LatestReleaseURLTemplate, c.owner, c.repo)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "samuel-cli")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch latest release: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		// No releases found - this is not an error, just no releases yet
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API error: %s", resp.Status)
	}

	var release Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("failed to parse release data: %w", err)
	}

	return &release, nil
}

// GetLatestVersionOrBranch returns the latest version, falling back to "dev" if no releases
func (c *Client) GetLatestVersionOrBranch() (version string, isBranch bool, err error) {
	release, err := c.GetLatestRelease()
	if err != nil {
		return "", false, err
	}

	if release != nil {
		// We have a release
		version := release.TagName
		if len(version) > 0 && version[0] == 'v' {
			version = version[1:]
		}
		return version, false, nil
	}

	// No releases - fall back to main branch
	return DevVersion, true, nil
}

// GetTags fetches available tags
func (c *Client) GetTags() ([]Tag, error) {
	url := fmt.Sprintf(TagsURLTemplate, c.owner, c.repo)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "samuel-cli")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tags: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API error: %s", resp.Status)
	}

	var tags []Tag
	if err := json.NewDecoder(resp.Body).Decode(&tags); err != nil {
		return nil, fmt.Errorf("failed to parse tags: %w", err)
	}

	return tags, nil
}

// GetArchiveURL returns the URL to download a specific version
func (c *Client) GetArchiveURL(version string) string {
	return fmt.Sprintf(ArchiveURLTemplate, c.owner, c.repo, version)
}

// GetBranchArchiveURL returns the URL to download from a branch
func (c *Client) GetBranchArchiveURL(branch string) string {
	return fmt.Sprintf(BranchArchiveURLTemplate, c.owner, c.repo, branch)
}

// DownloadArchive downloads the archive for a specific version
func (c *Client) DownloadArchive(version string) (io.ReadCloser, int64, error) {
	url := c.GetArchiveURL(version)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, 0, err
	}

	req.Header.Set("User-Agent", "samuel-cli")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to download archive: %w", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		resp.Body.Close()
		return nil, 0, fmt.Errorf("version %s not found", version)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, 0, fmt.Errorf("download failed: %s", resp.Status)
	}

	return resp.Body, resp.ContentLength, nil
}

// DownloadBranchArchive downloads the archive for a branch
func (c *Client) DownloadBranchArchive(branch string) (io.ReadCloser, int64, error) {
	url := c.GetBranchArchiveURL(branch)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, 0, err
	}

	req.Header.Set("User-Agent", "samuel-cli")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to download archive: %w", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		resp.Body.Close()
		return nil, 0, fmt.Errorf("branch %s not found", branch)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, 0, fmt.Errorf("download failed: %s", resp.Status)
	}

	return resp.Body, resp.ContentLength, nil
}

// DownloadFile downloads a single file from the repository
func (c *Client) DownloadFile(version, path string) ([]byte, error) {
	// Use raw.githubusercontent.com for direct file access
	url := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/v%s/%s",
		c.owner, c.repo, version, path)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "samuel-cli")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("file not found: %s", path)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download failed: %s", resp.Status)
	}

	data, err := io.ReadAll(io.LimitReader(resp.Body, MaxDownloadFileSize+1))
	if err != nil {
		return nil, fmt.Errorf("failed to read file content: %w", err)
	}
	if int64(len(data)) > MaxDownloadFileSize {
		return nil, fmt.Errorf("file %q exceeds maximum download size (%d bytes)", path, MaxDownloadFileSize)
	}
	return data, nil
}

// VersionInfo contains version comparison information
type VersionInfo struct {
	Current      string
	Latest       string
	UpdateNeeded bool
	ReleaseNotes string
}

// CheckForUpdates compares current version with latest
func (c *Client) CheckForUpdates(currentVersion string) (*VersionInfo, error) {
	release, err := c.GetLatestRelease()
	if err != nil {
		return nil, err
	}

	if release == nil {
		return nil, fmt.Errorf("no releases found for %s/%s", c.owner, c.repo)
	}

	// Strip 'v' prefix if present
	latestVersion := release.TagName
	if len(latestVersion) > 0 && latestVersion[0] == 'v' {
		latestVersion = latestVersion[1:]
	}

	return &VersionInfo{
		Current:      currentVersion,
		Latest:       latestVersion,
		UpdateNeeded: latestVersion != currentVersion,
		ReleaseNotes: release.Body,
	}, nil
}

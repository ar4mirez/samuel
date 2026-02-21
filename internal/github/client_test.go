package github

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// redirectTransport rewrites all outgoing requests to hit the test server.
type redirectTransport struct {
	server *httptest.Server
}

func (t *redirectTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.URL.Scheme = "http"
	req.URL.Host = t.server.Listener.Addr().String()
	return http.DefaultTransport.RoundTrip(req)
}

// newTestClient creates a Client whose HTTP requests go to the test server.
func newTestClient(server *httptest.Server) *Client {
	return &Client{
		httpClient: &http.Client{
			Transport: &redirectTransport{server: server},
		},
		owner: "testowner",
		repo:  "testrepo",
	}
}

func TestNewClient(t *testing.T) {
	c := NewClient("myowner", "myrepo")
	if c.owner != "myowner" {
		t.Errorf("owner = %q, want %q", c.owner, "myowner")
	}
	if c.repo != "myrepo" {
		t.Errorf("repo = %q, want %q", c.repo, "myrepo")
	}
	if c.httpClient == nil {
		t.Error("httpClient should not be nil")
	}
}

func TestGetLatestRelease(t *testing.T) {
	tests := []struct {
		name    string
		handler http.HandlerFunc
		wantNil bool
		wantErr bool
		wantTag string
	}{
		{
			name: "success",
			handler: func(w http.ResponseWriter, r *http.Request) {
				_ = json.NewEncoder(w).Encode(Release{
					TagName: "v1.2.3",
					Name:    "Release 1.2.3",
					Body:    "Release notes",
				})
			},
			wantTag: "v1.2.3",
		},
		{
			name: "not_found_returns_nil",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			},
			wantNil: true,
		},
		{
			name: "server_error",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			wantErr: true,
		},
		{
			name: "invalid_json",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				_, _ = w.Write([]byte("not json"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.handler)
			defer server.Close()
			client := newTestClient(server)

			release, err := client.GetLatestRelease()
			if (err != nil) != tt.wantErr {
				t.Fatalf("error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantNil && release != nil {
				t.Errorf("got %+v, want nil", release)
			}
			if tt.wantTag != "" && (release == nil || release.TagName != tt.wantTag) {
				t.Errorf("TagName = %v, want %q", release, tt.wantTag)
			}
		})
	}
}

func TestGetLatestRelease_Headers(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Accept"); got != "application/vnd.github.v3+json" {
			t.Errorf("Accept header = %q, want %q", got, "application/vnd.github.v3+json")
		}
		if got := r.Header.Get("User-Agent"); got != "samuel-cli" {
			t.Errorf("User-Agent = %q, want %q", got, "samuel-cli")
		}
		_ = json.NewEncoder(w).Encode(Release{TagName: "v1.0.0"})
	}))
	defer server.Close()

	client := newTestClient(server)
	_, _ = client.GetLatestRelease()
}

func TestGetLatestVersionOrBranch(t *testing.T) {
	tests := []struct {
		name       string
		handler    http.HandlerFunc
		wantVer    string
		wantBranch bool
		wantErr    bool
	}{
		{
			name: "strips_v_prefix",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				_ = json.NewEncoder(w).Encode(Release{TagName: "v1.2.3"})
			},
			wantVer: "1.2.3",
		},
		{
			name: "no_v_prefix",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				_ = json.NewEncoder(w).Encode(Release{TagName: "1.0.0"})
			},
			wantVer: "1.0.0",
		},
		{
			name: "no_releases_falls_back",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			},
			wantVer:    DevVersion,
			wantBranch: true,
		},
		{
			name: "api_error",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.handler)
			defer server.Close()
			client := newTestClient(server)

			ver, isBranch, err := client.GetLatestVersionOrBranch()
			if (err != nil) != tt.wantErr {
				t.Fatalf("error = %v, wantErr %v", err, tt.wantErr)
			}
			if ver != tt.wantVer {
				t.Errorf("version = %q, want %q", ver, tt.wantVer)
			}
			if isBranch != tt.wantBranch {
				t.Errorf("isBranch = %v, want %v", isBranch, tt.wantBranch)
			}
		})
	}
}

func TestGetTags(t *testing.T) {
	tests := []struct {
		name    string
		handler http.HandlerFunc
		wantLen int
		wantErr bool
	}{
		{
			name: "multiple_tags",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				_ = json.NewEncoder(w).Encode([]Tag{
					{Name: "v1.0.0"},
					{Name: "v0.9.0"},
				})
			},
			wantLen: 2,
		},
		{
			name: "empty_tags",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				_ = json.NewEncoder(w).Encode([]Tag{})
			},
			wantLen: 0,
		},
		{
			name: "server_error",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			wantErr: true,
		},
		{
			name: "invalid_json",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				_, _ = w.Write([]byte("{bad"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.handler)
			defer server.Close()
			client := newTestClient(server)

			tags, err := client.GetTags()
			if (err != nil) != tt.wantErr {
				t.Fatalf("error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && len(tags) != tt.wantLen {
				t.Errorf("len(tags) = %d, want %d", len(tags), tt.wantLen)
			}
		})
	}
}

func TestGetArchiveURL(t *testing.T) {
	c := NewClient("owner", "repo")
	got := c.GetArchiveURL("1.2.3")
	want := "https://github.com/owner/repo/archive/refs/tags/v1.2.3.tar.gz"
	if got != want {
		t.Errorf("GetArchiveURL() = %q, want %q", got, want)
	}
}

func TestGetBranchArchiveURL(t *testing.T) {
	c := NewClient("owner", "repo")
	got := c.GetBranchArchiveURL("main")
	want := "https://github.com/owner/repo/archive/refs/heads/main.tar.gz"
	if got != want {
		t.Errorf("GetBranchArchiveURL() = %q, want %q", got, want)
	}
}

func TestDownloadArchive(t *testing.T) {
	tests := []struct {
		name    string
		handler http.HandlerFunc
		wantErr bool
		errMsg  string
		wantBody string
	}{
		{
			name: "success",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				_, _ = w.Write([]byte("archive data"))
			},
			wantBody: "archive data",
		},
		{
			name: "not_found",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			},
			wantErr: true,
			errMsg:  "not found",
		},
		{
			name: "server_error",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			wantErr: true,
			errMsg:  "download failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.handler)
			defer server.Close()
			client := newTestClient(server)

			body, _, err := client.DownloadArchive("1.0.0")
			if (err != nil) != tt.wantErr {
				t.Fatalf("error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("error = %q, want containing %q", err, tt.errMsg)
				}
				return
			}
			defer body.Close()
			data, _ := io.ReadAll(body)
			if string(data) != tt.wantBody {
				t.Errorf("body = %q, want %q", string(data), tt.wantBody)
			}
		})
	}
}

func TestDownloadBranchArchive(t *testing.T) {
	tests := []struct {
		name     string
		handler  http.HandlerFunc
		wantErr  bool
		errMsg   string
		wantBody string
	}{
		{
			name: "success",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				_, _ = w.Write([]byte("branch data"))
			},
			wantBody: "branch data",
		},
		{
			name: "not_found",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			},
			wantErr: true,
			errMsg:  "not found",
		},
		{
			name: "server_error",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			wantErr: true,
			errMsg:  "download failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.handler)
			defer server.Close()
			client := newTestClient(server)

			body, _, err := client.DownloadBranchArchive("main")
			if (err != nil) != tt.wantErr {
				t.Fatalf("error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("error = %q, want containing %q", err, tt.errMsg)
				}
				return
			}
			defer body.Close()
			data, _ := io.ReadAll(body)
			if string(data) != tt.wantBody {
				t.Errorf("body = %q, want %q", string(data), tt.wantBody)
			}
		})
	}
}

func TestDownloadFile(t *testing.T) {
	tests := []struct {
		name    string
		handler http.HandlerFunc
		wantErr bool
		errMsg  string
		want    string
	}{
		{
			name: "success",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				_, _ = w.Write([]byte("file content"))
			},
			want: "file content",
		},
		{
			name: "not_found",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			},
			wantErr: true,
			errMsg:  "file not found",
		},
		{
			name: "server_error",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			wantErr: true,
			errMsg:  "download failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.handler)
			defer server.Close()
			client := newTestClient(server)

			data, err := client.DownloadFile("1.0.0", "README.md")
			if (err != nil) != tt.wantErr {
				t.Fatalf("error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("error = %q, want containing %q", err, tt.errMsg)
				}
				return
			}
			if string(data) != tt.want {
				t.Errorf("data = %q, want %q", string(data), tt.want)
			}
		})
	}
}

func TestDownloadFile_SizeLimit(t *testing.T) {
	origLimit := MaxDownloadFileSize
	MaxDownloadFileSize = 512 // 512-byte limit for testing
	defer func() { MaxDownloadFileSize = origLimit }()

	t.Run("oversized_file_rejected", func(t *testing.T) {
		oversized := strings.Repeat("x", 1024) // 1KB > 512-byte limit
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			_, _ = w.Write([]byte(oversized))
		}))
		defer server.Close()
		client := newTestClient(server)

		_, err := client.DownloadFile("1.0.0", "big-file.txt")
		if err == nil {
			t.Fatal("expected error for oversized file, got nil")
		}
		if !strings.Contains(err.Error(), "exceeds maximum download size") {
			t.Errorf("expected 'exceeds maximum download size' error, got: %v", err)
		}
	})

	t.Run("file_at_limit_succeeds", func(t *testing.T) {
		exactSize := strings.Repeat("y", 512) // exactly at limit
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			_, _ = w.Write([]byte(exactSize))
		}))
		defer server.Close()
		client := newTestClient(server)

		data, err := client.DownloadFile("1.0.0", "ok-file.txt")
		if err != nil {
			t.Fatalf("file at exact limit should succeed, got: %v", err)
		}
		if len(data) != 512 {
			t.Errorf("expected 512 bytes, got %d", len(data))
		}
	})
}

func TestCheckForUpdates(t *testing.T) {
	tests := []struct {
		name       string
		handler    http.HandlerFunc
		current    string
		wantUpdate bool
		wantLatest string
		wantErr    bool
	}{
		{
			name: "update_available",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				_ = json.NewEncoder(w).Encode(Release{
					TagName: "v2.0.0",
					Body:    "New features",
				})
			},
			current:    "1.0.0",
			wantUpdate: true,
			wantLatest: "2.0.0",
		},
		{
			name: "up_to_date",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				_ = json.NewEncoder(w).Encode(Release{TagName: "v1.0.0"})
			},
			current:    "1.0.0",
			wantLatest: "1.0.0",
		},
		{
			name: "no_releases_returns_error",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			},
			current: "1.0.0",
			wantErr: true,
		},
		{
			name: "api_error",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			current: "1.0.0",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.handler)
			defer server.Close()
			client := newTestClient(server)

			info, err := client.CheckForUpdates(tt.current)
			if (err != nil) != tt.wantErr {
				t.Fatalf("error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if info.UpdateNeeded != tt.wantUpdate {
				t.Errorf("UpdateNeeded = %v, want %v", info.UpdateNeeded, tt.wantUpdate)
			}
			if info.Current != tt.current {
				t.Errorf("Current = %q, want %q", info.Current, tt.current)
			}
			if info.Latest != tt.wantLatest {
				t.Errorf("Latest = %q, want %q", info.Latest, tt.wantLatest)
			}
		})
	}
}

func TestConstants(t *testing.T) {
	if DefaultBranch != "main" {
		t.Errorf("DefaultBranch = %q, want %q", DefaultBranch, "main")
	}
	if DevVersion != "dev" {
		t.Errorf("DevVersion = %q, want %q", DevVersion, "dev")
	}
}

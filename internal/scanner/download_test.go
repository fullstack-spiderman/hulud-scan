package scanner

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDownloadBlocklist(t *testing.T) {
	// Create mock HTTP server with Wiz format
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		csv := `Package,Version
lodash,= 4.17.20
express,= 4.17.1`
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(csv))
	}))
	defer server.Close()

	// Test download
	blocklist, err := DownloadBlocklist(server.URL)

	require.NoError(t, err)
	require.NotNil(t, blocklist)
	assert.Len(t, blocklist.Entries, 2)
	assert.Equal(t, "lodash", blocklist.Entries[0].PackageName)
	assert.Equal(t, "4.17.20", blocklist.Entries[0].Version)
	assert.Equal(t, SeverityCritical, blocklist.Entries[0].Severity)
	assert.Contains(t, blocklist.Entries[0].Reason, "Shai-Hulud")
}

func TestDownloadBlocklist_FullFormat(t *testing.T) {
	// Create mock HTTP server with full format
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		csv := `package_name,version,severity,reason,cve
lodash,4.17.20,critical,Prototype pollution,CVE-2020-8203
express,4.17.1,high,DoS vulnerability,CVE-2022-24999`
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(csv))
	}))
	defer server.Close()

	blocklist, err := DownloadBlocklist(server.URL)

	require.NoError(t, err)
	assert.Len(t, blocklist.Entries, 2)
	assert.Equal(t, "Prototype pollution", blocklist.Entries[0].Reason)
	assert.Equal(t, "CVE-2020-8203", blocklist.Entries[0].CVE)
}

func TestDownloadBlocklist_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	_, err := DownloadBlocklist(server.URL)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "HTTP 404")
}

func TestConvertToRawURL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "GitHub blob URL",
			input:    "https://github.com/user/repo/blob/main/file.csv",
			expected: "https://raw.githubusercontent.com/user/repo/main/file.csv",
		},
		{
			name:     "Already raw URL",
			input:    "https://raw.githubusercontent.com/user/repo/main/file.csv",
			expected: "https://raw.githubusercontent.com/user/repo/main/file.csv",
		},
		{
			name:     "Non-GitHub URL",
			input:    "https://example.com/blocklist.csv",
			expected: "https://example.com/blocklist.csv",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convertToRawURL(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestLoadOrDownloadBlocklist_LocalFile(t *testing.T) {
	blocklist, err := LoadOrDownloadBlocklist("../../testdata/sample-blocklist.csv", "")

	require.NoError(t, err)
	assert.NotEmpty(t, blocklist.Entries)
}

func TestLoadOrDownloadBlocklist_WithCaching(t *testing.T) {
	tmpDir := t.TempDir()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		csv := `Package,Version
test-package,= 1.0.0`
		w.Write([]byte(csv))
	}))
	defer server.Close()

	// First load - should download
	blocklist1, err := LoadOrDownloadBlocklist(server.URL, tmpDir)
	require.NoError(t, err)
	assert.Equal(t, "test-package", blocklist1.Entries[0].PackageName)

	// Verify cache file was created
	files, err := os.ReadDir(tmpDir)
	require.NoError(t, err)
	assert.Len(t, files, 1)
	assert.Contains(t, files[0].Name(), "blocklist-")

	// Second load - should use cache (within TTL)
	blocklist2, err := LoadOrDownloadBlocklist(server.URL, tmpDir)
	require.NoError(t, err)
	assert.Equal(t, blocklist1.Entries[0].PackageName, blocklist2.Entries[0].PackageName)
}

func TestLoadOrDownloadBlocklist_CacheExpiry(t *testing.T) {
	tmpDir := t.TempDir()

	downloadCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		downloadCount++
		csv := `Package,Version
test,= 1.0.0`
		w.Write([]byte(csv))
	}))
	defer server.Close()

	// First download
	_, err := LoadOrDownloadBlocklist(server.URL, tmpDir)
	require.NoError(t, err)
	assert.Equal(t, 1, downloadCount)

	// Get cache file
	files, err := os.ReadDir(tmpDir)
	require.NoError(t, err)
	cacheFile := filepath.Join(tmpDir, files[0].Name())

	// Modify cache file timestamp to make it expired
	oldTime := time.Now().Add(-2 * time.Hour)
	os.Chtimes(cacheFile, oldTime, oldTime)

	// Second load - should re-download due to expiry
	_, err = LoadOrDownloadBlocklist(server.URL, tmpDir)
	require.NoError(t, err)
	assert.Equal(t, 2, downloadCount, "Should have downloaded twice due to cache expiry")
}

func TestParseWizEntry_MultipleVersions(t *testing.T) {
	record := []string{"test-package", "= 1.0.0 || = 1.0.1"}

	entry := parseWizEntry(record)

	assert.Equal(t, "test-package", entry.PackageName)
	assert.Equal(t, "1.0.0", entry.Version) // Should take first version
	assert.Equal(t, SeverityCritical, entry.Severity)
}

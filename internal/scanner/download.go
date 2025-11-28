package scanner

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// DownloadBlocklist downloads a blocklist from a URL
func DownloadBlocklist(url string) (*Blocklist, error) {
	// Convert GitHub web URL to raw URL if needed
	url = convertToRawURL(url)

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	fmt.Printf("   Downloading from: %s\n", url)

	// Download the CSV
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to download blocklist: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			fmt.Printf("Warning: failed to close response body: %v\n", closeErr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download blocklist: HTTP %d", resp.StatusCode)
	}

	// Read response body
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Parse CSV
	reader := csv.NewReader(strings.NewReader(string(data)))
	return parseBlocklistCSV(reader)
}

// LoadOrDownloadBlocklist loads from file or downloads from URL
func LoadOrDownloadBlocklist(path string, cacheDir string) (*Blocklist, error) {
	// Check if it's a URL
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		// Check cache first
		if cacheDir != "" {
			cached, err := loadFromCache(path, cacheDir)
			if err == nil {
				fmt.Printf("   Using cached blocklist\n")
				return cached, nil
			}
		}

		// Download
		blocklist, err := DownloadBlocklist(path)
		if err != nil {
			// Try to use expired cache as fallback
			if cacheDir != "" {
				cached, cacheErr := loadFromCacheIgnoreExpiry(path, cacheDir)
				if cacheErr == nil {
					fmt.Printf("   ⚠️  Download failed, using cached version (may be outdated)\n")
					return cached, nil
				}
			}
			return nil, err
		}

		// Save to cache
		if cacheDir != "" {
			if err := saveToCache(path, cacheDir, blocklist); err != nil {
				// Non-fatal - just log
				fmt.Printf("   Warning: failed to cache blocklist: %v\n", err)
			}
		}

		return blocklist, nil
	}

	// Load from local file
	return LoadBlocklist(path)
}

// convertToRawURL converts GitHub web URLs to raw content URLs
func convertToRawURL(url string) string {
	// Convert github.com/user/repo/blob/branch/file
	// To: raw.githubusercontent.com/user/repo/branch/file
	if strings.Contains(url, "github.com") && strings.Contains(url, "/blob/") {
		url = strings.Replace(url, "github.com", "raw.githubusercontent.com", 1)
		url = strings.Replace(url, "/blob/", "/", 1)
	}
	return url
}

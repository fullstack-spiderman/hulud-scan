package scanner

import (
	"crypto/sha256"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const cacheTTL = 1 * time.Hour // Cache for 1 hour

// loadFromCache loads blocklist from cache if not expired
func loadFromCache(url string, cacheDir string) (*Blocklist, error) {
	cachePath := getCachePath(url, cacheDir)

	// Check if cache file exists
	info, err := os.Stat(cachePath)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("cache miss")
	}

	// Check if cache is expired
	if time.Since(info.ModTime()) > cacheTTL {
		return nil, fmt.Errorf("cache expired")
	}

	// Load from cache
	return LoadBlocklist(cachePath)
}

// loadFromCacheIgnoreExpiry loads cache even if expired (for fallback)
func loadFromCacheIgnoreExpiry(url string, cacheDir string) (*Blocklist, error) {
	cachePath := getCachePath(url, cacheDir)

	// Check if cache file exists
	if _, err := os.Stat(cachePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("cache miss")
	}

	// Load from cache (ignore expiry)
	return LoadBlocklist(cachePath)
}

// saveToCache saves blocklist to cache
func saveToCache(url string, cacheDir string, blocklist *Blocklist) error {
	cachePath := getCachePath(url, cacheDir)

	// Ensure cache directory exists
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return fmt.Errorf("failed to create cache dir: %w", err)
	}

	// Write CSV
	file, err := os.Create(cachePath)
	if err != nil {
		return fmt.Errorf("failed to create cache file: %w", err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			fmt.Printf("Warning: failed to close cache file: %v\n", closeErr)
		}
	}()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	if err := writer.Write([]string{"package_name", "version", "severity", "reason", "cve"}); err != nil {
		return err
	}

	// Write entries
	for _, entry := range blocklist.Entries {
		record := []string{
			entry.PackageName,
			entry.Version,
			string(entry.Severity),
			entry.Reason,
			entry.CVE,
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return writer.Error()
}

// getCachePath generates cache file path from URL
func getCachePath(url string, cacheDir string) string {
	// Hash the URL to create a unique filename
	hash := sha256.Sum256([]byte(url))
	filename := fmt.Sprintf("blocklist-%x.csv", hash[:8])
	return filepath.Join(cacheDir, filename)
}

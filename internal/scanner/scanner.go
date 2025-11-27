package scanner

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"github.com/arjunu/hulud-scan/internal/graph"
)

// LoadBlocklist loads a blocklist from a CSV file
func LoadBlocklist(path string) (*Blocklist, error) {
	// Open the CSV file
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open blocklist file: %w", err)
	}
	defer file.Close() // Close file when function returns

	// Create CSV reader
	reader := csv.NewReader(file)

	// Read all records
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV: %w", err)
	}

	// Check if file is empty
	if len(records) < 2 {
		return nil, fmt.Errorf("blocklist file is empty or missing header")
	}

	// First row is header, skip it
	records = records[1:]

	// Parse entries
	entries := make([]BlocklistEntry, 0, len(records))
	index := make(map[string][]int)

	for i, record := range records {
		// CSV format: package_name,version,severity,reason,cve
		if len(record) < 4 {
			continue // Skip malformed rows
		}

		entry := BlocklistEntry{
			PackageName: strings.TrimSpace(record[0]),
			Version:     strings.TrimSpace(record[1]),
			Severity:    Severity(strings.TrimSpace(record[2])),
			Reason:      strings.TrimSpace(record[3]),
		}

		// CVE is optional (column 5)
		if len(record) >= 5 {
			entry.CVE = strings.TrimSpace(record[4])
		}

		entries = append(entries, entry)

		// Build index for fast lookup
		index[entry.PackageName] = append(index[entry.PackageName], i)
	}

	return &Blocklist{
		Entries: entries,
		Index:   index,
	}, nil
}

// IsBlocked checks if a specific package version is in the blocklist
func (b *Blocklist) IsBlocked(packageName, version string) *BlocklistEntry {
	// Use index to find entries for this package
	indices, exists := b.Index[packageName]
	if !exists {
		return nil // Package not in blocklist at all
	}

	// Check each entry for this package
	for _, idx := range indices {
		entry := &b.Entries[idx]
		if entry.Version == version {
			return entry // Found a match!
		}
	}

	return nil // Package exists in blocklist but not this version
}

// ScanGraph scans a dependency graph against a blocklist
func ScanGraph(g *graph.Graph, blocklist *Blocklist) *ScanResult {
	result := &ScanResult{
		Findings:      make([]Finding, 0),
		TotalPackages: len(g.Nodes),
		IssuesFound:   0,
	}

	// Scan each package in the graph
	for path, node := range g.Nodes {
		pkg := node.Package

		// Check if this package/version is blocklisted
		entry := blocklist.IsBlocked(pkg.Name, pkg.Version)
		if entry != nil {
			// Found a compromised package!
			finding := Finding{
				PackageName: pkg.Name,
				Version:     pkg.Version,
				Path:        g.FindPath(path),
				Severity:    entry.Severity,
				Reason:      entry.Reason,
				CVE:         entry.CVE,
				IsDirect:    node.IsDirect,
			}

			result.Findings = append(result.Findings, finding)
			result.IssuesFound++
		}
	}

	return result
}

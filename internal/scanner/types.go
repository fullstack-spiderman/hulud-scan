package scanner

import "github.com/fullstack-spiderman/hulud-scan/internal/graph"

// Severity levels for security findings
type Severity string

const (
	SeverityCritical Severity = "critical"
	SeverityHigh     Severity = "high"
	SeverityMedium   Severity = "medium"
	SeverityLow      Severity = "low"
	SeverityInfo     Severity = "info"
)

// BlocklistEntry represents a known compromised package version
type BlocklistEntry struct {
	PackageName string   // Name of the compromised package
	Version     string   // Affected version
	Severity    Severity // How serious is this?
	Reason      string   // Why is it flagged?
	CVE         string   // CVE identifier (if applicable)
}

// Blocklist is a collection of known compromised packages
type Blocklist struct {
	Entries []BlocklistEntry // All blocklist entries
	Index   map[string][]int // Index: package name -> entry indices (for fast lookup)
}

// Finding represents a security issue found during scanning
type Finding struct {
	PackageName string               // Package that was flagged
	Version     string               // Version that was flagged
	Path        graph.DependencyPath // How we got to this package
	Severity    Severity             // Severity of the issue
	Reason      string               // Why it was flagged
	CVE         string               // CVE if applicable
	IsDirect    bool                 // Is this a direct dependency?
}

// ScanResult contains all findings from a scan
type ScanResult struct {
	Findings      []Finding // All security findings
	TotalPackages int       // Total packages scanned
	IssuesFound   int       // Number of issues found
}

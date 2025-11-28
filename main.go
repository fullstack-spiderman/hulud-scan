package main

import (
	"os"

	"github.com/fullstack-spiderman/hulud-scan/cmd"
)

// Version information set via ldflags during build
var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

func main() {
	// Set version info for cmd package
	cmd.Version = Version
	cmd.Commit = Commit
	cmd.Date = Date

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

package main

import (
	"os"

	"github.com/arjunu/hulud-scan/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

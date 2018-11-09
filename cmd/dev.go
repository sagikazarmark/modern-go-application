// +build dev

package main

import (
	"time"
)

func init() {
	// Load defaults for info variables
	if Version == "" {
		Version = "dev"
	}

	if CommitHash == "" {
		CommitHash = "dev"
	}

	if BuildDate == "" {
		BuildDate = time.Now().Format(time.RFC3339)
	}
}

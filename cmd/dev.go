// +build dev

package main

import (
	"time"
)

func init() {
	// Load defaults for info variables
	if version == "" {
		version = "dev"
	}

	if commitHash == "" {
		commitHash = "dev"
	}

	if buildDate == "" {
		buildDate = time.Now().Format(time.RFC3339)
	}
}

// +build dev

package main

import (
	"path"
	"runtime"
	"time"

	"github.com/joho/godotenv"
)

func init() {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("cannot get current dir: no caller information")
	}

	projectRoot := path.Clean(path.Dir(path.Dir(filename)))

	godotenv.Load(path.Join(projectRoot, ".env"))
	godotenv.Load(path.Join(projectRoot, ".env.dist"))

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

	if Build == "" {
		Build = "dev"
	}
}

package main

import (
	"encoding/json"
	"net/http"
	"runtime"

	"github.com/pkg/errors"
)

// Provisioned by ldflags
var (
	Version    string
	CommitHash string
	BuildDate  string
)

// BuildInfo represents all available build information.
type BuildInfo struct {
	Version    string `json:"version"`
	CommitHash string `json:"commit_hash"`
	BuildDate  string `json:"build_date"`
	GoVersion  string `json:"go_version"`
	Os         string `json:"os"`
	Arch       string `json:"arch"`
	Compiler   string `json:"compiler"`
}

// NewBuildInfo returns all available build information.
func NewBuildInfo() BuildInfo {
	return BuildInfo{
		Version:    Version,
		CommitHash: CommitHash,
		BuildDate:  BuildDate,
		GoVersion:  runtime.Version(),
		Os:         runtime.GOOS,
		Arch:       runtime.GOARCH,
		Compiler:   runtime.Compiler,
	}
}

// Handler returns an HTTP handler for version information.
func (i BuildInfo) Handler() http.Handler {
	var body []byte

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if body == nil {
			var err error

			body, err = json.Marshal(i)
			if err != nil {
				panic(errors.Wrap(err, "failed to render version information"))
			}
		}

		_, _ = w.Write(body)
	})
}

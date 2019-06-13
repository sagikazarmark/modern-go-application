package buildinfo

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	buildinfo := New("version", "commit", "date")

	assert.Equal(t, "version", buildinfo.Version)
	assert.Equal(t, "commit", buildinfo.CommitHash)
	assert.Equal(t, "date", buildinfo.BuildDate)
}

func TestBuildInfo_Fields(t *testing.T) {
	buildinfo := New("version", "commit", "date")

	actualFields := buildinfo.Fields()

	expectedFields := map[string]interface{}{
		"version":     "version",
		"commit_hash": "commit",
		"build_date":  "date",
		"go_version":  runtime.Version(),
		"os":          runtime.GOOS,
		"arch":        runtime.GOARCH,
		"compiler":    runtime.Compiler,
	}

	assert.Equal(t, expectedFields, actualFields)
}

package buildinfo

import (
	"encoding/json"
	"net/http/httptest"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	buildinfo := New("version", "commit", "date")

	server := httptest.NewServer(Handler(buildinfo))
	defer server.Close()

	resp, err := server.Client().Get(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	var actualFields map[string]interface{}

	decoder := json.NewDecoder(resp.Body)

	err = decoder.Decode(&actualFields)
	if err != nil {
		t.Fatal(err)
	}

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

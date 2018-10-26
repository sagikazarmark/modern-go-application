package maintenance

import (
	"encoding/json"
	"net/http"
	"runtime"

	"github.com/pkg/errors"
)

// NewVersionHandler returns an HTTP handler for returning version information.
func NewVersionHandler(version string, commitHash string, buildDate string) http.Handler {
	data := map[string]interface{}{
		"version":     version,
		"commit_hash": commitHash,
		"build_date":  buildDate,
		"go_version":  runtime.Version(),
		"os":          runtime.GOOS,
		"arch":        runtime.GOARCH,
		"compiler":    runtime.Compiler,
	}
	body, err := json.Marshal(data)
	if err != nil {
		panic(errors.Wrap(err, "failed to render version information"))
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(body)
	})
}

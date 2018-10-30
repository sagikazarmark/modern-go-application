package maintenance

import (

"encoding/json"
"net/http"

"github.com/pkg/errors"
)

// NewVersionHandler returns an HTTP handler for returning version information.
func NewVersionHandler(buildInfo interface{}) http.Handler {
	body, err := json.Marshal(buildInfo)
	if err != nil {
		panic(errors.Wrap(err, "failed to render version information"))
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(body)
	})
}

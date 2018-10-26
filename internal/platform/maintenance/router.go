package maintenance

import "net/http"

// NewRouter returns a new maintenance router.
func NewRouter() *http.ServeMux {
	router := http.NewServeMux()

	return router
}

package appkit

import (
	"emperror.dev/errors/match"
	"github.com/go-kit/kit/endpoint"
	kitxendpoint "github.com/sagikazarmark/kitx/endpoint"
)

// ClientErrorMiddleware checks if a returned error is a client error and wraps it in a failer response if it is.
func ClientErrorMiddleware(e endpoint.Endpoint) endpoint.Endpoint {
	return kitxendpoint.FailerMiddleware(match.ErrorMatcherFunc(IsClientError))(e)
}

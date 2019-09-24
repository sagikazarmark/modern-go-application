package http

import (
	"context"
	"encoding/json"
	"net/http"

	"emperror.dev/errors"
	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/moogar0880/problems"
)

// NopResponseEncoder can be used for operations without output parameters.
// It returns 200 OK status code without a response body.
func NopResponseEncoder(_ context.Context, _ http.ResponseWriter, _ interface{}) error {
	return nil
}

// StatusCodeResponseEncoder can be used for operations without output parameters.
// It returns 200 OK status code without a response body.
func StatusCodeResponseEncoder(code int) kithttp.EncodeResponseFunc {
	return func(_ context.Context, w http.ResponseWriter, _ interface{}) error {
		w.WriteHeader(code)

		return nil
	}
}

// JSONResponseEncoder encodes the passed response object to the HTTP response writer in JSON format.
func JSONResponseEncoder(ctx context.Context, w http.ResponseWriter, resp interface{}) error {
	err := kithttp.EncodeJSONResponse(ctx, w, resp)
	if err != nil {
		return errors.Wrap(err, "failed to encode response")
	}

	return nil
}

// EncodeErrorResponseFunc encodes the passed error to the HTTP response writer.
// It's designed to be used in HTTP servers, for server-side endpoints.
// An EncodeErrorResponseFunc is supposed to return an error with the proper HTTP status code.
type EncodeErrorResponseFunc func(context.Context, http.ResponseWriter, error) error

// ErrorResponseEncoder encodes the passed response object to the HTTP response writer in JSON format.
func ErrorResponseEncoder(
	encoder kithttp.EncodeResponseFunc,
	errorEncoder EncodeErrorResponseFunc,
) kithttp.EncodeResponseFunc {
	return func(ctx context.Context, w http.ResponseWriter, resp interface{}) error {
		if f, ok := resp.(endpoint.Failer); ok && f.Failed() != nil {
			return errorEncoder(ctx, w, f.Failed())
		}

		return encoder(ctx, w, resp)
	}
}

// ProblemErrorEncoder encodes errors in the Problem RFC format.
func ProblemErrorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	problem := problems.NewDetailedProblem(http.StatusInternalServerError, "something went wrong")

	w.Header().Set("Content-Type", problems.ProblemMediaType)
	w.WriteHeader(problem.Status)

	_ = json.NewEncoder(w).Encode(problem)
}

package tododriver

import (
	"context"
	"encoding/json"
	"net/http"

	"emperror.dev/errors"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/moogar0880/problems"
	kitxhttp "github.com/sagikazarmark/kitx/transport/http"

	"github.com/sagikazarmark/modern-go-application/internal/app/mga/todo"
)

// RegisterHTTPHandlers mounts all of the service endpoints into a router.
func RegisterHTTPHandlers(endpoints Endpoints, router *mux.Router, options ...kithttp.ServerOption) {
	router.Methods(http.MethodPost).Path("").Handler(kithttp.NewServer(
		endpoints.CreateTodo,
		decodeCreateTodoHTTPRequest,
		encodeCreateTodoHTTPResponse,
		options...,
	))

	router.Methods(http.MethodGet).Path("").Handler(kithttp.NewServer(
		endpoints.ListTodos,
		kithttp.NopRequestDecoder,
		kitxhttp.JSONResponseEncoder,
		options...,
	))

	router.Methods(http.MethodPost).Path("/{id}/done").Handler(kithttp.NewServer(
		endpoints.MarkAsDone,
		decodeMarkAsDoneHTTPRequest,
		kitxhttp.ErrorResponseEncoder(kitxhttp.NopResponseEncoder, errorEncoder),
		options...,
	))
}

func decodeCreateTodoHTTPRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req createTodoRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return req, errors.Wrap(err, "failed to decode request")
	}

	return req, nil
}

func encodeCreateTodoHTTPResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return kitxhttp.JSONResponseEncoder(ctx, w, kitxhttp.WithStatusCode(response, http.StatusCreated))
}

func decodeMarkAsDoneHTTPRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)

	id, ok := vars["id"]
	if !ok || id == "" {
		return nil, errors.NewWithDetails("missing parameter from the URL", "param", "id")
	}

	return markAsDoneRequest{
		ID: id,
	}, nil
}

func errorEncoder(_ context.Context, w http.ResponseWriter, err error) error {
	status := http.StatusInternalServerError

	// nolint: gocritic
	switch {
	case errors.As(err, &todo.NotFoundError{}):
		status = http.StatusNotFound
	}

	problem := problems.NewDetailedProblem(status, err.Error())

	w.Header().Set("Content-Type", problems.ProblemMediaType)
	w.WriteHeader(status)
	e := json.NewEncoder(w).Encode(problem)

	return e
}

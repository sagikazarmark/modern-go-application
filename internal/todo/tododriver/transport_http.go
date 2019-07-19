package tododriver

import (
	"context"
	"encoding/json"
	"net/http"

	"emperror.dev/emperror"
	"emperror.dev/errors"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/moogar0880/problems"
	"github.com/sagikazarmark/ocmux"

	api "github.com/sagikazarmark/modern-go-application/.gen/api/openapi/todo/go"
	"github.com/sagikazarmark/modern-go-application/internal/todo"
)

// MakeHTTPHandler mounts all of the service endpoints into an http.Handler.
func MakeHTTPHandler(endpoints Endpoints, errorHandler todo.ErrorHandler) http.Handler {
	r := mux.NewRouter().PathPrefix("/todos").Subrouter()
	r.Use(ocmux.Middleware())

	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodeHTTPError),
		httptransport.ServerErrorHandler(emperror.MakeContextAware(errorHandler)),
	}

	r.Methods(http.MethodPost).Path("/").Handler(httptransport.NewServer(
		endpoints.Create,
		decodeCreateTodoHTTPRequest,
		encodeCreateTodoHTTPResponse,
		options...,
	))

	r.Methods(http.MethodGet).Path("/").Handler(httptransport.NewServer(
		endpoints.List,
		decodeListTodosHTTPRequest,
		encodeListTodosHTTPResponse,
		options...,
	))

	r.Methods(http.MethodPost).Path("/{id}/done").Handler(httptransport.NewServer(
		endpoints.MarkAsDone,
		decodeMarkAsDoneHTTPRequest,
		encodeMarkAsDoneHTTPResponse,
		options...,
	))

	return r
}

func encodeHTTPError(_ context.Context, err error, w http.ResponseWriter) {
	problem := problems.NewDetailedProblem(http.StatusInternalServerError, err.Error())

	w.Header().Set("Content-Type", problems.ProblemMediaType)
	_ = json.NewEncoder(w).Encode(problem)
}

func decodeCreateTodoHTTPRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req createTodoRequest

	var apiReq api.CreateTodoRequest
	err := json.NewDecoder(r.Body).Decode(&apiReq)
	if err != nil {
		return req, errors.Wrap(err, "failed to decode request")
	}

	req = createTodoRequest{
		Text: apiReq.Text,
	}

	return req, nil
}

func encodeCreateTodoHTTPResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	resp := response.(createTodoResponse)

	apiResp := api.CreateTodoResponse{
		Id: resp.ID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err := json.NewEncoder(w).Encode(apiResp)

	return errors.Wrap(err, "failed to send response")
}

func decodeListTodosHTTPRequest(_ context.Context, _ *http.Request) (interface{}, error) {
	return nil, nil
}

func encodeListTodosHTTPResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	resp := response.(listTodosResponse)

	apiResp := api.TodoList{
		Todos: make([]api.Todo, len(resp.Todos)),
	}

	for i, t := range resp.Todos {
		apiResp.Todos[i] = api.Todo{
			Id:   t.ID,
			Text: t.Text,
			Done: t.Done,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(apiResp)

	return errors.Wrap(err, "failed to send response")
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

func encodeMarkAsDoneHTTPResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	if f, ok := response.(endpoint.Failer); ok && f.Failed() != nil {
		errorEncoder(f.Failed(), w)

		return nil
	}

	w.WriteHeader(http.StatusOK)

	return nil
}

func errorEncoder(failed error, w http.ResponseWriter) {
	status := http.StatusInternalServerError

	// nolint: gocritic
	switch errors.Cause(failed).(type) {
	case todo.NotFoundError:
		status = http.StatusNotFound
	}

	problem := problems.NewDetailedProblem(status, failed.Error())

	w.Header().Set("Content-Type", problems.ProblemMediaType)
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(problem)
}

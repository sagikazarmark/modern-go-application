package tododriver

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	kitoc "github.com/go-kit/kit/tracing/opencensus"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/goph/emperror"
	"github.com/gorilla/mux"
	"github.com/moogar0880/problems"
	"github.com/pkg/errors"
	"github.com/sagikazarmark/ocmux"

	api "github.com/sagikazarmark/modern-go-application/.gen/api/openapi/todo/go"
	"github.com/sagikazarmark/modern-go-application/internal/todo"
)

// MakeHTTPHandler mounts all of the service endpoints into an http.Handler.
func MakeHTTPHandler(todoList TodoList, errorHandler todo.ErrorHandler) http.Handler {
	r := mux.NewRouter().PathPrefix("/todos").Subrouter()
	r.Use(ocmux.Middleware())

	e := MakeEndpoints(todoList)

	errorEncoder := httptransport.ServerErrorEncoder(func(ctx context.Context, err error, w http.ResponseWriter) {
		// This replaces server error log
		errorHandler.Handle(err)

		err = encodeHTTPError(err, w)
		if err != nil {
			errorHandler.Handle(err)
		}
	})

	r.Methods(http.MethodPost).Path("/").Handler(httptransport.NewServer(
		kitoc.TraceEndpoint("todo.CreateTodo")(e.Create),
		decodeCreateTodoHTTPRequest,
		encodeCreateTodoHTTPResponse,
		errorEncoder,
	))

	r.Methods(http.MethodGet).Path("/").Handler(httptransport.NewServer(
		kitoc.TraceEndpoint("todo.ListTodos")(e.List),
		decodeListTodosHTTPRequest,
		encodeListTodosHTTPResponse,
		errorEncoder,
	))

	r.Methods(http.MethodPost).Path("/{id}/done").Handler(httptransport.NewServer(
		kitoc.TraceEndpoint("todo.MarkAsDone")(e.MarkAsDone),
		decodeMarkAsDoneHTTPRequest,
		encodeMarkAsDoneHTTPResponse,
		errorEncoder,
	))

	return r
}

func encodeHTTPError(err error, w http.ResponseWriter) error {
	problem := problems.NewDetailedProblem(http.StatusInternalServerError, err.Error())

	w.Header().Set("Content-Type", problems.ProblemMediaType)
	err = json.NewEncoder(w).Encode(problem)

	return emperror.WrapWith(err, "failed to respond with error", "error", problem.Detail)
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
		return nil, emperror.With(errors.New("missing parameter from the URL"), "param", "id")
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

	if e, ok := failed.(*todoError); ok && e.Code() == codeNotFound {
		status = http.StatusNotFound
	}

	problem := problems.NewDetailedProblem(status, failed.Error())

	w.Header().Set("Content-Type", problems.ProblemMediaType)
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(problem)
}

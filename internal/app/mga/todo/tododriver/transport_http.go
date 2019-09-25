package tododriver

import (
	"context"
	"encoding/json"
	"net/http"

	"emperror.dev/errors"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/moogar0880/problems"

	api "github.com/sagikazarmark/modern-go-application/.gen/api/openapi/todo/go"
	"github.com/sagikazarmark/modern-go-application/internal/app/mga/todo"
	kitxhttp "github.com/sagikazarmark/modern-go-application/pkg/kitx/transport/http"
)

// RegisterHTTPHandlers mounts all of the service endpoints into a router.
func RegisterHTTPHandlers(endpoints Endpoints, factory kitxhttp.ServerFactory, router *mux.Router) {
	router.Methods(http.MethodPost).Path("").Handler(factory.NewServer(
		endpoints.CreateTodo,
		decodeCreateTodoHTTPRequest,
		encodeCreateTodoHTTPResponse,
	))

	router.Methods(http.MethodGet).Path("").Handler(factory.NewServer(
		endpoints.ListTodos,
		kithttp.NopRequestDecoder,
		encodeListTodosHTTPResponse,
	))

	router.Methods(http.MethodPost).Path("/{id}/done").Handler(factory.NewServer(
		endpoints.MarkAsDone,
		decodeMarkAsDoneHTTPRequest,
		kitxhttp.ErrorResponseEncoder(kitxhttp.NopResponseEncoder, errorEncoder),
	))
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

func errorEncoder(_ context.Context, w http.ResponseWriter, err error) error {
	status := http.StatusInternalServerError

	// nolint: gocritic
	switch errors.Cause(err).(type) {
	case todo.NotFoundError:
		status = http.StatusNotFound
	}

	problem := problems.NewDetailedProblem(status, err.Error())

	w.Header().Set("Content-Type", problems.ProblemMediaType)
	w.WriteHeader(status)
	e := json.NewEncoder(w).Encode(problem)

	return e
}

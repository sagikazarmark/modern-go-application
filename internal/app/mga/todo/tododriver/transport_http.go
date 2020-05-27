package tododriver

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"emperror.dev/errors"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	appkithttp "github.com/sagikazarmark/appkit/transport/http"
	kitxhttp "github.com/sagikazarmark/kitx/transport/http"

	api "github.com/sagikazarmark/modern-go-application/.gen/api/openapi/todo/go"
)

// RegisterHTTPHandlers mounts all of the service endpoints into a router.
func RegisterHTTPHandlers(endpoints Endpoints, router *mux.Router, options ...kithttp.ServerOption) {
	errorEncoder := kitxhttp.NewJSONProblemErrorResponseEncoder(appkithttp.NewDefaultProblemConverter())

	router.Methods(http.MethodPost).Path("").Handler(kithttp.NewServer(
		endpoints.CreateTodo,
		decodeCreateTodoHTTPRequest,
		kitxhttp.ErrorResponseEncoder(encodeCreateTodoHTTPResponse, errorEncoder),
		options...,
	))

	router.Methods(http.MethodGet).Path("").Handler(kithttp.NewServer(
		endpoints.ListTodos,
		kithttp.NopRequestDecoder,
		kitxhttp.ErrorResponseEncoder(encodeListTodosHTTPResponse, errorEncoder),
		options...,
	))

	router.Methods(http.MethodDelete).Path("").Handler(kithttp.NewServer(
		endpoints.DeleteAll,
		kithttp.NopRequestDecoder,
		kitxhttp.ErrorResponseEncoder(kitxhttp.StatusCodeResponseEncoder(http.StatusNoContent), errorEncoder),
		options...,
	))

	router.Methods(http.MethodPost).Path("/{id}/complete").Handler(kithttp.NewServer(
		endpoints.MarkAsComplete,
		decodeMarkAsCompleteHTTPRequest,
		kitxhttp.ErrorResponseEncoder(kitxhttp.NopResponseEncoder, errorEncoder),
		options...,
	))
}

func decodeCreateTodoHTTPRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var apiRequest api.CreateTodoRequest

	err := json.NewDecoder(r.Body).Decode(&apiRequest)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode request")
	}

	return CreateTodoRequest{
		Title: apiRequest.Title,
	}, nil
}

func encodeCreateTodoHTTPResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	resp := response.(CreateTodoResponse)

	host, _ := ctx.Value(kithttp.ContextKeyRequestHost).(string)
	path, _ := ctx.Value(kithttp.ContextKeyRequestPath).(string)

	apiResponse := api.Todo{
		Id:        resp.Todo.ID,
		Title:     resp.Todo.Title,
		Completed: resp.Todo.Completed,
		Url:       fmt.Sprintf("%s%s/%s", host, path, resp.Todo.ID),
	}

	return kitxhttp.JSONResponseEncoder(ctx, w, kitxhttp.WithStatusCode(apiResponse, http.StatusCreated))
}

func encodeListTodosHTTPResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	resp := response.(ListTodosResponse)

	host, _ := ctx.Value(kithttp.ContextKeyRequestHost).(string)
	path, _ := ctx.Value(kithttp.ContextKeyRequestPath).(string)

	todos := make([]api.Todo, 0, len(resp.Todos))

	for _, todo := range resp.Todos {
		todos = append(todos, api.Todo{
			Id:        todo.ID,
			Title:     todo.Title,
			Completed: todo.Completed,
			Url:       fmt.Sprintf("%s%s/%s", host, path, todo.ID),
		})
	}

	return kitxhttp.JSONResponseEncoder(ctx, w, todos)
}

func decodeMarkAsCompleteHTTPRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)

	id, ok := vars["id"]
	if !ok || id == "" {
		return nil, errors.NewWithDetails("missing parameter from the URL", "param", "id")
	}

	return MarkAsCompleteRequest{
		Id: id,
	}, nil
}

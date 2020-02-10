package tododriver

import (
	"context"
	"encoding/json"
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

	router.Methods(http.MethodPost).Path("/{id}/done").Handler(kithttp.NewServer(
		endpoints.MarkAsDone,
		decodeMarkAsDoneHTTPRequest,
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
		Text: apiRequest.Text,
	}, nil
}

func encodeCreateTodoHTTPResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	resp := response.(CreateTodoResponse)

	apiResponse := api.CreateTodoResponse{
		Id: resp.Id,
	}

	return kitxhttp.JSONResponseEncoder(ctx, w, kitxhttp.WithStatusCode(apiResponse, http.StatusCreated))
}

func encodeListTodosHTTPResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	resp := response.(ListTodosResponse)

	apiResponse := api.TodoList{}

	for _, todo := range resp.Todos {
		apiResponse.Todos = append(apiResponse.Todos, api.Todo{
			Id:   todo.ID,
			Text: todo.Text,
			Done: todo.Done,
		})
	}

	return kitxhttp.JSONResponseEncoder(ctx, w, apiResponse)
}

func decodeMarkAsDoneHTTPRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)

	id, ok := vars["id"]
	if !ok || id == "" {
		return nil, errors.NewWithDetails("missing parameter from the URL", "param", "id")
	}

	return MarkAsDoneRequest{
		Id: id,
	}, nil
}

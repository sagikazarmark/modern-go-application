package tododriver

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/pkg/errors"

	"github.com/sagikazarmark/modern-go-application/internal/todo"
)

// TodoList manages a list of todos.
type TodoList interface {
	// CreateTodo adds a new todo to the todo list.
	CreateTodo(ctx context.Context, req todo.CreateTodoRequest) (*todo.CreateTodoResponse, error)

	// ListTodos returns the list of todos on the list.
	ListTodos(ctx context.Context) (*todo.ListTodosResponse, error)

	// MarkAsDone marks a todo as done.
	MarkAsDone(ctx context.Context, req todo.MarkAsDoneRequest) error
}

// Endpoints collects all of the endpoints that compose a todo list service. It's
// meant to be used as a helper struct, to collect all of the endpoints into a
// single parameter.
type Endpoints struct {
	Create     endpoint.Endpoint
	List       endpoint.Endpoint
	MarkAsDone endpoint.Endpoint
}

// MakeEndpoints returns an Endpoints struct where each endpoint invokes
// the corresponding method on the provided service.
func MakeEndpoints(t TodoList) Endpoints {
	return Endpoints{
		Create:     MakeCreateEndpoint(t),
		List:       MakeListEndpoint(t),
		MarkAsDone: MakeMarkAsDoneEndpoint(t),
	}
}

type createTodoRequest struct {
	Text string
}

type createTodoResponse struct {
	ID string

	Err error
}

func (r createTodoResponse) Failed() error {
	return r.Err
}

// MakeCreateEndpoint returns an endpoint for the matching method of the underlying service.
func MakeCreateEndpoint(t TodoList) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		ereq := request.(createTodoRequest)

		req := todo.CreateTodoRequest{
			Text: ereq.Text,
		}

		resp, err := t.CreateTodo(ctx, req)

		return createTodoResponse{
			ID: resp.ID,
		}, err
	}
}

type todoListItem struct {
	ID   string
	Text string
	Done bool
}

type listTodosResponse struct {
	Todos []todoListItem

	Err error
}

func (r listTodosResponse) Failed() error {
	return r.Err
}

// MakeListEndpoint returns an endpoint for the matching method of the underlying service.
func MakeListEndpoint(t TodoList) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		resp, err := t.ListTodos(ctx)

		eresp := listTodosResponse{
			Todos: make([]todoListItem, len(resp.Todos)),
		}

		for i, t := range resp.Todos {
			eresp.Todos[i] = todoListItem(t)
		}

		return eresp, err
	}
}

type markAsDoneRequest struct {
	ID string
}

type markAsDoneResponse struct {
	Err error
}

func (r markAsDoneResponse) Failed() error {
	return r.Err
}

// MakeMarkAsDoneEndpoint returns an endpoint for the matching method of the underlying service.
func MakeMarkAsDoneEndpoint(t TodoList) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		ereq := request.(markAsDoneRequest)

		req := todo.MarkAsDoneRequest{
			ID: ereq.ID,
		}

		err = t.MarkAsDone(ctx, req)

		if _, ok := errors.Cause(err).(todo.TodoNotFoundError); ok {
			return markAsDoneResponse{
				Err: err,
			}, nil
		}

		return nil, err
	}
}

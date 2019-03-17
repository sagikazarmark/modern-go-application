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
	CreateTodo(ctx context.Context, text string) (string, error)

	// ListTodos returns the list of todos.
	ListTodos(ctx context.Context) ([]todo.Todo, error)

	// MarkAsDone marks a todo as done.
	MarkAsDone(ctx context.Context, id string) error
}

const (
	codeNotFound int = 1
)

type todoError struct {
	msg  string
	code int
}

func (e *todoError) Error() string {
	return e.msg
}

func (e *todoError) Code() int {
	return e.code
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
}

// MakeCreateEndpoint returns an endpoint for the matching method of the underlying service.
func MakeCreateEndpoint(t TodoList) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(createTodoRequest)

		id, err := t.CreateTodo(ctx, req.Text)

		return createTodoResponse{
			ID: id,
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
}

// MakeListEndpoint returns an endpoint for the matching method of the underlying service.
func MakeListEndpoint(t TodoList) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		todos, err := t.ListTodos(ctx)

		resp := listTodosResponse{
			Todos: make([]todoListItem, len(todos)),
		}

		for i, t := range todos {
			resp.Todos[i] = todoListItem(t)
		}

		return resp, err
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
		req := request.(markAsDoneRequest)

		err = t.MarkAsDone(ctx, req.ID)

		if _, ok := errors.Cause(err).(todo.NotFoundError); ok {
			return markAsDoneResponse{
				Err: &todoError{
					msg:  "todo not found",
					code: codeNotFound,
				},
			}, nil
		}

		return nil, err
	}
}

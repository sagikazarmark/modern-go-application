package tododriver

import (
	"context"

	"emperror.dev/errors"
	"github.com/go-kit/kit/endpoint"
	kitoc "github.com/go-kit/kit/tracing/opencensus"

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

type businessError interface {
	// IsBusinessError tells the transport layer whether this error should be translated into the transport format
	// or an internal error should be returned instead.
	IsBusinessError() bool
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
		Create:     kitoc.TraceEndpoint("todo.CreateTodo")(MakeCreateEndpoint(t)),
		List:       kitoc.TraceEndpoint("todo.ListTodos")(MakeListEndpoint(t)),
		MarkAsDone: kitoc.TraceEndpoint("todo.MarkAsDone")(MakeMarkAsDoneEndpoint(t)),
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
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(createTodoRequest)

		id, err := t.CreateTodo(ctx, req.Text)

		return createTodoResponse{
			ID: id,
		}, err
	}
}

type listTodosResponse struct {
	Todos []todo.Todo
}

// MakeListEndpoint returns an endpoint for the matching method of the underlying service.
func MakeListEndpoint(t TodoList) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		todos, err := t.ListTodos(ctx)

		resp := listTodosResponse{
			Todos: todos,
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
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(markAsDoneRequest)

		err := t.MarkAsDone(ctx, req.ID)

		if b, ok := errors.Cause(err).(businessError); ok && b.IsBusinessError() {
			return markAsDoneResponse{
				Err: err,
			}, nil
		}

		return nil, err
	}
}

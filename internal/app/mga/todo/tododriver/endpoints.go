package tododriver

import (
	"context"

	"emperror.dev/errors"
	"github.com/go-kit/kit/endpoint"
	kitoc "github.com/go-kit/kit/tracing/opencensus"
	kitxendpoint "github.com/sagikazarmark/kitx/endpoint"

	"github.com/sagikazarmark/modern-go-application/internal/app/mga/todo"
)

type businessError interface {
	// IsBusinessError tells the transport layer whether this error should be translated into the transport format
	// or an internal error should be returned instead.
	IsBusinessError() bool
}

// Endpoints collects all of the endpoints that compose a todo list service. It's
// meant to be used as a helper struct, to collect all of the endpoints into a
// single parameter.
type Endpoints struct {
	CreateTodo endpoint.Endpoint
	ListTodos  endpoint.Endpoint
	MarkAsDone endpoint.Endpoint
}

// MakeEndpoints returns an Endpoints struct where each endpoint invokes
// the corresponding method on the provided service.
func MakeEndpoints(service todo.Service, middleware ...endpoint.Middleware) Endpoints {
	mw := kitxendpoint.Chain(middleware...)
	return Endpoints{
		CreateTodo: mw(MakeCreateEndpoint(service)),
		ListTodos:  mw(MakeListEndpoint(service)),
		MarkAsDone: mw(MakeMarkAsDoneEndpoint(service)),
	}
}

// TraceEndpoints returns an Endpoints struct where each endpoint invokes
// the corresponding method on the provided service.
func TraceEndpoints(endpoints Endpoints) Endpoints {
	return Endpoints{
		CreateTodo: kitoc.TraceEndpoint("todo.CreateTodo")(endpoints.CreateTodo),
		ListTodos:  kitoc.TraceEndpoint("todo.ListTodos")(endpoints.ListTodos),
		MarkAsDone: kitoc.TraceEndpoint("todo.MarkAsDone")(endpoints.MarkAsDone),
	}
}

type createTodoRequest struct {
	Text string
}

type createTodoResponse struct {
	ID string
}

// MakeCreateEndpoint returns an endpoint for the matching method of the underlying service.
func MakeCreateEndpoint(service todo.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(createTodoRequest)

		id, err := service.CreateTodo(ctx, req.Text)

		return createTodoResponse{
			ID: id,
		}, err
	}
}

type listTodosResponse struct {
	Todos []todo.Todo
}

// MakeListEndpoint returns an endpoint for the matching method of the underlying service.
func MakeListEndpoint(service todo.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		todos, err := service.ListTodos(ctx)

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
func MakeMarkAsDoneEndpoint(service todo.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(markAsDoneRequest)

		err := service.MarkAsDone(ctx, req.ID)

		var berr businessError
		if errors.As(err, &berr) && berr.IsBusinessError() {
			return markAsDoneResponse{
				Err: err,
			}, nil
		}

		return nil, err
	}
}

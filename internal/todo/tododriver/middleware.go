package tododriver

import (
	"context"
	"time"

	"go.opencensus.io/trace"

	"github.com/sagikazarmark/modern-go-application/internal/todo"
)

// Middleware describes a service middleware.
type Middleware func(TodoList) TodoList

// LoggingMiddleware is a service level logging middleware for TodoList.
func LoggingMiddleware(logger todo.Logger) Middleware {
	return func(next TodoList) TodoList {
		return &loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

type loggingMiddleware struct {
	next   TodoList
	logger todo.Logger
}

func (mw *loggingMiddleware) CreateTodo(ctx context.Context, text string) (string, error) {
	logger := mw.logger.WithContext(ctx)

	logger.Trace("processing request", map[string]interface{}{
		"operation": "todo.CreateTodo",
	})

	defer func(begin time.Time) {
		logger.Trace("processing request finished", map[string]interface{}{
			"operation": "todo.CreateTodo",
			"took":      time.Since(begin),
		})
	}(time.Now())

	return mw.next.CreateTodo(ctx, text)
}

func (mw *loggingMiddleware) ListTodos(ctx context.Context) ([]todo.Todo, error) {
	logger := mw.logger.WithContext(ctx)

	logger.Trace("processing request", map[string]interface{}{
		"operation": "todo.ListTodos",
	})

	defer func(begin time.Time) {
		logger.Trace("processing request finished", map[string]interface{}{
			"operation": "todo.ListTodos",
			"took":      time.Since(begin),
		})
	}(time.Now())

	return mw.next.ListTodos(ctx)
}

func (mw *loggingMiddleware) MarkAsDone(ctx context.Context, id string) error {
	logger := mw.logger.WithContext(ctx)

	logger.Trace("processing request", map[string]interface{}{
		"operation": "todo.MarkAsDone",
	})

	defer func(begin time.Time) {
		logger.Trace("processing request finished", map[string]interface{}{
			"operation": "todo.MarkAsDone",
			"took":      time.Since(begin),
		})
	}(time.Now())

	return mw.next.MarkAsDone(ctx, id)
}

// TracingMiddleware is a service level tracing middleware for TodoList.
func TracingMiddleware() Middleware {
	return func(next TodoList) TodoList {
		return &tracingMiddleware{
			next: next,
		}
	}
}

type tracingMiddleware struct {
	next TodoList
}

func (mw *tracingMiddleware) CreateTodo(ctx context.Context, text string) (string, error) {
	id, err := mw.next.CreateTodo(ctx, text)

	if span := trace.FromContext(ctx); span != nil {
		span.AddAttributes(trace.StringAttribute("todo_id", id))
	}

	return id, err
}

func (mw *tracingMiddleware) ListTodos(ctx context.Context) ([]todo.Todo, error) {
	return mw.next.ListTodos(ctx)
}

func (mw *tracingMiddleware) MarkAsDone(ctx context.Context, id string) error {
	if span := trace.FromContext(ctx); span != nil {
		span.AddAttributes(trace.StringAttribute("todo_id", id))
	}

	return mw.next.MarkAsDone(ctx, id)
}

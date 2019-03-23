package tododriver

import (
	"context"
	"time"

	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
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

// InstrumentationMiddleware is a service level tracing middleware for TodoList.
func InstrumentationMiddleware() Middleware {
	return func(next TodoList) TodoList {
		return &instrumentationMiddleware{
			next: next,
		}
	}
}

// Todo business metrics
// nolint: gochecknoglobals
var (
	CreatedTodoCount = stats.Int64("created_todo_count", "Number of TODOs created", stats.UnitDimensionless)
	DoneTodoCount    = stats.Int64("done_todo_count", "Number of TODOs marked done", stats.UnitDimensionless)
)

// nolint: gochecknoglobals
var (
	CreatedTodoCountView = &view.View{
		Name:        "todo_created_count",
		Description: "Count of TODOs created",
		Measure:     CreatedTodoCount,
		Aggregation: view.Count(),
	}

	DoneTodoCountView = &view.View{
		Name:        "todo_done_count",
		Description: "Count of TODOs done",
		Measure:     DoneTodoCount,
		Aggregation: view.Count(),
	}
)

type instrumentationMiddleware struct {
	next TodoList
}

func (mw *instrumentationMiddleware) CreateTodo(ctx context.Context, text string) (string, error) {
	id, err := mw.next.CreateTodo(ctx, text)

	if span := trace.FromContext(ctx); span != nil {
		span.AddAttributes(trace.StringAttribute("todo_id", id))
	}

	stats.Record(ctx, CreatedTodoCount.M(1))

	return id, err
}

func (mw *instrumentationMiddleware) ListTodos(ctx context.Context) ([]todo.Todo, error) {
	return mw.next.ListTodos(ctx)
}

func (mw *instrumentationMiddleware) MarkAsDone(ctx context.Context, id string) error {
	if span := trace.FromContext(ctx); span != nil {
		span.AddAttributes(trace.StringAttribute("todo_id", id))
	}

	stats.Record(ctx, DoneTodoCount.M(1))

	return mw.next.MarkAsDone(ctx, id)
}

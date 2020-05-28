package tododriver

import (
	"context"

	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"

	"github.com/sagikazarmark/modern-go-application/internal/app/mga/todo"
)

// Middleware is a service middleware.
type Middleware func(todo.Service) todo.Service

// defaultMiddleware helps implementing partial middleware.
type defaultMiddleware struct {
	service todo.Service
}

func (m defaultMiddleware) AddItem(ctx context.Context, newItem todo.NewItem) (todo.Item, error) {
	return m.service.AddItem(ctx, newItem)
}

func (m defaultMiddleware) ListItems(ctx context.Context) ([]todo.Item, error) {
	return m.service.ListItems(ctx)
}

func (m defaultMiddleware) DeleteItems(ctx context.Context) error {
	return m.service.DeleteItems(ctx)
}

func (m defaultMiddleware) GetItem(ctx context.Context, id string) (todo.Item, error) {
	return m.service.GetItem(ctx, id)
}

func (m defaultMiddleware) UpdateItem(ctx context.Context, id string, itemUpdate todo.ItemUpdate) (todo.Item, error) {
	return m.service.UpdateItem(ctx, id, itemUpdate)
}

func (m defaultMiddleware) DeleteItem(ctx context.Context, id string) error {
	return m.service.DeleteItem(ctx, id)
}

func (m defaultMiddleware) MarkAsComplete(ctx context.Context, id string) error {
	return m.service.MarkAsComplete(ctx, id)
}

// LoggingMiddleware is a service level logging middleware.
func LoggingMiddleware(logger todo.Logger) Middleware {
	return func(next todo.Service) todo.Service {
		return loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

type loggingMiddleware struct {
	next   todo.Service
	logger todo.Logger
}

func (mw loggingMiddleware) AddItem(ctx context.Context, newItem todo.NewItem) (todo.Item, error) {
	logger := mw.logger.WithContext(ctx)

	logger.Info("adding item")

	id, err := mw.next.AddItem(ctx, newItem)
	if err != nil {
		return id, err
	}

	logger.Info("added item", map[string]interface{}{"item_id": id})

	return id, err
}

func (mw loggingMiddleware) ListItems(ctx context.Context) ([]todo.Item, error) {
	logger := mw.logger.WithContext(ctx)

	logger.Info("listing item")

	return mw.next.ListItems(ctx)
}

func (mw loggingMiddleware) DeleteItems(ctx context.Context) error {
	logger := mw.logger.WithContext(ctx)

	logger.Info("deleting all items")

	return mw.next.DeleteItems(ctx)
}

func (mw loggingMiddleware) GetItem(ctx context.Context, id string) (todo.Item, error) {
	logger := mw.logger.WithContext(ctx)

	logger.Info("getting item details", map[string]interface{}{"item_id": id})

	return mw.next.GetItem(ctx, id)
}

func (mw loggingMiddleware) UpdateItem(ctx context.Context, id string, itemUpdate todo.ItemUpdate) (todo.Item, error) { // nolint: lll
	logger := mw.logger.WithContext(ctx)

	logger.Info("updating item", map[string]interface{}{"item_id": id})

	return mw.next.UpdateItem(ctx, id, itemUpdate)
}

func (mw loggingMiddleware) DeleteItem(ctx context.Context, id string) error {
	logger := mw.logger.WithContext(ctx)

	logger.Info("deleting item", map[string]interface{}{"item_id": id})

	return mw.next.DeleteItem(ctx, id)
}

func (mw loggingMiddleware) MarkAsComplete(ctx context.Context, id string) error {
	logger := mw.logger.WithContext(ctx)

	logger.Info("marking item as complete", map[string]interface{}{"item_id": id})

	return mw.next.MarkAsComplete(ctx, id)
}

// Business metrics
// nolint: gochecknoglobals,lll
var (
	CreatedTodoItemCount  = stats.Int64("created_todo_item_count", "Number of todo items created", stats.UnitDimensionless)
	CompleteTodoItemCount = stats.Int64("complete_todo_item_count", "Number of todo items marked complete", stats.UnitDimensionless)
)

// nolint: gochecknoglobals
var (
	CreatedTodoItemCountView = &view.View{
		Name:        "todo_item_created_count",
		Description: "Count of todo items created",
		Measure:     CreatedTodoItemCount,
		Aggregation: view.Count(),
	}

	CompleteTodoItemCountView = &view.View{
		Name:        "todo_item_complete_count",
		Description: "Count of todo items complete",
		Measure:     CompleteTodoItemCount,
		Aggregation: view.Count(),
	}
)

// InstrumentationMiddleware is a service level instrumentation middleware.
func InstrumentationMiddleware() Middleware {
	return func(next todo.Service) todo.Service {
		return instrumentationMiddleware{
			Service: defaultMiddleware{next},
			next:    next,
		}
	}
}

type instrumentationMiddleware struct {
	todo.Service
	next todo.Service
}

func (mw instrumentationMiddleware) AddItem(ctx context.Context, newItem todo.NewItem) (todo.Item, error) {
	item, err := mw.next.AddItem(ctx, newItem)
	if err != nil {
		return item, err
	}

	if span := trace.FromContext(ctx); span != nil {
		span.AddAttributes(trace.StringAttribute("item_id", item.ID))
	}

	stats.Record(ctx, CreatedTodoItemCount.M(1))

	return item, nil
}

func (mw instrumentationMiddleware) MarkAsComplete(ctx context.Context, id string) error {
	if span := trace.FromContext(ctx); span != nil {
		span.AddAttributes(trace.StringAttribute("item_id", id))
	}

	stats.Record(ctx, CompleteTodoItemCount.M(1))

	return mw.next.MarkAsComplete(ctx, id)
}

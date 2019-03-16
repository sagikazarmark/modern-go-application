package todoadapter

import (
	"context"

	"github.com/pkg/errors"

	"github.com/sagikazarmark/modern-go-application/internal/todo"
)

// EventBus dispatches events to event handlers.
type EventBus interface {
	// Publish sends an event to the event bus.
	Publish(ctx context.Context, event interface{}) error
}

// TodoEvents dispatches todo events.
type TodoEvents struct {
	eventBus EventBus
}

// NewTodoEvents returns a new TodoEvents instance.
func NewTodoEvents(eventBus EventBus) *TodoEvents {
	return &TodoEvents{
		eventBus: eventBus,
	}
}

// MarkedAsDone dispatches a MarkedAsDone event.
func (e *TodoEvents) MarkedAsDone(ctx context.Context, event todo.MarkedAsDone) error {
	err := e.eventBus.Publish(ctx, event)
	if err != nil {
		return errors.WithMessage(err, "failed to dispatch event")
	}

	return nil
}

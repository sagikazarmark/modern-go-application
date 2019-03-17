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

// EventDispatcher dispatches todo events.
type EventDispatcher struct {
	eventBus EventBus
}

// NewEventDispatcher returns a new EventDispatcher instance.
func NewEventDispatcher(eventBus EventBus) *EventDispatcher {
	return &EventDispatcher{
		eventBus: eventBus,
	}
}

// MarkedAsDone dispatches a MarkedAsDone event.
func (e *EventDispatcher) MarkedAsDone(ctx context.Context, event todo.MarkedAsDone) error {
	err := e.eventBus.Publish(ctx, event)
	if err != nil {
		return errors.WithMessage(err, "failed to dispatch event")
	}

	return nil
}

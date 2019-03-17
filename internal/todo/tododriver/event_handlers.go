package tododriver

import (
	"context"
	"fmt"

	"github.com/goph/emperror"
	"github.com/pkg/errors"

	"github.com/sagikazarmark/modern-go-application/internal/todo"
)

// MarkedAsDoneHandler handles MarkedAsDone events.
type MarkedAsDoneHandler interface {
	// MarkedAsDone handles a MarkedAsDone event.
	MarkedAsDone(ctx context.Context, event todo.MarkedAsDone) error
}

// MarkedAsDoneEventHandler handles a MarkedAsDone events.
type MarkedAsDoneEventHandler struct {
	handler MarkedAsDoneHandler
}

// NewMarkedAsDoneEventHandler returns a new MarkedAsDoneEventHandler instance.
func NewMarkedAsDoneEventHandler(handler MarkedAsDoneHandler) *MarkedAsDoneEventHandler {
	return &MarkedAsDoneEventHandler{
		handler: handler,
	}
}

func (*MarkedAsDoneEventHandler) NewEvent() interface{} {
	return &todo.MarkedAsDone{}
}

func (h *MarkedAsDoneEventHandler) Handle(ctx context.Context, event interface{}) error {
	e, ok := event.(*todo.MarkedAsDone)
	if !ok {
		return emperror.With(
			errors.New("unexpected event type"),
			"type", fmt.Sprintf("%T", event),
		)
	}

	return h.handler.MarkedAsDone(ctx, *e)
}

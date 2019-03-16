package tododriver

import (
	"context"
	"fmt"

	"github.com/goph/emperror"
	"github.com/pkg/errors"

	"github.com/sagikazarmark/modern-go-application/internal/todo"
)

// MarkedAsDoneTodoEventHandler handles MarkedAsDone events.
type MarkedAsDoneTodoEventHandler interface {
	// MarkedAsDone handles a MarkedAsDone event.
	MarkedAsDone(ctx context.Context, event todo.MarkedAsDone) error
}

// MarkedAsDoneEventHandler handles a MarkedAsDone events.
type MarkedAsDoneEventHandler struct {
	handler MarkedAsDoneTodoEventHandler
}

// NewMarkedAsDoneEventHandler returns a new MarkedAsDoneEventHandler instance.
func NewMarkedAsDoneEventHandler(handler MarkedAsDoneTodoEventHandler) *MarkedAsDoneEventHandler {
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

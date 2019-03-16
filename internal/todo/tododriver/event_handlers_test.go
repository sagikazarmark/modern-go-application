package tododriver

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sagikazarmark/modern-go-application/internal/todo"
)

type markedAsDoneTodoEventHandlerStub struct {
	ctx   context.Context
	event todo.MarkedAsDone
}

func (s *markedAsDoneTodoEventHandlerStub) MarkedAsDone(ctx context.Context, event todo.MarkedAsDone) error {
	s.ctx = ctx
	s.event = event

	return nil
}

func TestMarkedAsDoneEventHandler_NewEvent(t *testing.T) {
	handler := NewMarkedAsDoneEventHandler(&markedAsDoneTodoEventHandlerStub{})

	event := handler.NewEvent()

	assert.IsType(t, &todo.MarkedAsDone{}, event)
}

func TestMarkedAsDoneEventHandler_Handle(t *testing.T) {
	h := &markedAsDoneTodoEventHandlerStub{}
	handler := NewMarkedAsDoneEventHandler(h)

	ctx := context.Background()
	event := todo.MarkedAsDone{
		ID: "1234",
	}

	err := handler.Handle(ctx, &event)
	require.NoError(t, err)

	assert.Equal(t, h.ctx, ctx)
	assert.Equal(t, h.event, event)
}

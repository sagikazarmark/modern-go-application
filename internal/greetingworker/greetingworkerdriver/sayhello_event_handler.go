package greetingworkerdriver

import (
	"context"
	"encoding/json"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/pkg/errors"

	"github.com/sagikazarmark/modern-go-application/internal/greetingworker"
)

// GreeterEventSubscriber subscribes to hello world events.
type GreeterEventSubscriber interface {
	// SaidHello handles a SaidHello event.
	SaidHelloTo(ctx context.Context, event greetingworker.SaidHello) error
}

type GreeterEventHandler struct {
	subscriber GreeterEventSubscriber
}

// NewGreeterEventHandler returns a new GreeterEventHandler.
func NewGreeterEventHandler(subscriber GreeterEventSubscriber) *GreeterEventHandler {
	return &GreeterEventHandler{
		subscriber: subscriber,
	}
}

func (h *GreeterEventHandler) SaidHelloTo(msg *message.Message) (messages []*message.Message, e error) {
	// TODO: get correlation ID from message and add it to the context
	//correlationID := middleware.MessageCorrelationID(msg)

	var event greetingworker.SaidHello

	err := json.Unmarshal(msg.Payload, &event)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal event payload")
	}

	err = h.subscriber.SaidHelloTo(context.Background(), event)

	return nil, err
}

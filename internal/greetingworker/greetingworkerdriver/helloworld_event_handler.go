package greetingworkerdriver

import (
	"context"
	"encoding/json"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/pkg/errors"
	"github.com/sagikazarmark/modern-go-application/internal/greetingworker"
)

// HelloWorldEventSubscriber subscribes to hello world events.
type HelloWorldEventSubscriber interface {
	// SaidHello handles a SaidHello event.
	SaidHello(ctx context.Context, event greetingworker.SaidHello) error
}

type HelloWorldEventHandler struct {
	subscriber HelloWorldEventSubscriber
}

// NewHelloWorldEventHandler returns a new HelloWorldEventHandler.
func NewHelloWorldEventHandler(subscriber HelloWorldEventSubscriber) *HelloWorldEventHandler {
	return &HelloWorldEventHandler{
		subscriber: subscriber,
	}
}

func (h *HelloWorldEventHandler) SaidHello(msg *message.Message) (messages []*message.Message, e error) {
	// TODO: get correlation ID from message and add it to the context
	//correlationID := middleware.MessageCorrelationID(msg)

	var event greetingworker.SaidHello

	err := json.Unmarshal(msg.Payload, &event)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal event payload")
	}

	err = h.subscriber.SaidHello(context.Background(), event)

	return nil, err
}

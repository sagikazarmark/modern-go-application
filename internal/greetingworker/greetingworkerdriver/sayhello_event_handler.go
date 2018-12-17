package greetingworkerdriver

import (
	"context"
	"encoding/json"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/pkg/errors"
	"github.com/sagikazarmark/modern-go-application/internal/greetingworker"
)

// SayHelloEventSubscriber subscribes to hello world events.
type SayHelloEventSubscriber interface {
	// SaidHelloTo handles a SaidHelloTo event.
	SaidHelloTo(ctx context.Context, event greetingworker.SaidHelloTo) error
}

type SayHelloEventHandler struct {
	subscriber SayHelloEventSubscriber
}

// NewSayHelloEventHandler returns a new SayHelloEventHandler.
func NewSayHelloEventHandler(subscriber SayHelloEventSubscriber) *SayHelloEventHandler {
	return &SayHelloEventHandler{
		subscriber: subscriber,
	}
}

func (h *SayHelloEventHandler) SaidHelloTo(msg *message.Message) (messages []*message.Message, e error) {
	// TODO: get correlation ID from message and add it to the context
	//correlationID := middleware.MessageCorrelationID(msg)

	var event greetingworker.SaidHelloTo

	err := json.Unmarshal(msg.Payload, &event)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal event payload")
	}

	err = h.subscriber.SaidHelloTo(context.Background(), event)

	return nil, err
}

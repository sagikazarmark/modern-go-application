// nolint: dupl
package greetingadapter

import (
	"context"
	"encoding/json"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"github.com/sagikazarmark/modern-go-application/internal/greeting"
)

const (
	saidHelloTopic = "said_hello"
)

// HelloWorldEvents is the dispatcher for hello world events.
type HelloWorldEvents struct {
	publisher message.Publisher
}

// NewHelloWorldEvents returns a new HelloWorldEvents instance.
func NewHelloWorldEvents(publisher message.Publisher) *HelloWorldEvents {
	return &HelloWorldEvents{
		publisher: publisher,
	}
}

// SaidHello dispatches a SaidHello event.
func (e *HelloWorldEvents) SaidHello(ctx context.Context, event greeting.SaidHello) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return errors.Wrap(err, "failed to marshal message payload")
	}

	msgID, err := uuid.NewV4()
	if err != nil {
		return errors.Wrap(err, "failed to generate message ID")
	}
	msg := message.NewMessage(msgID.String(), payload)

	// TODO: set from context
	corrID, err := uuid.NewV4()
	if err != nil {
		return errors.Wrap(err, "failed to generate correlation ID")
	}
	middleware.SetCorrelationID(corrID.String(), msg)

	err = e.publisher.Publish(saidHelloTopic, msg)
	if err != nil {
		return errors.WithMessage(err, "failed to publish message")
	}

	return nil
}

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
	saidHelloToTopic = "said_hello_to"
)

// SayHelloEvents is the dispatcher for hello world events.
type SayHelloEvents struct {
	publisher message.Publisher
}

// NewSayHelloEvents returns a new SayHelloEvents instance.
func NewSayHelloEvents(publisher message.Publisher) *SayHelloEvents {
	return &SayHelloEvents{
		publisher: publisher,
	}
}

// SaidHelloTo dispatches a SaidHelloTo event.
func (e *SayHelloEvents) SaidHelloTo(ctx context.Context, event greeting.SaidHelloTo) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return errors.Wrap(err, "failed to marshal event payload")
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

	err = e.publisher.Publish(saidHelloToTopic, msg)
	if err != nil {
		return errors.WithMessage(err, "failed to publish event")
	}

	return nil
}

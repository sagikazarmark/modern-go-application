package watermill

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	watermilllog "logur.dev/integration/watermill"
	"logur.dev/logur"
)

// NewPubSub returns a new PubSub.
func NewPubSub(logger logur.Logger) (message.Publisher, message.Subscriber) {
	pubsub := gochannel.NewGoChannel(
		gochannel.Config{},
		watermilllog.New(logur.WithField(logger, "component", "watermill")),
	)

	return pubsub, pubsub
}

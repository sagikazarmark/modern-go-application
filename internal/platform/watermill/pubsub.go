package watermill

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/infrastructure/gochannel"
	"github.com/goph/logur"
	"github.com/goph/logur/integrations/watermilllog"
)

// NewPubSub returns a new PubSub.
func NewPubSub(logger logur.Logger) message.PubSub {
	return gochannel.NewGoChannel(
		gochannel.Config{},
		watermilllog.New(logur.WithFields(logger, map[string]interface{}{"component": "watermill"})),
	)
}

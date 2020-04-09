package watermill

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/sagikazarmark/kitx/correlation"
)

// PublisherCorrelationID decorates a publisher with a correlation ID middleware.
func PublisherCorrelationID(publisher message.Publisher) message.Publisher {
	publisher, _ = message.MessageTransformPublisherDecorator(func(msg *message.Message) {
		if cid, ok := correlation.FromContext(msg.Context()); ok {
			middleware.SetCorrelationID(cid, msg)
		}
	})(publisher)

	return publisher
}

// SubscriberCorrelationID decorates a subscriber with a correlation ID middleware.
func SubscriberCorrelationID(subscriber message.Subscriber) message.Subscriber {
	subscriber, _ = message.MessageTransformSubscriberDecorator(func(msg *message.Message) {
		if cid := middleware.MessageCorrelationID(msg); cid != "" {
			msg.SetContext(correlation.ToContext(msg.Context(), cid))
		}
	})(subscriber)

	return subscriber
}

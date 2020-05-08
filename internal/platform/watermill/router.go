package watermill

import (
	"time"

	"emperror.dev/errors"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	watermilllog "logur.dev/integration/watermill"
	"logur.dev/logur"
)

// NewRouter returns a new message router for message subscription logic.
func NewRouter(logger logur.Logger) (*message.Router, error) {
	h, err := message.NewRouter(
		message.RouterConfig{},
		watermilllog.New(logur.WithField(logger, "component", "watermill")),
	)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to create message router")
	}

	retryMiddleware := middleware.Retry{}
	retryMiddleware.MaxRetries = 1
	retryMiddleware.MaxInterval = time.Millisecond * 10

	h.AddMiddleware(
		// if retries limit was exceeded, message is sent to poison queue (poison_queue topic)
		retryMiddleware.Middleware,

		// recovered recovers panic from handlers
		middleware.Recoverer,

		// correlation ID middleware adds to every produced message correlation id of consumed message,
		// useful for debugging
		middleware.CorrelationID,
	)

	return h, nil
}

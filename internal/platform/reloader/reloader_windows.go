package reloader

import (
	"context"
	"net"

	"github.com/oklog/run"
	"logur.dev/logur"
)

type UnsupportedReloader struct {
}

func Create(logger logur.LoggerFacade) Reloader {
	logger.Warn("graceful reload is not supported on this platform")

	return &UnsupportedReloader{}
}

func (t *UnsupportedReloader) Listen(network, address string) (net.Listener, error) {
	return net.Listen(network, address)
}

func (t *UnsupportedReloader) SetupGracefulRestart(context context.Context, group run.Group) {
	// no-op since it isn't supported
}

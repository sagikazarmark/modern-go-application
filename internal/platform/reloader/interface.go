package reloader

import (
	"context"
	"net"

	"github.com/oklog/run"
)

type Reloader interface {
	Listen(network, address string) (net.Listener, error)
	SetupGracefulRestart(context.Context, run.Group)
}

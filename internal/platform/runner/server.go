package runner

import (
	"context"
	"net"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/goph/emperror"
)

type server interface {
	Serve(l net.Listener) error
	Shutdown(ctx context.Context) error
	Close() error
}

// ServerRunner configures server run funcs.
type ServerRunner struct {
	shutdownTimeout time.Duration

	logger       log.Logger
	errorHandler emperror.Handler
}

// NewServerRunner returns a new ServerRunner instance.
func NewServerRunner(shutdownTimeout time.Duration, logger log.Logger, errorHandler emperror.Handler) *ServerRunner {
	return &ServerRunner{
		shutdownTimeout: shutdownTimeout,

		logger:       logger,
		errorHandler: errorHandler,
	}
}

func (r *ServerRunner) RunFuncs(s server, ln net.Listener, name string) (func() error, func(e error)) {
	return func() error {
		level.Info(r.logger).Log("msg", "starting server", "name", name)

		return s.Serve(ln)
	},
		func(e error) {
			ctx, cancel := context.WithTimeout(context.Background(), r.shutdownTimeout)
			defer cancel()

			level.Info(r.logger).Log("msg", "shutting server down", "name", name)

			err := s.Shutdown(ctx)
			if err != nil {
				r.errorHandler.Handle(emperror.With(err, "name", name))
			}

			s.Close()
		}
}

// +build !windows

package reloader

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/cloudflare/tableflip"
	"github.com/oklog/run"
	appkitrun "github.com/sagikazarmark/appkit/run"
	"logur.dev/logur"
)

type TableflipReloader struct {
	*tableflip.Upgrader
}

func Create(logger logur.LoggerFacade) *Reloader {
	upg, _ := tableflip.New(tableflip.Options{})

	// Do an upgrade on SIGHUP
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGHUP)
		for range ch {
			logger.Info("graceful reloading")

			_ = upg.Upgrade()
		}
	}()

	return &TableflipReloader{upg}
}

func (t *TableflipReloader) SetupGracefulRestart(context context.Context, group run.Group) {
	group.Add(appkitrun.GracefulRestart(context, t))
}

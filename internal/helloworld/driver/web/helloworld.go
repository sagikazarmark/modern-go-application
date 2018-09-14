package web

import (
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/goph/emperror"
)

// HelloWorldDriverOption configures a HelloWorldDriver.
type HelloWorldDriverOption interface {
	apply(d *HelloWorldDriver)
}

// HelloWorldDriverOptionFunc makes a function a HelloWorldDriverOption if it's signature is compatible.
type HelloWorldDriverOptionFunc func(d *HelloWorldDriver)

func (fn HelloWorldDriverOptionFunc) apply(d *HelloWorldDriver) {
	fn(d)
}

// Logger configures a logger instance in a HelloWorldDriver.
func Logger(l log.Logger) HelloWorldDriverOption {
	return HelloWorldDriverOptionFunc(func(d *HelloWorldDriver) {
		d.logger = l
	})
}

// ErrorHandler configures an error handler instance in a HelloWorldDriver.
func ErrorHandler(h emperror.Handler) HelloWorldDriverOption {
	return HelloWorldDriverOptionFunc(func(d *HelloWorldDriver) {
		d.errorHandler = h
	})
}

type helloWorlder interface {
	HelloWorld() string
}

// HelloWorldDriver exposes the UseCase on an HTTP interface.
type HelloWorldDriver struct {
	helloWorld helloWorlder

	logger       log.Logger
	errorHandler emperror.Handler
}

// NewHelloWorldDriver returns a new HelloWorldDriver instance.
func NewHelloWorldDriver(helloWorld helloWorlder, opts ...HelloWorldDriverOption) *HelloWorldDriver {
	d := &HelloWorldDriver{
		helloWorld: helloWorld,
	}

	for _, opt := range opts {
		opt.apply(d)
	}

	// Default logger instance
	if d.logger == nil {
		d.logger = log.NewNopLogger()
	}

	// Default error handler instance
	if d.errorHandler == nil {
		d.errorHandler = emperror.NewNopHandler()
	}

	return d
}

func (d *HelloWorldDriver) HelloWorld(rw http.ResponseWriter, r *http.Request) {
	level.Info(d.logger).Log("msg", "Hello, World!")

	_, err := rw.Write([]byte(d.helloWorld.HelloWorld()))
	if err != nil {
		d.errorHandler.Handle(err)
	}
}

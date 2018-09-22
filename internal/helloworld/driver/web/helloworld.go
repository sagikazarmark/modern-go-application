package web

import (
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/goph/emperror"
	"github.com/sagikazarmark/go-service-project-boilerplate/.gen/openapi/go"
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
	SayHello(who string) string
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

func (d *HelloWorldDriver) HelloWorld(w http.ResponseWriter, r *http.Request) {
	level.Info(d.logger).Log("msg", "Hello, World!")

	response := api.Hello{
		Message: d.helloWorld.HelloWorld(),
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	encoder := json.NewEncoder(w)

	if err := encoder.Encode(response); err != nil {
		d.errorHandler.Handle(err)
	}
}

func (d *HelloWorldDriver) SayHello(w http.ResponseWriter, r *http.Request) {
	level.Info(d.logger).Log("msg", "Say hello")

	decoder := json.NewDecoder(r.Body)

	var request api.HelloRequest

	if err := decoder.Decode(&request); err != nil {
		d.errorHandler.Handle(err)

		http.Error(w, "invalid request", http.StatusBadRequest)

		return
	}

	response := api.Hello{
		Message: d.helloWorld.SayHello(request.Who),
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	encoder := json.NewEncoder(w)

	if err := encoder.Encode(response); err != nil {
		d.errorHandler.Handle(err)
	}
}

package internal

import (
	"net/http"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/goph/emperror"
	"github.com/goph/logur"
	"github.com/gorilla/mux"
	"github.com/mccutchen/go-httpbin/httpbin"
	"google.golang.org/grpc"

	greetingpb "github.com/sagikazarmark/modern-go-application/.gen/proto/greeting"
	"github.com/sagikazarmark/modern-go-application/internal/greeting"
	"github.com/sagikazarmark/modern-go-application/internal/greeting/greetingadapter"
	"github.com/sagikazarmark/modern-go-application/internal/greeting/greetingdriver"
	"github.com/sagikazarmark/modern-go-application/internal/greetingworker"
	"github.com/sagikazarmark/modern-go-application/internal/greetingworker/greetingworkeradapter"
	"github.com/sagikazarmark/modern-go-application/internal/greetingworker/greetingworkerdriver"
	"github.com/sagikazarmark/modern-go-application/internal/landing/landingdriver"
)

// NewApp returns a new HTTP and a new gRPC application.
func NewApp(
	logger logur.Logger,
	publisher message.Publisher,
	errorHandler emperror.Handler,
) (http.Handler, func(*grpc.Server)) {
	greeter := greeting.NewGreeter(
		greetingadapter.NewGreeterEvents(publisher),
		greetingadapter.NewLogger(logger),
		errorHandler,
	)

	router := mux.NewRouter()

	router.Path("/").Methods("GET").Handler(landingdriver.NewHTTPHandler())
	router.PathPrefix("/greeting").Methods("POST").Handler(
		http.StripPrefix("/greeting", greetingdriver.NewHTTPHandler(greeter, errorHandler)),
	)
	router.PathPrefix("/httpbin").Handler(
		http.StripPrefix(
			"/httpbin",
			httpbin.New(
				httpbin.WithObserver(func(result httpbin.Result) {
					logger.Info(
						"httpbin call",
						map[string]interface{}{
							"status":      result.Status,
							"method":      result.Method,
							"uri":         result.URI,
							"size_bytes":  result.Size,
							"duration_ms": result.Duration.Seconds() * 1e3, // https://github.com/golang/go/issues/5491#issuecomment-66079585
						},
					)
				}),
			).Handler(),
		),
	)

	helloWorldGRPCController := greetingdriver.NewGRPCController(greeter, errorHandler)

	return router, func(s *grpc.Server) {
		greetingpb.RegisterGreeterServer(s, helloWorldGRPCController)
	}
}

// RegisterEventHandlers registers event handlers in a message router.
func RegisterEventHandlers(router *message.Router, subscriber message.Subscriber, logger logur.Logger) error {
	sayHelloHandler := greetingworkerdriver.NewGreeterEventHandler(
		greetingworker.NewGreeterEventLogger(greetingworkeradapter.NewLogger(logger)),
	)

	err := router.AddNoPublisherHandler(
		"log_said_hello",
		"said_hello",
		subscriber,
		sayHelloHandler.SaidHelloTo,
	)
	if err != nil {
		return err
	}

	return nil
}

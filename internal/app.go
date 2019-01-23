package internal

import (
	"net/http"
	"strings"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/goph/emperror"
	"github.com/goph/logur"
	"github.com/gorilla/mux"
	"github.com/sagikazarmark/modern-go-application/.gen/proto/greeting"
	"google.golang.org/grpc"

	"github.com/sagikazarmark/modern-go-application/internal/greeting"
	"github.com/sagikazarmark/modern-go-application/internal/greeting/greetingadapter"
	"github.com/sagikazarmark/modern-go-application/internal/greeting/greetingdriver"
	"github.com/sagikazarmark/modern-go-application/internal/greetingworker"
	"github.com/sagikazarmark/modern-go-application/internal/greetingworker/greetingworkeradapter"
	"github.com/sagikazarmark/modern-go-application/internal/greetingworker/greetingworkerdriver"
	"github.com/sagikazarmark/modern-go-application/internal/httpbin"
)

// NewApp returns a new HTTP application.
func NewApp(logger logur.Logger, publisher message.Publisher, errorHandler emperror.Handler) http.Handler {
	helloService := greeting.NewHelloService(
		greetingadapter.NewSayHelloEvents(publisher),
		greetingadapter.NewLogger(logger),
		errorHandler,
	)
	helloWorldController := greetingdriver.NewHTTPController(helloService, errorHandler)

	router := mux.NewRouter()

	router.Path("/").Methods("GET").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html")
		_, _ = w.Write([]byte(template))
	})

	router.Path("/hello").Methods("POST").HandlerFunc(helloWorldController.SayHello)

	router.PathPrefix("/httpbin").Handler(http.StripPrefix("/httpbin", httpbin.New()))

	helloWorldGRPCController := greetingdriver.NewGRPCController(helloService, errorHandler)

	grpcServer := grpc.NewServer()
	greetingpb.RegisterHelloServiceServer(grpcServer, helloWorldGRPCController)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// This is a partial recreation of gRPC's internal checks:
		// https://github.com/grpc/grpc-go/blob/7346c871b018d255a1d89b3f814a645cc9c5e356/transport/handler_server.go#L61-L75
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			router.ServeHTTP(w, r)
		}
	})
}

// RegisterEventHandlers registers event handlers in a message router.
func RegisterEventHandlers(router *message.Router, subscriber message.Subscriber, logger logur.Logger) error {
	sayHelloHandler := greetingworkerdriver.NewSayHelloEventHandler(
		greetingworker.NewSayHelloEventLogger(greetingworkeradapter.NewLogger(logger)),
	)

	err := router.AddNoPublisherHandler(
		"log_said_hello_to",
		"said_hello_to",
		subscriber,
		sayHelloHandler.SaidHelloTo,
	)
	if err != nil {
		return err
	}

	return nil
}

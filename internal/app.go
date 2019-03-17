package internal

import (
	"net/http"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/goph/emperror"
	"github.com/goph/idgen/ulidgen"
	"github.com/goph/logur"
	"github.com/goph/logur/integrations/watermilllog"
	"github.com/goph/watermillx"
	"github.com/gorilla/mux"
	"github.com/mccutchen/go-httpbin/httpbin"
	"google.golang.org/grpc"

	greetingpb "github.com/sagikazarmark/modern-go-application/.gen/api/proto/greeting"
	"github.com/sagikazarmark/modern-go-application/internal/greeting"
	"github.com/sagikazarmark/modern-go-application/internal/greeting/greetingadapter"
	"github.com/sagikazarmark/modern-go-application/internal/greeting/greetingdriver"
	"github.com/sagikazarmark/modern-go-application/internal/greetingworker"
	"github.com/sagikazarmark/modern-go-application/internal/greetingworker/greetingworkeradapter"
	"github.com/sagikazarmark/modern-go-application/internal/greetingworker/greetingworkerdriver"
	"github.com/sagikazarmark/modern-go-application/internal/landing/landingdriver"
	"github.com/sagikazarmark/modern-go-application/internal/todo"
	"github.com/sagikazarmark/modern-go-application/internal/todo/todoadapter"
	"github.com/sagikazarmark/modern-go-application/internal/todo/tododriver"
)

const todoTopic = "todo"

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

	todoList := todo.NewTodoList(
		ulidgen.NewGenerator(),
		todo.NewInmemoryTodoStore(),
		todoadapter.NewTodoEvents(cqrs.NewEventBus(
			publisher,
			todoTopic,
			watermillx.NewStructNameMarshaler(cqrs.JSONMarshaler{}),
		)),
	)

	router := mux.NewRouter()

	router.Path("/").Methods("GET").Handler(landingdriver.NewHTTPHandler())
	router.PathPrefix("/greeting").Methods("POST").Handler(
		http.StripPrefix("/greeting", greetingdriver.MakeHTTPHandler(greeter, errorHandler)),
	)
	router.PathPrefix("/todos").Handler(
		http.StripPrefix("/todos", tododriver.MakeHTTPHandler(todoList, errorHandler)),
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

	router.AddNoPublisherHandler(
		"log_said_hello",
		"said_hello",
		subscriber,
		sayHelloHandler.SaidHelloTo,
	)

	todoEventProcessor := cqrs.NewEventProcessor(
		[]cqrs.EventHandler{
			tododriver.NewMarkedAsDoneEventHandler(todo.NewLogEventHandler(todoadapter.NewLogger(logger))),
		},
		todoTopic,
		subscriber,
		watermillx.NewStructNameMarshaler(cqrs.JSONMarshaler{}),
		watermilllog.New(logur.WithFields(logger, map[string]interface{}{"component": "watermill"})),
	)

	err := todoEventProcessor.AddHandlersToRouter(router)
	if err != nil {
		return err
	}

	return nil
}

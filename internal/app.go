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
	"github.com/sagikazarmark/ocmux"
	"google.golang.org/grpc"

	todov1beta1 "github.com/sagikazarmark/modern-go-application/.gen/api/proto/todo/v1beta1"
	"github.com/sagikazarmark/modern-go-application/internal/landing/landingdriver"
	"github.com/sagikazarmark/modern-go-application/internal/todo"
	"github.com/sagikazarmark/modern-go-application/internal/todo/todoadapter"
	"github.com/sagikazarmark/modern-go-application/internal/todo/tododriver"
	"github.com/sagikazarmark/modern-go-application/pkg/correlation"
)

const todoTopic = "todo"

// NewApp returns a new HTTP and a new gRPC application.
func NewApp(
	logger logur.Logger,
	publisher message.Publisher,
	errorHandler emperror.Handler,
) (http.Handler, func(*grpc.Server)) {
	var todoList tododriver.TodoList
	{
		todoList = todo.NewList(
			ulidgen.NewGenerator(),
			todo.NewInmemoryStore(),
			todoadapter.NewEventDispatcher(cqrs.NewEventBus(
				publisher,
				todoTopic,
				watermillx.NewStructNameMarshaler(cqrs.JSONMarshaler{}),
			)),
		)
		logger := todoadapter.NewContextAwareLogger(logger, &correlation.ContextExtractor{}).
			WithFields(map[string]interface{}{
				"module": "todo",
			})
		todoList = tododriver.LoggingMiddleware(logger)(todoList)
		todoList = tododriver.InstrumentationMiddleware()(todoList)
	}

	todoListEndpoint := tododriver.MakeEndpoints(todoList)

	router := mux.NewRouter()
	router.Use(ocmux.Middleware())
	router.Use(correlation.HTTPMiddleware(ulidgen.NewGenerator()))

	router.Path("/").Methods("GET").Handler(landingdriver.NewHTTPHandler())
	router.PathPrefix("/todos").Handler(tododriver.MakeHTTPHandler(todoListEndpoint, errorHandler))
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

	return router, func(s *grpc.Server) {
		todov1beta1.RegisterTodoListServer(s, tododriver.MakeGRPCServer(todoListEndpoint, errorHandler))
	}
}

// RegisterEventHandlers registers event handlers in a message router.
func RegisterEventHandlers(router *message.Router, subscriber message.Subscriber, logger logur.Logger) error {
	todoLogger := todoadapter.NewContextAwareLogger(logger, &correlation.ContextExtractor{})
	todoEventProcessor := cqrs.NewEventProcessor(
		[]cqrs.EventHandler{
			tododriver.NewMarkedAsDoneEventHandler(todo.NewLogEventHandler(todoLogger)),
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

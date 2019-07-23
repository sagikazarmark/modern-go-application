package internal

import (
	"net/http"

	"emperror.dev/emperror"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/goph/idgen/ulidgen"
	"github.com/goph/logur"
	"github.com/goph/logur/integrations/watermilllog"
	"github.com/gorilla/mux"
	"github.com/sagikazarmark/ocmux"
	"google.golang.org/grpc"

	todov1beta1 "github.com/sagikazarmark/modern-go-application/.gen/api/proto/todo/v1beta1"
	"github.com/sagikazarmark/modern-go-application/internal/common/commonadapter"
	"github.com/sagikazarmark/modern-go-application/internal/httpbin"
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
	commonLogger := commonadapter.NewContextAwareLogger(logger, &correlation.ContextExtractor{})

	var todoList tododriver.TodoList
	{
		eventBus, _ := cqrs.NewEventBus(
			publisher,
			func(eventName string) string { return todoTopic },
			cqrs.JSONMarshaler{GenerateName: cqrs.StructName},
		)
		todoList = todo.NewList(
			ulidgen.NewGenerator(),
			todo.NewInmemoryStore(),
			todoadapter.NewEventDispatcher(eventBus),
		)
		logger := commonLogger.WithFields(map[string]interface{}{"module": "todo"})
		todoList = tododriver.LoggingMiddleware(logger)(todoList)
		todoList = tododriver.InstrumentationMiddleware()(todoList)
	}

	todoListEndpoint := tododriver.MakeEndpoints(todoList)

	router := mux.NewRouter()
	router.Use(ocmux.Middleware())
	router.Use(correlation.HTTPMiddleware(ulidgen.NewGenerator()))

	router.Path("/").Methods("GET").Handler(landingdriver.NewHTTPHandler())
	router.PathPrefix("/todos").Handler(tododriver.MakeHTTPHandler(todoListEndpoint, errorHandler))
	router.PathPrefix("/graphql").Handler(tododriver.MakeGraphQLHandler(todoListEndpoint, errorHandler))
	router.PathPrefix("/httpbin").Handler(http.StripPrefix(
		"/httpbin",
		httpbin.MakeHTTPHandler(commonLogger.WithFields(map[string]interface{}{"module": "httpbin"})),
	))

	return router, func(s *grpc.Server) {
		todov1beta1.RegisterTodoListServer(s, tododriver.MakeGRPCServer(todoListEndpoint, errorHandler))
	}
}

// RegisterEventHandlers registers event handlers in a message router.
func RegisterEventHandlers(router *message.Router, subscriber message.Subscriber, logger logur.Logger) error {
	commonLogger := commonadapter.NewContextAwareLogger(logger, &correlation.ContextExtractor{})
	todoEventProcessor, _ := cqrs.NewEventProcessor(
		[]cqrs.EventHandler{
			tododriver.NewMarkedAsDoneEventHandler(todo.NewLogEventHandler(commonLogger)),
		},
		func(eventName string) string { return todoTopic },
		func(handlerName string) (message.Subscriber, error) { return subscriber, nil },
		cqrs.JSONMarshaler{GenerateName: cqrs.StructName},
		watermilllog.New(logur.WithFields(logger, map[string]interface{}{"component": "watermill"})),
	)

	err := todoEventProcessor.AddHandlersToRouter(router)
	if err != nil {
		return err
	}

	return nil
}

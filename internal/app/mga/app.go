package mga

import (
	"net/http"

	"emperror.dev/emperror"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/go-kit/kit/endpoint"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/goph/idgen/ulidgen"
	"github.com/gorilla/mux"
	"github.com/sagikazarmark/kitx/correlation"
	kitxendpoint "github.com/sagikazarmark/kitx/endpoint"
	kitxgrpc "github.com/sagikazarmark/kitx/transport/grpc"
	kitxhttp "github.com/sagikazarmark/kitx/transport/http"
	"google.golang.org/grpc"
	watermilllog "logur.dev/integration/watermill"
	"logur.dev/logur"

	todov1beta1 "github.com/sagikazarmark/modern-go-application/.gen/api/proto/todo/v1beta1"
	"github.com/sagikazarmark/modern-go-application/internal/app/mga/httpbin"
	"github.com/sagikazarmark/modern-go-application/internal/app/mga/landing/landingdriver"
	"github.com/sagikazarmark/modern-go-application/internal/app/mga/todo"
	"github.com/sagikazarmark/modern-go-application/internal/app/mga/todo/tododriver"
	"github.com/sagikazarmark/modern-go-application/internal/app/mga/todo/todogen"
	"github.com/sagikazarmark/modern-go-application/internal/common/commonadapter"
	"github.com/sagikazarmark/modern-go-application/internal/platform/appkit"
)

const todoTopic = "todo"

// InitializeApp initializes a new HTTP and a new gRPC application.
func InitializeApp(
	httpRouter *mux.Router,
	grpcServer *grpc.Server,
	publisher message.Publisher,
	logger logur.Logger,
	errorHandler emperror.Handler,
) {
	commonLogger := commonadapter.NewContextAwareLogger(logger, appkit.ContextExtractor{})

	endpointMiddleware := []endpoint.Middleware{
		correlation.Middleware(),
	}

	httpServerOptions := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(emperror.MakeContextAware(errorHandler)),
		kithttp.ServerErrorEncoder(kitxhttp.ProblemErrorEncoder),
		kithttp.ServerBefore(correlation.HTTPToContext()),
	}

	grpcServerOptions := []kitgrpc.ServerOption{
		kitgrpc.ServerErrorHandler(emperror.MakeContextAware(errorHandler)),
		kitgrpc.ServerBefore(correlation.GRPCToContext()),
	}

	{
		logger := commonLogger.WithFields(map[string]interface{}{"module": "todo"})
		errorHandler := emperror.MakeContextAware(emperror.WithDetails(errorHandler, "module", "todo"))

		eventBus, _ := cqrs.NewEventBus(
			publisher,
			func(eventName string) string { return todoTopic },
			cqrs.JSONMarshaler{GenerateName: cqrs.StructName},
		)

		service := todo.NewService(
			ulidgen.NewGenerator(),
			todo.NewInMemoryStore(),
			todogen.NewEventDispatcher(eventBus),
		)
		service = tododriver.LoggingMiddleware(logger)(service)
		service = tododriver.InstrumentationMiddleware()(service)

		endpoints := tododriver.TraceEndpoints(tododriver.MakeEndpoints(
			service,
			kitxendpoint.Chain(endpointMiddleware...),
			appkit.EndpointLogger(logger),
		))

		tododriver.RegisterHTTPHandlers(
			endpoints,
			httpRouter.PathPrefix("/todos").Subrouter(),
			kitxhttp.ServerOptions(httpServerOptions),
			kithttp.ServerErrorHandler(errorHandler),
		)

		todov1beta1.RegisterTodoListServer(
			grpcServer,
			tododriver.MakeGRPCServer(
				endpoints,
				kitxgrpc.ServerOptions(grpcServerOptions),
				kitgrpc.ServerErrorHandler(errorHandler),
			),
		)

		httpRouter.PathPrefix("/graphql").Handler(tododriver.MakeGraphQLHandler(endpoints, errorHandler))
	}

	landingdriver.RegisterHTTPHandlers(httpRouter)
	httpRouter.PathPrefix("/httpbin").Handler(http.StripPrefix(
		"/httpbin",
		httpbin.MakeHTTPHandler(commonLogger.WithFields(map[string]interface{}{"module": "httpbin"})),
	))
}

// RegisterEventHandlers registers event handlers in a message router.
func RegisterEventHandlers(router *message.Router, subscriber message.Subscriber, logger logur.Logger) error {
	commonLogger := commonadapter.NewContextAwareLogger(logger, appkit.ContextExtractor{})
	todoEventProcessor, _ := cqrs.NewEventProcessor(
		[]cqrs.EventHandler{
			todogen.NewMarkedAsDoneEventHandler(todo.NewLogEventHandler(commonLogger), "marked_as_done"),
		},
		func(eventName string) string { return todoTopic },
		func(handlerName string) (message.Subscriber, error) { return subscriber, nil },
		cqrs.JSONMarshaler{GenerateName: cqrs.StructName},
		watermilllog.New(logur.WithField(logger, "component", "watermill")),
	)

	err := todoEventProcessor.AddHandlersToRouter(router)
	if err != nil {
		return err
	}

	return nil
}

package mga

import (
	"context"
	"net/http"

	"emperror.dev/emperror"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/go-kit/kit/endpoint"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/goph/idgen/ulidgen"
	"github.com/gorilla/mux"
	appkitendpoint "github.com/sagikazarmark/appkit/endpoint"
	appkiterrors "github.com/sagikazarmark/appkit/errors"
	appkithttp "github.com/sagikazarmark/appkit/transport/http"
	"github.com/sagikazarmark/kitx/correlation"
	kitxendpoint "github.com/sagikazarmark/kitx/endpoint"
	kitxtransport "github.com/sagikazarmark/kitx/transport"
	kitxgrpc "github.com/sagikazarmark/kitx/transport/grpc"
	kitxhttp "github.com/sagikazarmark/kitx/transport/http"
	"go.opencensus.io/trace"
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
)

const todoTopic = "todo"

// InitializeApp initializes a new HTTP and a new gRPC application.
func InitializeApp(
	httpRouter *mux.Router,
	grpcServer *grpc.Server,
	publisher message.Publisher,
	logger logur.LoggerFacade,
	errorHandler emperror.ErrorHandlerFacade,
) {
	logger = logur.WithContextExtractor(logger, contextExtractor)
	errorHandler = emperror.WithContextExtractor(errorHandler, contextExtractor)

	commonLogger := commonadapter.NewContextAwareLogger(logger, contextExtractor)

	endpointMiddleware := []endpoint.Middleware{
		correlation.Middleware(),
		appkitendpoint.LoggingMiddleware(logger),
		appkitendpoint.ClientErrorMiddleware,
	}

	transportErrorHandler := kitxtransport.NewErrorHandler(emperror.WithFilter(
		errorHandler,
		appkiterrors.IsClientError, // filter out client errors
	))

	httpServerOptions := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transportErrorHandler),
		kithttp.ServerErrorEncoder(kitxhttp.NewJSONProblemErrorEncoder(appkithttp.NewDefaultProblemConverter())),
		kithttp.ServerBefore(correlation.HTTPToContext()),
	}

	grpcServerOptions := []kitgrpc.ServerOption{
		kitgrpc.ServerErrorHandler(transportErrorHandler),
		kitgrpc.ServerBefore(correlation.GRPCToContext()),
	}

	{
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
		service = tododriver.LoggingMiddleware(commonLogger)(service)
		service = tododriver.InstrumentationMiddleware()(service)

		endpoints := tododriver.TraceEndpoints(tododriver.MakeEndpoints(
			service,
			kitxendpoint.Combine(endpointMiddleware...),
		))

		tododriver.RegisterHTTPHandlers(
			endpoints,
			httpRouter.PathPrefix("/todos").Subrouter(),
			kitxhttp.ServerOptions(httpServerOptions),
		)

		todov1beta1.RegisterTodoListServer(
			grpcServer,
			tododriver.MakeGRPCServer(
				endpoints,
				kitxgrpc.ServerOptions(grpcServerOptions),
				kitgrpc.ServerErrorHandler(kitxtransport.NewErrorHandler(errorHandler)),
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
func RegisterEventHandlers(router *message.Router, subscriber message.Subscriber, logger logur.LoggerFacade) error {
	logger = logur.WithContextExtractor(logger, contextExtractor)
	commonLogger := commonadapter.NewContextAwareLogger(logger, contextExtractor)

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

func contextExtractor(ctx context.Context) map[string]interface{} {
	fields := make(map[string]interface{})

	if correlationID, ok := correlation.FromContext(ctx); ok {
		fields["correlation_id"] = correlationID
	}

	if operationName, ok := kitxendpoint.OperationName(ctx); ok {
		fields["operation_name"] = operationName
	}

	if span := trace.FromContext(ctx); span != nil {
		spanCtx := span.SpanContext()

		fields["trace_id"] = spanCtx.TraceID.String()
		fields["span_id"] = spanCtx.SpanID.String()
	}

	return fields
}

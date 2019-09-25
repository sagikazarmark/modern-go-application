package mga

import (
	"context"
	"net/http"

	"emperror.dev/emperror"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/go-kit/kit/endpoint"
	kitoc "github.com/go-kit/kit/tracing/opencensus"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/goph/idgen/ulidgen"
	"github.com/gorilla/mux"
	"github.com/sagikazarmark/kitx/correlation"
	kitxendpoint "github.com/sagikazarmark/kitx/endpoint"
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
	"github.com/sagikazarmark/modern-go-application/internal/common"
	"github.com/sagikazarmark/modern-go-application/internal/common/commonadapter"
	"github.com/sagikazarmark/modern-go-application/internal/platform/appkit"
)

const todoTopic = "todo"

// ContextExtractor extracts values from a context.
type ContextExtractor struct{}

// Extract extracts values from a context.
func (ContextExtractor) Extract(ctx context.Context) map[string]interface{} {
	fields := make(map[string]interface{})

	if correlationID, ok := correlation.FromContext(ctx); ok {
		fields["correlation_id"] = correlationID
	}

	if span := trace.FromContext(ctx); span != nil {
		spanCtx := span.SpanContext()
		fields["trace_id"] = spanCtx.TraceID.String()
		fields["span_id"] = spanCtx.SpanID.String()
	}

	return fields
}

// InitializeApp initializes a new HTTP and a new gRPC application.
func InitializeApp(
	httpRouter *mux.Router,
	grpcServer *grpc.Server,
	publisher message.Publisher,
	logger logur.Logger,
	errorHandler emperror.Handler,
) {
	commonLogger := commonadapter.NewContextAwareLogger(logger, ContextExtractor{})

	endpointFactory := func(logger common.Logger) kitxendpoint.Factory {
		return kitxendpoint.NewFactory(
			kitxendpoint.Middleware(correlation.Middleware()),
			func(name string) endpoint.Middleware { return kitoc.TraceEndpoint(name) },
			appkit.EndpointLoggerFactory(logger),
		)
	}

	httpServerFactory := func(errorHandler common.ErrorHandler) kitxhttp.ServerFactory {
		return kitxhttp.NewServerFactory(
			kithttp.ServerErrorHandler(errorHandler),
			kithttp.ServerBefore(correlation.HTTPToContext()),
		)
	}

	grpcServerFactory := func(errorHandler common.ErrorHandler) kitxgrpc.ServerFactory {
		return kitxgrpc.NewServerFactory(
			kitgrpc.ServerErrorHandler(errorHandler),
			kitgrpc.ServerBefore(correlation.GRPCToContext()),
		)
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

		endpoints := tododriver.MakeEndpoints(service, endpointFactory(logger))

		tododriver.RegisterHTTPHandlers(
			endpoints,
			httpServerFactory(errorHandler),
			httpRouter.PathPrefix("/todos").Subrouter(),
		)

		todov1beta1.RegisterTodoListServer(
			grpcServer,
			tododriver.MakeGRPCServer(endpoints, grpcServerFactory(errorHandler)),
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
	commonLogger := commonadapter.NewContextAwareLogger(logger, ContextExtractor{})
	todoEventProcessor, _ := cqrs.NewEventProcessor(
		[]cqrs.EventHandler{
			todogen.NewMarkedAsDoneEventHandler(todo.NewLogEventHandler(commonLogger), "marked_as_done"),
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

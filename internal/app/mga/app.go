package mga

import (
	"context"
	"net/http"
	"time"

	"emperror.dev/emperror"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/go-kit/kit/endpoint"
	kitoc "github.com/go-kit/kit/tracing/opencensus"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/goph/idgen/ulidgen"
	"github.com/gorilla/mux"
	"github.com/sagikazarmark/ocmux"
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
	"github.com/sagikazarmark/modern-go-application/pkg/kitx/correlation"
	kitxendpoint "github.com/sagikazarmark/modern-go-application/pkg/kitx/endpoint"
	kitxgrpc "github.com/sagikazarmark/modern-go-application/pkg/kitx/transport/grpc"
	kitxhttp "github.com/sagikazarmark/modern-go-application/pkg/kitx/transport/http"
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

// NewApp returns a new HTTP and a new gRPC application.
func NewApp(
	logger logur.Logger,
	publisher message.Publisher,
	errorHandler emperror.Handler,
) (http.Handler, func(*grpc.Server)) {
	commonLogger := commonadapter.NewContextAwareLogger(logger, ContextExtractor{})

	var todoList todo.Service
	{
		eventBus, _ := cqrs.NewEventBus(
			publisher,
			func(eventName string) string { return todoTopic },
			cqrs.JSONMarshaler{GenerateName: cqrs.StructName},
		)
		todoList = todo.NewService(
			ulidgen.NewGenerator(),
			todo.NewInmemoryStore(),
			todogen.NewEventDispatcher(eventBus),
		)
		logger := commonLogger.WithFields(map[string]interface{}{"module": "todo"})
		todoList = tododriver.LoggingMiddleware(logger)(todoList)
		todoList = tododriver.InstrumentationMiddleware()(todoList)
	}

	endpointFactory := kitxendpoint.NewFactory(
		kitxendpoint.Middleware(correlation.Middleware()),
		func(name string) endpoint.Middleware { return kitoc.TraceEndpoint(name) },
		func(name string) endpoint.Middleware {
			return func(e endpoint.Endpoint) endpoint.Endpoint {
				return func(ctx context.Context, request interface{}) (interface{}, error) {
					logger := commonLogger.WithContext(ctx).WithFields(map[string]interface{}{"module": "todo"})

					logger.Trace("processing request", map[string]interface{}{
						"operation": name,
					})

					defer func(begin time.Time) {
						logger.Trace("processing request finished", map[string]interface{}{
							"operation": name,
							"took":      time.Since(begin),
						})
					}(time.Now())

					return e(ctx, request)
				}
			}
		},
	)

	todoListEndpoint := tododriver.MakeEndpoints(todoList, endpointFactory)

	ctxErrorHandler := emperror.MakeContextAware(errorHandler)

	router := mux.NewRouter()
	router.Use(ocmux.Middleware())

	httpServerFactory := kitxhttp.NewServerFactory(
		kithttp.ServerErrorHandler(ctxErrorHandler),
		kithttp.ServerBefore(correlation.HTTPToContext()),
	)

	landingdriver.RegisterHTTPHandlers(router)
	tododriver.RegisterHTTPHandlers(todoListEndpoint, httpServerFactory, router.PathPrefix("/todos").Subrouter())
	router.PathPrefix("/graphql").Handler(tododriver.MakeGraphQLHandler(todoListEndpoint, ctxErrorHandler))
	router.PathPrefix("/httpbin").Handler(http.StripPrefix(
		"/httpbin",
		httpbin.MakeHTTPHandler(commonLogger.WithFields(map[string]interface{}{"module": "httpbin"})),
	))

	grpcServerFactory := kitxgrpc.NewServerFactory(
		kitgrpc.ServerErrorHandler(ctxErrorHandler),
		kitgrpc.ServerBefore(correlation.GRPCToContext()),
	)

	return router, func(s *grpc.Server) {
		todov1beta1.RegisterTodoListServer(s, tododriver.MakeGRPCServer(todoListEndpoint, grpcServerFactory))
	}
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

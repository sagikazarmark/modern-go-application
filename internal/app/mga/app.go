package mga

import (
	"context"
	"database/sql"
	"net/http"

	entsql "entgo.io/ent/dialect/sql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/tracing/opencensus"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/goph/idgen/ulidgen"
	"github.com/gorilla/mux"
	appkitendpoint "github.com/sagikazarmark/appkit/endpoint"
	appkithttp "github.com/sagikazarmark/appkit/transport/http"
	"github.com/sagikazarmark/kitx/correlation"
	kitxendpoint "github.com/sagikazarmark/kitx/endpoint"
	kitxtransport "github.com/sagikazarmark/kitx/transport"
	kitxgrpc "github.com/sagikazarmark/kitx/transport/grpc"
	kitxhttp "github.com/sagikazarmark/kitx/transport/http"
	todov1 "github.com/sagikazarmark/todobackend-go-kit/api/todo/v1"
	"github.com/sagikazarmark/todobackend-go-kit/todo"
	"github.com/sagikazarmark/todobackend-go-kit/todo/tododriver"
	"google.golang.org/grpc"
	watermilllog "logur.dev/integration/watermill"

	"github.com/sagikazarmark/modern-go-application/internal/app/mga/httpbin"
	"github.com/sagikazarmark/modern-go-application/internal/app/mga/landing/landingdriver"
	todo2 "github.com/sagikazarmark/modern-go-application/internal/app/mga/todo"
	"github.com/sagikazarmark/modern-go-application/internal/app/mga/todo/todoadapter"
	"github.com/sagikazarmark/modern-go-application/internal/app/mga/todo/todoadapter/ent"
	"github.com/sagikazarmark/modern-go-application/internal/app/mga/todo/todoadapter/ent/migrate"
	tododriver2 "github.com/sagikazarmark/modern-go-application/internal/app/mga/todo/tododriver"
	"github.com/sagikazarmark/modern-go-application/internal/app/mga/todo/todogen"
	"github.com/sagikazarmark/modern-go-application/static/templates"
)

const todoTopic = "todo"

// InitializeApp initializes a new HTTP and a new gRPC application.
func InitializeApp(
	httpRouter *mux.Router,
	grpcServer *grpc.Server,
	publisher message.Publisher,
	storage string,
	db *sql.DB,
	logger Logger,
	errorHandler ErrorHandler, // nolint: interfacer
) {
	endpointMiddleware := []endpoint.Middleware{
		correlation.Middleware(),
		opencensus.TraceEndpoint("", opencensus.WithSpanName(func(ctx context.Context, _ string) string {
			name, _ := kitxendpoint.OperationName(ctx)

			return name
		})),
		appkitendpoint.LoggingMiddleware(logger),
	}

	transportErrorHandler := kitxtransport.NewErrorHandler(errorHandler)

	httpServerOptions := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transportErrorHandler),
		kithttp.ServerErrorEncoder(kitxhttp.NewJSONProblemErrorEncoder(appkithttp.NewDefaultProblemConverter())),
		kithttp.ServerBefore(correlation.HTTPToContext(), kithttp.PopulateRequestContext),
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

		var store todo.Store = todo.NewInMemoryStore()
		if storage == "database" {
			client := ent.NewClient(ent.Driver(entsql.OpenDB("mysql", db)))
			err := client.Schema.Create(
				context.Background(),
				migrate.WithDropIndex(true),
				migrate.WithDropColumn(true),
			)
			if err != nil {
				panic(err)
			}

			store = todoadapter.NewEntStore(client)
		}

		service := todo.NewService(ulidgen.NewGenerator(), store)
		service = todo2.EventMiddleware(todogen.NewEventDispatcher(eventBus))(service)
		service = tododriver2.LoggingMiddleware(logger)(service)
		service = tododriver2.InstrumentationMiddleware()(service)

		endpoints := tododriver.MakeEndpoints(
			service,
			kitxendpoint.Combine(endpointMiddleware...),
		)

		tododriver.RegisterHTTPHandlers(
			endpoints,
			httpRouter.PathPrefix("/todos").Subrouter(),
			kitxhttp.ServerOptions(httpServerOptions),
		)
		todov1.RegisterTodoListServiceServer(
			grpcServer,
			tododriver.MakeGRPCServer(endpoints, kitxgrpc.ServerOptions(grpcServerOptions)),
		)
		httpRouter.PathPrefix("/graphql").Handler(handler.NewDefaultServer(
			tododriver.MakeGraphQLSchema(endpoints),
		))
	}

	landingdriver.RegisterHTTPHandlers(httpRouter, templates.Files())
	httpRouter.PathPrefix("/httpbin").Handler(http.StripPrefix(
		"/httpbin",
		httpbin.MakeHTTPHandler(logger.WithFields(map[string]interface{}{"module": "httpbin"})),
	))
}

// RegisterEventHandlers registers event handlers in a message router.
func RegisterEventHandlers(router *message.Router, subscriber message.Subscriber, logger Logger) error {
	todoEventProcessor, _ := cqrs.NewEventProcessor(
		[]cqrs.EventHandler{
			todogen.NewMarkedAsCompleteEventHandler(todo2.NewLogEventHandler(logger), "marked_as_complete"),
		},
		func(eventName string) string { return todoTopic },
		func(handlerName string) (message.Subscriber, error) { return subscriber, nil },
		cqrs.JSONMarshaler{GenerateName: cqrs.StructName},
		watermilllog.New(logger.WithFields(map[string]interface{}{"component": "watermill"})),
	)

	err := todoEventProcessor.AddHandlersToRouter(router)
	if err != nil {
		return err
	}

	return nil
}

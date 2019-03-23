package todocli

import (
	"github.com/goph/emperror"
	"github.com/spf13/cobra"
	"go.opencensus.io/exporter/jaeger"
	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/trace"
	"google.golang.org/grpc"

	"github.com/sagikazarmark/modern-go-application/internal/todocli/command"

	todov1beta1 "github.com/sagikazarmark/modern-go-application/.gen/api/proto/todo/v1alpha1"
)

// Configure configures a root command.
func Configure(rootCmd *cobra.Command) {
	var address string

	flags := rootCmd.PersistentFlags()

	flags.StringVar(&address, "address", "127.0.0.1:8001", "Todo service address")

	c := &context{}

	var grpcConn *grpc.ClientConn

	rootCmd.PersistentPreRunE = func(_ *cobra.Command, _ []string) error {
		conn, err := grpc.Dial(
			address,
			grpc.WithInsecure(),
			grpc.WithStatsHandler(&ocgrpc.ClientHandler{
				StartOptions: trace.StartOptions{
					Sampler:  trace.AlwaysSample(),
					SpanKind: trace.SpanKindClient,
				},
			}),
		)
		if err != nil {
			return emperror.Wrap(err, "failed to dial service")
		}

		exporter, err := jaeger.NewExporter(jaeger.Options{
			CollectorEndpoint: "http://localhost:14268/api/traces?format=jaeger.thrift",
			Process: jaeger.Process{
				ServiceName: "todocli",
			},
		})
		if err != nil {
			return emperror.Wrap(err, "failed to create exporter")
		}

		trace.RegisterExporter(exporter)

		grpcConn = conn

		c.client = todov1beta1.NewTodoListClient(conn)

		return nil
	}

	rootCmd.PersistentPostRunE = func(_ *cobra.Command, _ []string) error {
		return grpcConn.Close()
	}

	command.AddCommands(rootCmd, c)
}

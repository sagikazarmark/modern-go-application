package todocli

import (
	"contrib.go.opencensus.io/exporter/ocagent"
	"emperror.dev/errors"
	todov1 "github.com/sagikazarmark/todobackend-go-kit/api/todo/v1"
	"github.com/spf13/cobra"
	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/trace"
	"google.golang.org/grpc"

	"github.com/sagikazarmark/modern-go-application/internal/app/todocli/command"
)

// Configure configures a root command.
func Configure(rootCmd *cobra.Command) {
	var address string

	flags := rootCmd.PersistentFlags()

	flags.StringVar(&address, "address", "127.0.0.1:8001", "Todo service address")

	c := &context{}

	var grpcConn *grpc.ClientConn
	var ocagentExporter *ocagent.Exporter

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
			return errors.WrapIf(err, "failed to dial service")
		}

		// Configure OpenCensus exporter
		exporter, err := ocagent.NewExporter(ocagent.WithServiceName("todocli"), ocagent.WithInsecure())
		if err != nil {
			return errors.WrapIf(err, "failed to create exporter")
		}

		ocagentExporter = exporter

		trace.RegisterExporter(exporter)

		grpcConn = conn

		c.client = todov1.NewTodoListServiceClient(conn)

		return nil
	}

	rootCmd.PersistentPostRunE = func(_ *cobra.Command, _ []string) error {
		ocagentExporter.Flush()

		return grpcConn.Close()
	}

	command.AddCommands(rootCmd, c)
}

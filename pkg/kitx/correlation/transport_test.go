package correlation

import (
	"context"
	"net/http"
	"testing"

	"google.golang.org/grpc/metadata"
)

func TestHTTPToContext(t *testing.T) {
	reqFunc := HTTPToContext()

	t.Run("no_header", func(t *testing.T) {
		ctx := reqFunc(context.Background(), &http.Request{})

		if ctx.Value(correlationIDContextKey) != nil {
			t.Error("context should not contain the encoded correlation ID")
		}
	})

	t.Run("default_header", func(t *testing.T) {
		headerVal := "2314"

		header := http.Header{}
		header.Set("Correlation-ID", headerVal)

		ctx := reqFunc(context.Background(), &http.Request{Header: header})

		cid := ctx.Value(correlationIDContextKey).(string)
		if cid != headerVal {
			t.Errorf(
				"context does not contain the expected encoded correlation ID\nexpected: %s\nactual:   %s",
				headerVal,
				cid,
			)
		}
	})

	t.Run("custom_header", func(t *testing.T) {
		reqFunc := HTTPToContext("Correlation-ID", "X-Correlation-ID")

		headerVal := "3412"

		header := http.Header{}
		header.Set("X-Correlation-ID", headerVal)

		ctx := reqFunc(context.Background(), &http.Request{Header: header})

		cid := ctx.Value(correlationIDContextKey).(string)
		if cid != headerVal {
			t.Errorf(
				"context does not contain the expected encoded correlation ID\nexpected: %s\nactual:   %s",
				headerVal,
				cid,
			)
		}
	})
}

func TestGRPCToContext(t *testing.T) {
	reqFunc := GRPCToContext()

	t.Run("no_header", func(t *testing.T) {
		ctx := reqFunc(context.Background(), metadata.MD{})

		if ctx.Value(correlationIDContextKey) != nil {
			t.Error("context should not contain the encoded correlation ID")
		}
	})

	t.Run("default_header", func(t *testing.T) {
		headerVal := "2431"

		md := metadata.MD{}
		md.Set("correlation-id", headerVal)

		ctx := reqFunc(context.Background(), md)

		cid := ctx.Value(correlationIDContextKey).(string)
		if cid != headerVal {
			t.Errorf(
				"context does not contain the expected encoded correlation ID\nexpected: %s\nactual:   %s",
				headerVal,
				cid,
			)
		}
	})

	t.Run("custom_header", func(t *testing.T) {
		reqFunc := GRPCToContext("correlation-id", "x-correlation-id")

		headerVal := "1234"

		md := metadata.MD{}
		md.Set("x-correlation-id", headerVal)

		ctx := reqFunc(context.Background(), md)

		cid := ctx.Value(correlationIDContextKey).(string)
		if cid != headerVal {
			t.Errorf(
				"context does not contain the expected encoded correlation ID\nexpected: %s\nactual:   %s",
				headerVal,
				cid,
			)
		}
	})
}

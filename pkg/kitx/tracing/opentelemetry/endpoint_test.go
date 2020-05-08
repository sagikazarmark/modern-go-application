package opentelemetry_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/key"
	"go.opentelemetry.io/otel/api/trace/testtrace"
	"google.golang.org/grpc/codes"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/lb"

	"github.com/sagikazarmark/modern-go-application/pkg/kitx/tracing/opentelemetry"
)

const (
	span1 = ""
	span2 = "SPAN-2"
	span3 = "SPAN-3"
	span4 = "SPAN-4"
	span5 = "SPAN-5"
	span6 = "SPAN-6"
)

var (
	err1 = errors.New("some error")
	err2 = errors.New("other error")
	err3 = errors.New("some business error")
	err4 = errors.New("other business error")
)

// compile time assertion
var _ endpoint.Failer = failedResponse{}

type failedResponse struct {
	err error
}

func (r failedResponse) Failed() error { return r.err }

func passEndpoint(_ context.Context, req interface{}) (interface{}, error) {
	if err, _ := req.(error); err != nil {
		return nil, err
	}
	return req, nil
}

func withName(name string) opentelemetry.EndpointOption {
	return opentelemetry.WithSpanName(func(ctx context.Context, _ string) string {
		return name
	})
}

func TestTraceEndpoint(t *testing.T) {
	ctx := context.Background()

	tracer := testtrace.NewTracer()

	// span 1
	span1Attrs := []core.KeyValue{
		key.String("string", "value"),
		key.Int64("int64", 42),
	}
	mw := opentelemetry.TraceEndpoint(
		withName(span1),
		opentelemetry.WithTracer(tracer),
		opentelemetry.WithEndpointAttributes(span1Attrs...),
	)
	mw(endpoint.Nop)(ctx, nil)

	// span 2
	opts := opentelemetry.EndpointOptions{}
	mw = opentelemetry.TraceEndpoint(
		opentelemetry.WithEndpointConfig(opts),
		withName(span2),
		opentelemetry.WithTracer(tracer),
	)
	mw(passEndpoint)(ctx, err1)

	// span3
	mw = opentelemetry.TraceEndpoint(withName(span3), opentelemetry.WithTracer(tracer))
	ep := lb.Retry(5, 1*time.Second, lb.NewRoundRobin(sd.FixedEndpointer{passEndpoint}))
	mw(ep)(ctx, err2)

	// span4
	mw = opentelemetry.TraceEndpoint(withName(span4), opentelemetry.WithTracer(tracer))
	mw(passEndpoint)(ctx, failedResponse{err: err3})

	// span5
	mw = opentelemetry.TraceEndpoint(
		withName(span5),
		opentelemetry.WithTracer(tracer),
		opentelemetry.WithIgnoreBusinessError(true),
	)
	mw(passEndpoint)(ctx, failedResponse{err: err4})

	// span6
	span6Attrs := []core.KeyValue{
		key.String("string", "value"),
		key.Int64("int64", 42),
	}
	mw = opentelemetry.TraceEndpoint(
		opentelemetry.WithDefaultName(span6),
		opentelemetry.WithTracer(tracer),
		opentelemetry.WithSpanAttributes(func(ctx context.Context) []core.KeyValue {
			return span6Attrs
		}),
	)
	mw(endpoint.Nop)(ctx, nil)

	// TODO: add a test case with a global trace provider

	// check span count
	spans := tracer.Spans()
	if want, have := 6, len(spans); want != have {
		t.Fatalf("incorrected number of spans, wanted %d, got %d", want, have)
	}

	// test span 1
	span := spans[0]
	if want, have := codes.OK, span.StatusCode(); want != have {
		t.Errorf("incorrect status code, wanted %d, got %d", want, have)
	}

	if want, have := opentelemetry.TraceEndpointDefaultName, span.Name(); want != have {
		t.Errorf("incorrect span name, wanted %q, got %q", want, have)
	}

	if want, have := 2, len(span.Attributes()); want != have {
		t.Fatalf("incorrect attribute count, wanted %d, got %d", want, have)
	}

	// test span 2
	span = spans[1]
	if want, have := codes.Unknown, span.StatusCode(); want != have {
		t.Errorf("incorrect status code, wanted %d, got %d", want, have)
	}

	if want, have := span2, span.Name(); want != have {
		t.Errorf("incorrect span name, wanted %q, got %q", want, have)
	}

	if want, have := 0, len(span.Attributes()); want != have {
		t.Fatalf("incorrect attribute count, wanted %d, got %d", want, have)
	}

	// test span 3
	span = spans[2]
	if want, have := codes.Unknown, span.StatusCode(); want != have {
		t.Errorf("incorrect status code, wanted %d, got %d", want, have)
	}

	if want, have := span3, span.Name(); want != have {
		t.Errorf("incorrect span name, wanted %q, got %q", want, have)
	}

	if want, have := 5, len(span.Attributes()); want != have {
		t.Fatalf("incorrect attribute count, wanted %d, got %d", want, have)
	}

	// test span 4
	span = spans[3]
	if want, have := codes.Unknown, span.StatusCode(); want != have {
		t.Errorf("incorrect status code, wanted %d, got %d", want, have)
	}

	if want, have := span4, span.Name(); want != have {
		t.Errorf("incorrect span name, wanted %q, got %q", want, have)
	}

	if want, have := 1, len(span.Attributes()); want != have {
		t.Fatalf("incorrect attribute count, wanted %d, got %d", want, have)
	}

	// test span 5
	span = spans[4]
	if want, have := codes.OK, span.StatusCode(); want != have {
		t.Errorf("incorrect status code, wanted %d, got %d", want, have)
	}

	if want, have := span5, span.Name(); want != have {
		t.Errorf("incorrect span name, wanted %q, got %q", want, have)
	}

	if want, have := 1, len(span.Attributes()); want != have {
		t.Fatalf("incorrect attribute count, wanted %d, got %d", want, have)
	}

	// test span 6
	span = spans[5]
	if want, have := span6, span.Name(); want != have {
		t.Errorf("incorrect span name, wanted %q, got %q", want, have)
	}

	if want, have := 2, len(span.Attributes()); want != have {
		t.Fatalf("incorrect attribute count, wanted %d, got %d", want, have)
	}
}

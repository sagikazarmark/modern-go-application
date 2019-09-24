package correlation

import (
	"context"
	"testing"
)

func TestFromContext(t *testing.T) {
	ctx := context.WithValue(context.Background(), correlationIDContextKey, "id")

	id, ok := FromContext(ctx)
	if !ok {
		t.Fatal("correlation ID not found in the context")
	}

	if want, have := "id", id; want != have {
		t.Errorf("unexpected correlation ID\nexpected: %s\nactual:   %s", want, have)
	}
}

func TestToContext(t *testing.T) {
	ctx := ToContext(context.Background(), "id")

	id, ok := ctx.Value(correlationIDContextKey).(string)
	if !ok {
		t.Fatal("correlation ID not found in the context")
	}

	if want, have := "id", id; want != have {
		t.Errorf("unexpected correlation ID\nexpected: %s\nactual:   %s", want, have)
	}
}

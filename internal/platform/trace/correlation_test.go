package trace

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithCorrelationID(t *testing.T) {
	ctx := WithCorrelationID(context.Background(), "id")

	assert.Equal(t, "id", ctx.Value(correlationID))
}

func TestCorrelationID(t *testing.T) {
	ctx := context.WithValue(context.Background(), correlationID, "id")

	id, ok := CorrelationID(ctx)
	require.True(t, ok)
	assert.Equal(t, "id", id)
}

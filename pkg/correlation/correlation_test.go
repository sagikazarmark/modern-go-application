package correlation

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithID(t *testing.T) {
	ctx := WithID(context.Background(), "id")

	assert.Equal(t, "id", ctx.Value(correlationID))
}

func TestID(t *testing.T) {
	ctx := context.WithValue(context.Background(), correlationID, "id")

	id, ok := ID(ctx)
	require.True(t, ok)
	assert.Equal(t, "id", id)
}

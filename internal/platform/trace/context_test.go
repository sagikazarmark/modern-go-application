package trace

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContextExtractor_Extract(t *testing.T) {
	ctx := WithCorrelationID(context.Background(), "id")

	extractor := &ContextExtractor{}

	fields := extractor.Extract(ctx)

	assert.Equal(t, fields, map[string]interface{}{"correlation_id": "id"})
}

package greetingadapter

import (
	"context"
	"testing"

	"github.com/sagikazarmark/modern-go-application/internal/greeting"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSayHelloEvents_SaidHelloTo(t *testing.T) {
	publisher := &publisherStub{}

	events := NewSayHelloEvents(publisher)

	event := greeting.SaidHelloTo{
		Message: "Hello, John!",
		Who:     "John",
	}

	err := events.SaidHelloTo(context.Background(), event)
	require.NoError(t, err)

	assert.Equal(t, saidHelloToTopic, "said_hello_to")
	assert.Equal(t, saidHelloToTopic, publisher.topic)
	assert.Equal(t, string(publisher.messages[0].Payload), "{\"Message\":\"Hello, John!\",\"Who\":\"John\"}")
}

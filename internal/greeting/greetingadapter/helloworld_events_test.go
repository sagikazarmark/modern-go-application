package greetingadapter

import (
	"context"
	"testing"

	"github.com/sagikazarmark/modern-go-application/internal/greeting"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHelloWorldEvents_SaidHello(t *testing.T) {
	publisher := &publisherStub{}

	events := NewHelloWorldEvents(publisher)

	event := greeting.SaidHello{
		Message: "Hello, World!",
	}

	err := events.SaidHello(context.Background(), event)
	require.NoError(t, err)

	assert.Equal(t, saidHelloTopic, "said_hello")
	assert.Equal(t, saidHelloTopic, publisher.topic)
	assert.Equal(t, string(publisher.messages[0].Payload), "{\"Message\":\"Hello, World!\"}")
}

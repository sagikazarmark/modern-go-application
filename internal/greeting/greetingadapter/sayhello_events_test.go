package greetingadapter

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sagikazarmark/modern-go-application/internal/greeting"
)

func TestSayHelloEvents_SaidHelloTo(t *testing.T) {
	publisher := &publisherStub{}

	events := NewSayHelloEvents(publisher)

	event := greeting.SaidHello{
		Greeting: "welcome",
		Reply:    "hello",
	}

	err := events.SaidHello(context.Background(), event)
	require.NoError(t, err)

	assert.Equal(t, saidHelloTopic, "said_hello")
	assert.Equal(t, saidHelloTopic, publisher.topic)
	assert.Equal(t, string(publisher.messages[0].Payload), "{\"Greeting\":\"welcome\",\"Reply\":\"hello\"}")
}

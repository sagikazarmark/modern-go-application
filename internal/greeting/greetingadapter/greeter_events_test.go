package greetingadapter

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sagikazarmark/modern-go-application/internal/greeting"
)

func TestGreeterEvents_SaidHelloTo(t *testing.T) {
	publisher := &publisherStub{}

	events := NewGreeterEvents(publisher)

	event := greeting.SaidHello{
		Name:  "welcome",
		Reply: "hello",
	}

	err := events.SaidHello(context.Background(), event)
	require.NoError(t, err)

	assert.Equal(t, saidHelloTopic, "said_hello")
	assert.Equal(t, saidHelloTopic, publisher.topic)
	assert.Equal(t, string(publisher.messages[0].Payload), "{\"Name\":\"welcome\",\"Reply\":\"hello\"}")
}

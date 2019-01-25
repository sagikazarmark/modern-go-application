package greetingdriver

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ThreeDotsLabs/watermill/message/infrastructure/gochannel"
	"github.com/goph/emperror"
	"github.com/goph/logur"
	"github.com/goph/logur/integrations/watermilllog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	api "github.com/sagikazarmark/modern-go-application/.gen/openapi/greeting/go"
	"github.com/sagikazarmark/modern-go-application/internal/greeting"
	"github.com/sagikazarmark/modern-go-application/internal/greeting/greetingadapter"
)

func testSayHello(t *testing.T) {
	sayHello := greeting.NewGreeter(
		greetingadapter.NewGreeterEvents(gochannel.NewGoChannel(
			10,
			watermilllog.New(logur.WithFields(logur.NewNoopLogger(), map[string]interface{}{"component": "watermill"})),
			3*time.Second,
		)),
		greetingadapter.NewNoopLogger(),
		emperror.NewNoopHandler(),
	)
	controller := NewHTTPController(sayHello, emperror.NewNoopHandler())

	server := httptest.NewServer(http.HandlerFunc(controller.SayHello))

	var buf bytes.Buffer

	to := api.HelloRequest{
		Name: "John",
	}

	encoder := json.NewEncoder(&buf)

	err := encoder.Encode(to)
	require.NoError(t, err)

	resp, err := http.Post(server.URL, "application/json", &buf)
	require.NoError(t, err)
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	var hello api.HelloResponse

	err = decoder.Decode(&hello)
	require.NoError(t, err)

	assert.Equal(
		t,
		api.HelloResponse{
			Reply: "hello",
		},
		hello,
	)
}

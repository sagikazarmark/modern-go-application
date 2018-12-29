package greeting_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/goph/emperror"
	"github.com/sagikazarmark/modern-go-application/.gen/openapi/go"
	"github.com/sagikazarmark/modern-go-application/internal/greeting"
	"github.com/sagikazarmark/modern-go-application/internal/greeting/greetingadapter"
	"github.com/sagikazarmark/modern-go-application/internal/greeting/greetingdriver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testSayHello(t *testing.T) {
	events := &sayHelloEventsStub{}

	sayHello := greeting.NewSayHello(events, greetingadapter.NewNoopLogger(), emperror.NewNoopHandler())
	controller := greetingdriver.NewGreetingController(nil, sayHello, emperror.NewNoopHandler())

	server := httptest.NewServer(http.HandlerFunc(controller.SayHello))

	var buf bytes.Buffer

	to := api.HelloRequest{
		Who: "John",
	}

	encoder := json.NewEncoder(&buf)

	err := encoder.Encode(to)
	require.NoError(t, err)

	resp, err := http.Post(server.URL, "application/json", &buf)
	require.NoError(t, err)
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	var hello api.Hello

	err = decoder.Decode(&hello)
	require.NoError(t, err)

	assert.Equal(
		t,
		api.Hello{
			Message: "Hello, John!",
		},
		hello,
	)
}

package greetingdriver

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/goph/emperror"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sagikazarmark/modern-go-application/.gen/openapi/greeting/go"
	"github.com/sagikazarmark/modern-go-application/internal/greeting"
)

func TestHTTPController_SayHello(t *testing.T) {
	service := &helloServiceStub{
		resp: &greeting.HelloResponse{
			Reply: "hello",
		},
	}
	controller := NewHTTPController(service, emperror.NewNoopHandler())

	server := httptest.NewServer(http.HandlerFunc(controller.SayHello))

	var buf bytes.Buffer

	to := api.HelloRequest{
		Greeting: "welcome",
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

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

	api "github.com/sagikazarmark/modern-go-application/.gen/api/openapi/greeting/go"
	"github.com/sagikazarmark/modern-go-application/internal/greeting"
)

func TestHTTPController_SayHello(t *testing.T) {
	greeter := &greeterStub{
		resp: &greeting.HelloResponse{
			Reply: "hello",
		},
	}
	controller := NewHTTPController(greeter, emperror.NewNoopHandler())

	server := httptest.NewServer(http.HandlerFunc(controller.SayHello))

	var buf bytes.Buffer

	apiReq := api.HelloRequest{
		Name: "John",
	}

	encoder := json.NewEncoder(&buf)

	err := encoder.Encode(apiReq)
	require.NoError(t, err)

	resp, err := http.Post(server.URL, "application/json", &buf)
	require.NoError(t, err)
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	var apiResp api.HelloResponse

	err = decoder.Decode(&apiResp)
	require.NoError(t, err)

	assert.Equal(
		t,
		api.HelloResponse{
			Reply: "hello",
		},
		apiResp,
	)
}

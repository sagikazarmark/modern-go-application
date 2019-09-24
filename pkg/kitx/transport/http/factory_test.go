package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	kithttp "github.com/go-kit/kit/transport/http"
)

func TestNewServerFactory(t *testing.T) {
	var beforeCalled bool
	factory := NewServerFactory(
		kithttp.ServerBefore(func(i context.Context, request *http.Request) context.Context {
			beforeCalled = true

			return i
		}),
	)

	var endpointCalled bool
	ep := func(ctx context.Context, request interface{}) (interface{}, error) {
		endpointCalled = true

		return nil, nil
	}

	server := factory.NewServer(ep, kithttp.NopRequestDecoder, NopResponseEncoder)

	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	server.ServeHTTP(httptest.NewRecorder(), req)

	if !beforeCalled {
		t.Error("global before function is supposed to be called")
	}

	if !endpointCalled {
		t.Error("endpoint is supposed to be called")
	}
}

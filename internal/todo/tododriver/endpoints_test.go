package tododriver

import (
	"context"
	"testing"

	"github.com/go-kit/kit/endpoint"
	"github.com/stretchr/testify/require"
)

func TestCreateEndpoint(t *testing.T) {
	tests := []endpoint.Endpoint{
		MakeCreateEndpoint(),
		MakeEndpoints().Create,
	}

	for _, test := range tests {
		test := test

		t.Run("", func(t *testing.T) {
			req := createTodoRequest{
				Text: "My first todo",
			}

			resp, err := test(context.Background(), req)
			require.NoError(t, err)
		})
	}
}

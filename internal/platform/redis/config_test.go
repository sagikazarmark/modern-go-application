package redis

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_Validate(t *testing.T) {
	tests := map[string]Config{
		"redis host is required": {
			Port: 6379,
		},
		"redis port is required": {
			Host: "127.0.0.1",
		},
	}

	for name, test := range tests {
		name, test := name, test

		t.Run(name, func(t *testing.T) {
			err := test.Validate()

			assert.EqualError(t, err, name)
		})
	}
}

func TestConfig_Server(t *testing.T) {
	config := Config{
		Host: "127.0.0.1",
		Port: 6379,
	}

	server := config.Server()

	assert.Equal(t, "127.0.0.1:6379", server)
}

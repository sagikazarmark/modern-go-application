package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogConfig_Validate(t *testing.T) {
	tests := map[string]LogConfig{
		"log format is required": {},
		"invalid log format: xml": {
			Format: "xml",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := test.Validate()

			assert.EqualError(t, err, name)
		})
	}
}

package main

import (
	"bytes"
	"testing"
	"time"

	"github.com/goph/conf"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//nolint:deadcode
func TestConfig_Validate(t *testing.T) {
	tests := map[string]Config{
		"environment is required": {
			ShutdownTimeout: 15 * time.Second,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := test.Validate()

			assert.EqualError(t, err, name)
		})
	}
}

//nolint:deadcode
func TestConfig_Prepare(t *testing.T) {
	config := NewConfig()

	var buf bytes.Buffer

	configurator := conf.NewConfigurator("app", conf.ContinueOnError)
	configurator.SetOutput(&buf)

	config.Prepare(configurator)

	environment := map[string]string{
		"ENVIRONMENT": "staging",
		"DEBUG":       "true",
		"LOG_FORMAT":  "logfmt",
	}

	args := []string{"--shutdown-timeout", "5s"}

	err := configurator.Parse(args, environment)
	require.NoError(t, err)
}

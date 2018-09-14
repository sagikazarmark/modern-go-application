package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDatabaseConfig_Validate(t *testing.T) {
	tests := map[string]DatabaseConfig{
		"database host is required": {
			Port: 3306,
			User: "root",
			Pass: "",
			Name: "database",
		},
		"database port is required": {
			Host: "localhost",
			User: "root",
			Pass: "",
			Name: "database",
		},
		"database user is required": {
			Host: "localhost",
			Port: 3306,
			Pass: "",
			Name: "database",
		},
		"database name is required": {
			Host: "localhost",
			Port: 3306,
			User: "root",
			Pass: "",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := test.Validate()

			assert.EqualError(t, err, name)
		})
	}
}

func TestDatabaseConfig_DSN(t *testing.T) {
	config := DatabaseConfig{
		Host: "host",
		Port: 3306,
		User: "root",
		Pass: "",
		Name: "database",
		Params: map[string]string{
			"parseTime": "true",
		},
	}

	dsn := config.DSN()

	assert.Equal(t, "root:@tcp(host:3306)/database?parseTime=true", dsn)
}

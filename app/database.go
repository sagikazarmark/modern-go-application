package app

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql" // Importing mysql driver here
	"github.com/pkg/errors"
)

// DatabaseConfig holds information necessary for connecting to a database.
type DatabaseConfig struct {
	Host string
	Port int
	User string
	Pass string
	Name string

	Params map[string]string
}

// Validate checks that the configuration is valid.
func (c DatabaseConfig) Validate() error {
	if c.Host == "" {
		return errors.New("database host is required")
	}

	if c.Port == 0 {
		return errors.New("database port is required")
	}

	if c.User == "" {
		return errors.New("database user is required")
	}

	if c.Name == "" {
		return errors.New("database name is required")
	}

	return nil
}

// DSN returns a MySQL driver compatible data source name.
func (c DatabaseConfig) DSN() string {
	var params string

	if len(c.Params) > 0 {
		var query string

		for key, value := range c.Params {
			if query != "" {
				query += "&"
			}

			query += key + "=" + value
		}

		params = "?" + query
	}

	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s%s",
		c.User,
		c.Pass,
		c.Host,
		c.Port,
		c.Name,
		params,
	)
}

// NewDatabaseConfig returns a new DatabaseConfig instance with some defaults.
func NewDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		Host: "localhost",
		Port: 3306,
		User: "root",
		Params: map[string]string{
			"charset": "utf8mb4",
		},
	}
}

// NewDatabaseConnection returns a new database connection for the application.
func NewDatabaseConnection(config DatabaseConfig) (*sql.DB, error) {
	// Set some mandatory parameters
	config.Params["parseTime"] = "true"
	config.Params["rejectReadOnly"] = "true"

	db, err := sql.Open("mysql", config.DSN())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	db.SetMaxIdleConns(2)
	db.SetMaxOpenConns(10)
	db.SetConnMaxLifetime(time.Minute)

	return db, nil
}

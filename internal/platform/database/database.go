package database

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql" // Importing mysql driver here
	"github.com/pkg/errors"
)

// NewConnection returns a new database connection for the application.
func NewConnection(config Config) (*sql.DB, error) {
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

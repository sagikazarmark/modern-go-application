package database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql" // Importing mysql driver here
	"github.com/opencensus-integrations/ocsql"
	"github.com/pkg/errors"
)

// NewConnection returns a new database connection for the application.
func NewConnection(config Config) (*sql.DB, error) {
	driverName, err := ocsql.Register(
		"mysql",
		ocsql.WithOptions(ocsql.TraceOptions{
			AllowRoot:    false,
			Ping:         true,
			RowsNext:     true,
			RowsClose:    true,
			RowsAffected: true,
			LastInsertID: true,
			Query:        true,
			QueryParams:  false,
		}),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to register ocsql driver")
	}

	// Set some mandatory parameters
	config.Params["parseTime"] = "true"
	config.Params["rejectReadOnly"] = "true"

	db, err := sql.Open(driverName, config.DSN())

	return db, errors.WithStack(err)
}

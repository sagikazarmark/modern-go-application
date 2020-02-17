package database

import (
	"database/sql/driver"

	"contrib.go.opencensus.io/integrations/ocsql"
	"emperror.dev/errors"
	"github.com/go-sql-driver/mysql"
)

// NewConnector returns a new database connector for the application.
func NewConnector(config Config) (driver.Connector, error) {
	// Set some mandatory parameters
	config.Params["parseTime"] = "true"
	config.Params["rejectReadOnly"] = "true"

	// TODO: fill in the config instead of playing with DSN
	conf, err := mysql.ParseDSN(config.DSN())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	connector, err := mysql.NewConnector(conf)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return ocsql.WrapConnector(
		connector,
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
	), nil
}

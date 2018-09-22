package health

import (
	"database/sql"
	"time"

	"github.com/InVisionApp/go-health"
	"github.com/InVisionApp/go-health/checkers"
	"github.com/pkg/errors"
)

// New returns a new health checker.
func New(db *sql.DB) (*health.Health, error) {
	h := health.New()

	h.DisableLogging() // TODO: implement a logger

	dbCheck, err := checkers.NewSQL(&checkers.SQLConfig{
		Pinger: db,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create db health checker")
	}

	h.AddCheck(&health.Config{
		Name: "database",
		Checker: dbCheck,
		Interval: time.Duration(3) * time.Second,
		Fatal:    true,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to add health checker")
	}

	return h, nil
}

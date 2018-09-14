package main

import "github.com/sagikazarmark/go-service-project-boilerplate/app"

// NewAppContext returns a new application context.
func NewAppContext(config Config) app.Context {
	return app.Context{
		Build: app.Build{
			Version:    Version,
			CommitHash: CommitHash,
			Date:       BuildDate,
		},
		Name:         app.ServiceName,
		FriendlyName: app.FriendlyServiceName,

		Environment: config.Environment,
		Debug:       config.Debug,
	}
}

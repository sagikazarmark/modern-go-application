package main

const (
	// serviceName is an identifier-like name used anywhere this app needs to be identified.
	//
	// It identifies the service itself, the actual instance needs to be identified via environment
	// and other details.
	serviceName = "mga"

	// friendlyServiceName is the visible name of the service.
	friendlyServiceName = "Modern Go Application"

	// envPrefix is prepended to environment variables when processing configuration.
	envPrefix = "app"
)

package main

const (
	// ServiceName is an identifier-like name used anywhere this app needs to be identified.
	//
	// It identifies the service itself, the actual instance needs to be identified via environment
	// and other details.
	ServiceName = "mga"

	// FriendlyServiceName is the visible name of the service.
	FriendlyServiceName = "Modern Go Application"

	// EnvPrefix is prepended to environment variables when processing configuration.
	EnvPrefix = "app"
)

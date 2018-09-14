package app

// Context stores build and runtime application information.
type Context struct {
	// Build information
	Build Build

	// Name used for identification
	Name string

	// FriendlyName appearing in some logs
	FriendlyName string

	// Environment information
	Environment string
	Debug       bool
}

// Build stores build information about the application.
type Build struct {
	Version    string
	CommitHash string
	Date       string
}

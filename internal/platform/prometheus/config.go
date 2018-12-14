package prometheus

// Config holds information for configuring the Prometheus exporter
type Config struct {
	Namespace string
}

// Validate checks that the configuration is valid.
func (c Config) Validate() error {
	return nil
}

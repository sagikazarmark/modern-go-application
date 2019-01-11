package runner

// RunCloser is a runnable construct that can be stopped by calling close on it.
type RunCloser interface {
	// Run starts the process.
	Run() error

	// Close stops the process.
	Close() error
}

// RunCloserRunner is a runner that makes a RunCloser runner compatible.
type RunCloserRunner struct {
	RunCloser RunCloser
}

// Start starts the RunCloser and waits for it to return.
func (r *RunCloserRunner) Start() error {
	return r.RunCloser.Run()
}

// Stop stops the RunCloser.
func (r *RunCloserRunner) Stop(e error) {
	_ = r.RunCloser.Close()
}

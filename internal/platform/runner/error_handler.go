package runner

// ErrorHandler is responsible for handling an error.
type ErrorHandler interface {
	// Handle takes care of unhandled errors.
	Handle(err error)
}

// nolint: gochecknoglobals
var defaultErrorHandler ErrorHandler = &noopHandler{}

type noopHandler struct{}

func (*noopHandler) Handle(err error) {}

package greeting

// ErrorHandler handles an error.
type ErrorHandler interface {
	// Handle handles an error.
	Handle(err error)
}

package todo

// ErrorHandler handles an error.
type ErrorHandler interface {
	// Handle handles an error.
	Handle(err error)
}

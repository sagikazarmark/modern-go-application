package todo

import "context"

// TodoEvents is the dispatcher for todo events.
type TodoEvents interface {
	// MarkedAsDone dispatches a MarkedAsDone event.
	MarkedAsDone(ctx context.Context, event MarkedAsDone) error
}

// MarkedAsDone event is triggered when a todo gets marked as done.
type MarkedAsDone struct {
	ID string
}

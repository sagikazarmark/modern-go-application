package todo

import (
	"context"

	"github.com/pkg/errors"
)

// Todo is a note describing a task to be done.
type Todo struct {
	ID   string
	Text string
	Done bool
}

// List manages a list of todos.
type List struct {
	idgenerator IDGenerator
	store       Store
	events      Events
}

// IDGenerator generates a new ID.
type IDGenerator interface {
	// Generate generates a new ID.
	Generate() (string, error)
}

// Store stores todos.
type Store interface {
	// Store stores a todo.
	Store(ctx context.Context, todo Todo) error

	// All returns all todos.
	All(ctx context.Context) ([]Todo, error)

	// Get returns a single todo by its ID.
	Get(ctx context.Context, id string) (Todo, error)
}

// NotFoundError is returned if a todo cannot be found.
type NotFoundError struct {
	ID string
}

// Error implements the error interface.
func (NotFoundError) Error() string {
	return "todo not found"
}

// Context returns context parameters for the error.
func (e NotFoundError) Context() []interface{} {
	return []interface{}{"todo_id", e.ID}
}

// Events dispatches todo events.
type Events interface {
	// MarkedAsDone dispatches a MarkedAsDone event.
	MarkedAsDone(ctx context.Context, event MarkedAsDone) error
}

// MarkedAsDone event is triggered when a todo gets marked as done.
type MarkedAsDone struct {
	ID string
}

// NewList returns a new todo list.
func NewList(id IDGenerator, todos Store, events Events) *List {
	return &List{
		idgenerator: id,
		store:       todos,
		events:      events,
	}
}

// CreateTodo adds a new todo to the list.
func (l *List) CreateTodo(ctx context.Context, text string) (string, error) {
	id, err := l.idgenerator.Generate()
	if err != nil {
		return "", err
	}

	todo := Todo{
		ID:   id,
		Text: text,
	}

	err = l.store.Store(ctx, todo)

	return id, err
}

// ListTodos returns the list of todos.
func (l *List) ListTodos(ctx context.Context) ([]Todo, error) {
	return l.store.All(ctx)
}

// MarkAsDone marks a todo as done.
func (l *List) MarkAsDone(ctx context.Context, id string) error {
	todo, err := l.store.Get(ctx, id)
	if err != nil {
		return errors.WithMessage(err, "failed to mark todo as done")
	}

	todo.Done = true

	err = l.store.Store(ctx, todo)
	if err != nil {
		return errors.WithMessage(err, "failed to mark todo as done")
	}

	event := MarkedAsDone{
		ID: todo.ID,
	}

	err = l.events.MarkedAsDone(ctx, event)
	if err != nil {
		return errors.WithMessage(err, "failed to mark todo as done")
	}

	return nil
}

package todo

import (
	"context"

	"emperror.dev/errors"
)

// Todo is a note describing a task to be done.
type Todo struct {
	ID   string
	Text string
	Done bool
}

// Service manages a list of todos.
type Service interface {
	// CreateTodo adds a new todo to the todo list.
	CreateTodo(ctx context.Context, text string) (string, error)

	// ListTodos returns the list of todos.
	ListTodos(ctx context.Context) ([]Todo, error)

	// MarkAsDone marks a todo as done.
	MarkAsDone(ctx context.Context, id string) error
}

type service struct {
	idgenerator IDGenerator
	store       Store
	events      Events
}

// IDGenerator generates a new ID.
type IDGenerator interface {
	// Generate generates a new ID.
	Generate() (string, error)
}

// Store provides todo persistence.
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

// Details returns error details.
func (e NotFoundError) Details() []interface{} {
	return []interface{}{"todo_id", e.ID}
}

// IsBusinessError tells the transport layer whether this error should be translated into the transport format
// or an internal error should be returned instead.
func (NotFoundError) IsBusinessError() bool {
	return true
}

//go:generate mga gen kit endpoint --outdir tododriver --with-oc Service
//go:generate mga gen ev dispatcher Events
//go:generate mga gen ev handler MarkedAsDone

// Events dispatches todo events.
type Events interface {
	// MarkedAsDone dispatches a MarkedAsDone event.
	MarkedAsDone(ctx context.Context, event MarkedAsDone) error
}

// MarkedAsDone event is triggered when a todo gets marked as done.
type MarkedAsDone struct {
	ID string
}

// NewService returns a new Service.
func NewService(idgenerator IDGenerator, store Store, events Events) Service {
	return &service{
		idgenerator: idgenerator,
		store:       store,
		events:      events,
	}
}

func (s service) CreateTodo(ctx context.Context, text string) (string, error) {
	id, err := s.idgenerator.Generate()
	if err != nil {
		return "", err
	}

	todo := Todo{
		ID:   id,
		Text: text,
	}

	err = s.store.Store(ctx, todo)

	return id, err
}

func (s service) ListTodos(ctx context.Context) ([]Todo, error) {
	return s.store.All(ctx)
}

func (s service) MarkAsDone(ctx context.Context, id string) error {
	todo, err := s.store.Get(ctx, id)
	if err != nil {
		return errors.WithMessage(err, "failed to mark todo as done")
	}

	todo.Done = true

	err = s.store.Store(ctx, todo)
	if err != nil {
		return errors.WithMessage(err, "failed to mark todo as done")
	}

	event := MarkedAsDone{
		ID: todo.ID,
	}

	err = s.events.MarkedAsDone(ctx, event)
	if err != nil {
		return errors.WithMessage(err, "failed to mark todo as done")
	}

	return nil
}

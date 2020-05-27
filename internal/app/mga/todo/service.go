package todo

import (
	"context"

	"emperror.dev/errors"
)

// Todo is a note describing a task to be done.
type Todo struct {
	ID        string
	Title     string
	Completed bool
}

// +kit:endpoint:errorStrategy=service

// Service manages a list of todos.
type Service interface {
	// CreateTodo adds a new todo to the todo list.
	CreateTodo(ctx context.Context, title string) (todo Todo, err error)

	// ListTodos returns the list of todos.
	ListTodos(ctx context.Context) (todos []Todo, err error)

	// GetItem returns the details of an item.
	GetItem(ctx context.Context, id string) (todo Todo, err error)

	// MarkAsComplete marks an item as complete.
	MarkAsComplete(ctx context.Context, id string) error

	// DeleteAll deletes all items from the list.
	DeleteAll(ctx context.Context) error
}

// NewService returns a new Service.
func NewService(idgenerator IDGenerator, store Store, events Events) Service {
	return &service{
		idgenerator: idgenerator,
		store:       store,
		events:      events,
	}
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
	// Store stores an item.
	Store(ctx context.Context, todo Todo) error

	// All returns all items.
	All(ctx context.Context) ([]Todo, error)

	// Get returns a single item by its ID.
	Get(ctx context.Context, id string) (Todo, error)

	// DeleteAll deletes all items in the store.
	DeleteAll(ctx context.Context) error
}

// NotFoundError is returned if an item cannot be found.
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

// NotFound tells a client that this error is related to a resource being not found.
// Can be used to translate the error to eg. status code.
func (NotFoundError) NotFound() bool {
	return true
}

// ServiceError tells the transport layer whether this error should be translated into the transport format
// or an internal error should be returned instead.
func (NotFoundError) ServiceError() bool {
	return true
}

// +mga:event:dispatcher

// Events dispatches todo events.
type Events interface {
	// MarkedAsComplete dispatches a MarkedAsComplete event.
	MarkedAsComplete(ctx context.Context, event MarkedAsComplete) error
}

// +mga:event:handler

// MarkedAsComplete event is triggered when a todo gets marked as complete.
type MarkedAsComplete struct {
	ID string
}

type validationError struct {
	violations map[string][]string
}

func (validationError) Error() string {
	return "invalid todo"
}

func (e validationError) Violations() map[string][]string {
	return e.violations
}

// Validation tells a client that this error is related to a resource being invalid.
// Can be used to translate the error to eg. status code.
func (validationError) Validation() bool {
	return true
}

// ServiceError tells the transport layer whether this error should be translated into the transport format
// or an internal error should be returned instead.
func (validationError) ServiceError() bool {
	return true
}

func (s service) CreateTodo(ctx context.Context, text string) (Todo, error) {
	id, err := s.idgenerator.Generate()
	if err != nil {
		return Todo{}, err
	}

	if text == "" {
		return Todo{}, errors.WithStack(validationError{violations: map[string][]string{
			"text": {
				"text cannot be empty",
			},
		}})
	}

	todo := Todo{
		ID:    id,
		Title: text,
	}

	err = s.store.Store(ctx, todo)
	if err != nil {
		return Todo{}, err
	}

	return todo, nil
}

func (s service) ListTodos(ctx context.Context) ([]Todo, error) {
	return s.store.All(ctx)
}

func (s service) GetItem(ctx context.Context, id string) (Todo, error) {
	todo, err := s.store.Get(ctx, id)
	if err != nil {
		return Todo{}, err
	}

	return todo, nil
}

func (s service) MarkAsComplete(ctx context.Context, id string) error {
	todo, err := s.store.Get(ctx, id)
	if err != nil {
		return errors.WithMessage(err, "failed to mark todo as complete")
	}

	todo.Completed = true

	err = s.store.Store(ctx, todo)
	if err != nil {
		return errors.WithMessage(err, "failed to mark todo as complete")
	}

	event := MarkedAsComplete{
		ID: todo.ID,
	}

	err = s.events.MarkedAsComplete(ctx, event)
	if err != nil {
		return errors.WithMessage(err, "failed to mark todo as complete")
	}

	return nil
}

func (s service) DeleteAll(ctx context.Context) error {
	return s.store.DeleteAll(ctx)
}

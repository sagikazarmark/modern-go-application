package todo

import (
	"context"

	"github.com/pkg/errors"
)

// Todo is a note describing a task to be executed.
type Todo struct {
	ID   string
	Text string
	Done bool
}

// TodoList manages a list of todos.
type TodoList struct {
	id     IDGenerator
	todos  TodoStore
	events TodoEvents
}

// IDGenerator generates a new ID.
type IDGenerator interface {
	// New generates a new ID.
	New() (string, error)
}

// TodoStore stores the items on the todo list.
type TodoStore interface {
	// Store stores a todo.
	Store(ctx context.Context, todo Todo) error

	// All returns all todos.
	All(ctx context.Context) ([]Todo, error)

	// Get returns a single todo by its ID.
	Get(ctx context.Context, id string) (Todo, error)
}

// TodoNotFoundError is returned from the TodoStore if a Todo cannot be found.
type TodoNotFoundError struct {
	ID string
}

// Error implements the error interface.
func (TodoNotFoundError) Error() string {
	return "todo not found"
}

// Context returns context parameters for the error.
func (e TodoNotFoundError) Context() []interface{} {
	return []interface{}{"todo_id", e.ID}
}

// TodoEvents is the dispatcher for todo events.
type TodoEvents interface {
	// MarkedAsDone dispatches a MarkedAsDone event.
	MarkedAsDone(ctx context.Context, event MarkedAsDone) error
}

// MarkedAsDone event is triggered when a todo gets marked as done.
type MarkedAsDone struct {
	ID string
}

// NewTodoList returns a new TodoList instance.
func NewTodoList(id IDGenerator, todos TodoStore, events TodoEvents) *TodoList {
	return &TodoList{
		id:     id,
		todos:  todos,
		events: events,
	}
}

// CreateTodo adds a new todo to the todo list.
func (t *TodoList) CreateTodo(ctx context.Context, text string) (string, error) {
	id, err := t.id.New()
	if err != nil {
		return "", err
	}

	todo := Todo{
		ID:   id,
		Text: text,
	}

	err = t.todos.Store(ctx, todo)
	if err != nil {
		return "", err
	}

	return id, nil
}

// ListTodos returns the list of todos.
func (t *TodoList) ListTodos(ctx context.Context) ([]Todo, error) {
	todos, err := t.todos.All(ctx)
	if err != nil {
		return nil, err
	}

	return todos, nil
}

// MarkAsDone marks a todo as done.
func (t *TodoList) MarkAsDone(ctx context.Context, id string) error {
	todo, err := t.todos.Get(ctx, id)
	if err != nil {
		return errors.WithMessage(err, "failed to mark todo as done")
	}

	todo.Done = true

	err = t.todos.Store(ctx, todo)
	if err != nil {
		return errors.WithMessage(err, "failed to mark todo as done")
	}

	event := MarkedAsDone{
		ID: todo.ID,
	}

	err = t.events.MarkedAsDone(ctx, event)
	if err != nil {
		return errors.WithMessage(err, "failed to mark todo as done")
	}

	return nil
}

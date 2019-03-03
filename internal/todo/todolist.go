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
	id    IDGenerator
	todos TodoStore
	events TodoEvents
}

// NewTodoList returns a new TodoList instance.
func NewTodoList(id IDGenerator, todos TodoStore, events TodoEvents) *TodoList {
	return &TodoList{
		id:    id,
		todos: todos,
		events: events,
	}
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

// CreateTodoRequest contains a new todo.
type CreateTodoRequest struct {
	Text string
}

// CreateTodoResponse contains the ID of the newly created todo.
type CreateTodoResponse struct {
	ID string
}

// CreateTodo adds a new todo to the todo list.
func (t *TodoList) CreateTodo(ctx context.Context, req CreateTodoRequest) (*CreateTodoResponse, error) {
	id, err := t.id.New()
	if err != nil {
		return nil, err
	}

	todo := Todo{
		ID:   id,
		Text: req.Text,
	}

	err = t.todos.Store(ctx, todo)
	if err != nil {
		return nil, err
	}

	return &CreateTodoResponse{
		ID: id,
	}, nil
}

// ListTodosResponse contains the list of todos on the list.
type ListTodosResponse struct {
	Todos []Todo
}

// ListTodos returns the list of todos on the list.
func (t *TodoList) ListTodos(ctx context.Context) (*ListTodosResponse, error) {
	todos, err := t.todos.All(ctx)
	if err != nil {
		return nil, err
	}

	return &ListTodosResponse{
		Todos: todos,
	}, nil
}

type MarkAsDoneRequest struct {
	ID string
}

func (t *TodoList) MarkAsDone(ctx context.Context, req MarkAsDoneRequest) error {
	todo, err := t.todos.Get(ctx, req.ID)
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

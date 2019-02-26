package todo

import "context"

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
}

// NewTodoList returns a new TodoList instance.
func NewTodoList(id IDGenerator, todos TodoStore) *TodoList {
	return &TodoList{
		id:    id,
		todos: todos,
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
	Store(todo Todo) error
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

	err = t.todos.Store(todo)
	if err != nil {
		return nil, err
	}

	return &CreateTodoResponse{
		ID: id,
	}, nil
}

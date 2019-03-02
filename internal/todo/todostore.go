package todo

import (
	"context"
	"github.com/goph/emperror"
	"github.com/pkg/errors"
	"sort"
)

// InmemoryTodoStore keeps todos in the memory.
// Use it in tests or development/demo purposes.
type InmemoryTodoStore struct {
	todos map[string]Todo
}

// NewInmemoryTodoStore returns a new inmemory todo store.
func NewInmemoryTodoStore() *InmemoryTodoStore {
	return &InmemoryTodoStore{todos: make(map[string]Todo)}
}

// Store stores a todo.
func (s *InmemoryTodoStore) Store(ctx context.Context, todo Todo) error {
	s.todos[todo.ID] = todo

	return nil
}

// All returns all todos.
func (s *InmemoryTodoStore) All(ctx context.Context) ([]Todo, error) {
	todos := make([]Todo, len(s.todos))

	// This makes sure todos are always returned in the same, sorted representation
	var keys []string
	for k := range s.todos {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	i := 0
	for _, key := range keys {
		todos[i] = s.todos[key]

		i++
	}

	return todos, nil
}

// Get returns a single todo by its ID.
func (s *InmemoryTodoStore) Get(ctx context.Context, id string) (Todo, error) {
	todo, ok := s.todos[id]
	if !ok {
		return todo, TodoNotFoundError{ID: id}
	}

	return todo, nil
}

// ReadOnlyTodoStore cannot be modified.
type ReadOnlyTodoStore struct {
	todos TodoStore
}

// NewReadOnlyTodoStore returns a new read-only todo store instance.
func NewReadOnlyTodoStore(todos TodoStore) *ReadOnlyTodoStore {
	return &ReadOnlyTodoStore{todos: todos}
}

// Store stores a todo.
func (*ReadOnlyTodoStore) Store(ctx context.Context, todo Todo) error {
	return emperror.With(errors.New("read-only todo store cannot be modified"), "todo_id", todo.ID)
}

// All returns all todos.
func (s *ReadOnlyTodoStore) All(ctx context.Context) ([]Todo, error) {
	return s.todos.All(ctx)
}

// Get returns a single todo by its ID.
func (s *ReadOnlyTodoStore) Get(ctx context.Context, id string) (Todo, error) {
	return s.todos.Get(ctx, id)
}

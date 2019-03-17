package todo

import (
	"context"
	"sort"

	"github.com/goph/emperror"
	"github.com/pkg/errors"
)

// InmemoryStore keeps todos in the memory.
// Use it in tests or development/demo purposes.
type InmemoryStore struct {
	todos map[string]Todo
}

// NewInmemoryStore returns a new inmemory todo store.
func NewInmemoryStore() *InmemoryStore {
	return &InmemoryStore{
		todos: make(map[string]Todo),
	}
}

// Store stores a todo.
func (s *InmemoryStore) Store(ctx context.Context, todo Todo) error {
	s.todos[todo.ID] = todo

	return nil
}

// All returns all todos.
func (s *InmemoryStore) All(ctx context.Context) ([]Todo, error) {
	todos := make([]Todo, len(s.todos))

	// This makes sure todos are always returned in the same, sorted representation
	keys := make([]string, len(s.todos))
	i := 0
	for k := range s.todos {
		keys[i] = k
		i++
	}
	sort.Strings(keys)

	for i, key := range keys {
		todos[i] = s.todos[key]
	}

	return todos, nil
}

// Get returns a single todo by its ID.
func (s *InmemoryStore) Get(ctx context.Context, id string) (Todo, error) {
	todo, ok := s.todos[id]
	if !ok {
		return todo, NotFoundError{ID: id}
	}

	return todo, nil
}

// ReadOnlyStore cannot be modified.
type ReadOnlyStore struct {
	store Store
}

// NewReadOnlyStore returns a new read-only todo store instance.
func NewReadOnlyStore(store Store) *ReadOnlyStore {
	return &ReadOnlyStore{
		store: store,
	}
}

// Store stores a todo.
func (*ReadOnlyStore) Store(ctx context.Context, todo Todo) error {
	return emperror.With(errors.New("read-only todo store cannot be modified"), "todo_id", todo.ID)
}

// All returns all todos.
func (s *ReadOnlyStore) All(ctx context.Context) ([]Todo, error) {
	return s.store.All(ctx)
}

// Get returns a single todo by its ID.
func (s *ReadOnlyStore) Get(ctx context.Context, id string) (Todo, error) {
	return s.store.Get(ctx, id)
}

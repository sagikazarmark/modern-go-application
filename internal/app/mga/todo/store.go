package todo

import (
	"context"
	"sort"

	"emperror.dev/errors"
)

// InMemoryStore keeps todos in the memory.
// Use it in tests or for development/demo purposes.
type InMemoryStore struct {
	todos map[string]Todo
}

// NewInMemoryStore returns a new inmemory todo store.
func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		todos: make(map[string]Todo),
	}
}

// Store stores a todo.
func (s *InMemoryStore) Store(ctx context.Context, todo Todo) error {
	s.todos[todo.ID] = todo

	return nil
}

// All returns all todos.
func (s *InMemoryStore) All(ctx context.Context) ([]Todo, error) {
	todos := make([]Todo, len(s.todos))

	// This makes sure todos are always returned in the same, sorted order
	keys := make([]string, 0, len(s.todos))
	for k := range s.todos {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for i, key := range keys {
		todos[i] = s.todos[key]
	}

	return todos, nil
}

// Get returns a single todo by its ID.
func (s *InMemoryStore) Get(ctx context.Context, id string) (Todo, error) {
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
	return errors.NewWithDetails("read-only todo store cannot be modified", "todo_id", todo.ID)
}

// All returns all todos.
func (s *ReadOnlyStore) All(ctx context.Context) ([]Todo, error) {
	return s.store.All(ctx)
}

// Get returns a single todo by its ID.
func (s *ReadOnlyStore) Get(ctx context.Context, id string) (Todo, error) {
	return s.store.Get(ctx, id)
}

package todo

import (
	"context"
	"sort"
)

// InmemoryTodoStore keeps todos in the memory.
// Use it in tests or development/demo purposes.
type InmemoryTodoStore struct {
	todos map[string]Todo
}

// NewInmemoryTodoStore returns a new inmemory todo store.
func NewInmemoryTodoStore() *InmemoryTodoStore {
	return &InmemoryTodoStore{
		todos: make(map[string]Todo),
	}
}

func (s *InmemoryTodoStore) Store(todo Todo) error {
	s.todos[todo.ID] = todo

	return nil
}

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

func (s *InmemoryTodoStore) Get(ctx context.Context, id string) (Todo, error) {
	todo, ok := s.todos[id]
	if !ok {
		return todo, TodoNotFoundError{ID: id}
	}

	return todo, nil
}

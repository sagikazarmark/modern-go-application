package todo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInMemoryStore_StoresATodo(t *testing.T) {
	store := NewInMemoryStore()

	todo := Item{
		ID:    "id",
		Title: "Store me!",
	}

	err := store.Store(context.Background(), todo)
	require.NoError(t, err)

	assert.Equal(t, todo, store.todos[todo.ID])
}

func TestInMemoryStore_OverwritesAnExistingTodo(t *testing.T) {
	store := NewInMemoryStore()

	todo := Item{
		ID:        "id",
		Title:     "Store me first!",
		Completed: true,
	}

	err := store.Store(context.Background(), todo)
	require.NoError(t, err)

	todo = Item{
		ID:    "id",
		Title: "Store me!",
	}

	err = store.Store(context.Background(), todo)
	require.NoError(t, err)

	assert.Equal(t, todo, store.todos[todo.ID])
}

func TestInMemoryStore_ListsAllTodos(t *testing.T) {
	store := NewInMemoryStore()

	store.todos["id"] = Item{
		ID:        "id",
		Title:     "Store me first!",
		Completed: true,
	}

	store.todos["id2"] = Item{
		ID:        "id2",
		Title:     "Store me second!",
		Completed: true,
	}

	todos, err := store.All(context.Background())
	require.NoError(t, err)

	expectedTodos := []Item{store.todos["id"], store.todos["id2"]}

	assert.Equal(t, expectedTodos, todos)
}

func TestInMemoryStore_GetsATodo(t *testing.T) {
	store := NewInMemoryStore()

	id := "id"

	store.todos[id] = Item{
		ID:    id,
		Title: "Store me!",
	}

	todo, err := store.Get(context.Background(), id)
	require.NoError(t, err)

	assert.Equal(t, store.todos[id], todo)
}

func TestInMemoryStore_CannotReturnANonExistingTodo(t *testing.T) {
	store := NewInMemoryStore()

	_, err := store.Get(context.Background(), "id")
	require.Error(t, err)

	require.IsType(t, NotFoundError{}, err)

	e := err.(NotFoundError)
	assert.Equal(t, "id", e.ID)
}

func TestReadOnlyStore_IsReadOnly(t *testing.T) {
	todo := Item{
		ID:    "id",
		Title: "Store me!",
	}

	store := NewReadOnlyStore(NewInMemoryStore())

	err := store.Store(context.Background(), todo)
	require.Error(t, err)
}

func TestReadOnlyStore_ListsAllTodos(t *testing.T) {
	inmemStore := NewInMemoryStore()
	store := NewReadOnlyStore(inmemStore)

	inmemStore.todos["id"] = Item{
		ID:        "id",
		Title:     "Store me first!",
		Completed: true,
	}

	inmemStore.todos["id2"] = Item{
		ID:        "id2",
		Title:     "Store me second!",
		Completed: true,
	}

	todos, err := store.All(context.Background())
	require.NoError(t, err)

	expectedTodos := []Item{inmemStore.todos["id"], inmemStore.todos["id2"]}

	assert.Equal(t, expectedTodos, todos)
}

func TestReadOnlyStore_GetsATodo(t *testing.T) {
	inmemStore := NewInMemoryStore()
	store := NewReadOnlyStore(inmemStore)

	id := "id"

	inmemStore.todos[id] = Item{
		ID:    id,
		Title: "Store me!",
	}

	todo, err := store.Get(context.Background(), id)
	require.NoError(t, err)

	assert.Equal(t, inmemStore.todos[id], todo)
}

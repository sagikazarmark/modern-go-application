package todo

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestInmemoryTodoStore_StoresATodo(t *testing.T) {
	store := NewInmemoryTodoStore()

	todo := Todo{
		ID:   "id",
		Text: "Store me!",
	}

	err := store.Store(context.Background(), todo)
	require.NoError(t, err)

	assert.Equal(t, todo, store.todos[todo.ID])
}

func TestInmemoryTodoStore_OverwritesAnExistingTodo(t *testing.T) {
	store := NewInmemoryTodoStore()

	todo := Todo{
		ID:   "id",
		Text: "Store me first!",
		Done: true,
	}

	err := store.Store(context.Background(), todo)
	require.NoError(t, err)

	todo = Todo{
		ID:   "id",
		Text: "Store me!",
	}

	err = store.Store(context.Background(), todo)
	require.NoError(t, err)

	assert.Equal(t, todo, store.todos[todo.ID])
}

func TestInmemoryTodoStore_ListsAllTodos(t *testing.T) {
	store := NewInmemoryTodoStore()

	store.todos["id"] = Todo{
		ID:   "id",
		Text: "Store me first!",
		Done: true,
	}

	store.todos["id2"] = Todo{
		ID:   "id2",
		Text: "Store me second!",
		Done: true,
	}

	todos, err := store.All(context.Background())
	require.NoError(t, err)

	expectedTodos := []Todo{store.todos["id"], store.todos["id2"]}

	assert.Equal(t, expectedTodos, todos)
}

func TestInmemoryTodoStore_GetsATodo(t *testing.T) {
	store := NewInmemoryTodoStore()

	id := "id"

	store.todos[id] = Todo{
		ID:   id,
		Text: "Store me!",
	}

	todo, err := store.Get(context.Background(), id)
	require.NoError(t, err)

	assert.Equal(t, store.todos[id], todo)
}

func TestInmemoryTodoStore_ReturnsANotFoundErrorWhenATodoIsNotFound(t *testing.T) {
	store := NewInmemoryTodoStore()

	_, err := store.Get(context.Background(), "id")

	require.IsType(t, TodoNotFoundError{}, err)

	e := err.(TodoNotFoundError)
	assert.Equal(t, "id", e.ID)
}

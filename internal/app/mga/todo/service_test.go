package todo

import (
	"context"
	"testing"

	"github.com/goph/idgen"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestList_CreatesATodo(t *testing.T) {
	todoStore := NewInMemoryStore()

	const expectedID = "id"
	const text = "My first todo"

	todoList := NewService(idgen.NewConstantGenerator(expectedID), todoStore, nil)

	todo, err := todoList.AddItem(context.Background(), NewItem{Title: text})
	require.NoError(t, err)

	expectedTodo := Item{
		ID:    expectedID,
		Title: text,
	}

	assert.Equal(t, expectedTodo, todo)

	actualTodo, err := todoStore.Get(context.Background(), todo.ID)
	require.NoError(t, err)

	assert.Equal(t, expectedTodo, actualTodo)
}

func TestList_CannotCreateATodo(t *testing.T) {
	todoList := NewService(idgen.NewConstantGenerator("id"), NewReadOnlyStore(NewInMemoryStore()), nil)

	_, err := todoList.AddItem(context.Background(), NewItem{Title: "My first todo"})
	require.Error(t, err)
}

func TestList_ListTodos(t *testing.T) {
	todoStore := NewInMemoryStore()

	todo := Item{
		ID:    "id",
		Title: "Make the listing work",
	}
	require.NoError(t, todoStore.Store(context.Background(), todo))

	todoList := NewService(idgen.NewConstantGenerator("id"), todoStore, nil)

	todos, err := todoList.ListItems(context.Background())
	require.NoError(t, err)

	expectedTodos := []Item{todo}

	assert.Equal(t, expectedTodos, todos)
}

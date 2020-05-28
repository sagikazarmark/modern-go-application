package todo

import (
	"context"
	"testing"

	"emperror.dev/errors"
	"github.com/goph/idgen"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type todoEventsStub struct {
	markedAsComplete MarkedAsComplete
}

func (s *todoEventsStub) MarkedAsComplete(ctx context.Context, event MarkedAsComplete) error {
	s.markedAsComplete = event

	return nil
}

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

func TestList_MarkAsComplete(t *testing.T) {
	todoStore := NewInMemoryStore()

	const id = "id"

	todo := Item{
		ID:    id,
		Title: "Do me",
	}
	require.NoError(t, todoStore.Store(context.Background(), todo))

	events := &todoEventsStub{}
	todoList := NewService(nil, todoStore, events)

	err := todoList.MarkAsComplete(context.Background(), id)
	require.NoError(t, err)

	expectedTodo := todo
	expectedTodo.Completed = true

	actualTodo, err := todoStore.Get(context.Background(), todo.ID)
	require.NoError(t, err)

	assert.Equal(t, expectedTodo, actualTodo)

	expectedEvent := MarkedAsComplete{
		ID: "id",
	}

	assert.Equal(t, expectedEvent, events.markedAsComplete)
}

func TestList_CannotMarkANonExistingTodoComplete(t *testing.T) {
	todoStore := NewInMemoryStore()

	events := &todoEventsStub{}
	todoList := NewService(nil, todoStore, events)

	const id = "id"

	err := todoList.MarkAsComplete(context.Background(), id)
	require.Error(t, err)

	cause := errors.Cause(err)

	require.IsType(t, NotFoundError{}, cause)

	e := cause.(NotFoundError)
	assert.Equal(t, id, e.ID)
}

func TestList_StoringCompleteTodoFails(t *testing.T) {
	inmemTodoStore := NewInMemoryStore()

	todo := Item{
		ID:    "id",
		Title: "Do me",
	}
	require.NoError(t, inmemTodoStore.Store(context.Background(), todo))

	todoList := NewService(nil, NewReadOnlyStore(inmemTodoStore), &todoEventsStub{})

	err := todoList.MarkAsComplete(context.Background(), "id")
	require.Error(t, err)
}

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
	markedAsDone MarkedAsDone
}

func (s *todoEventsStub) MarkedAsDone(ctx context.Context, event MarkedAsDone) error {
	s.markedAsDone = event

	return nil
}

func TestList_CreatesATodo(t *testing.T) {
	todoStore := NewInMemoryStore()

	const expectedID = "id"
	const text = "My first todo"

	todoList := NewService(idgen.NewConstantGenerator(expectedID), todoStore, nil)

	id, err := todoList.CreateTodo(context.Background(), text)
	require.NoError(t, err)

	assert.Equal(t, expectedID, id)

	expectedTodo := Todo{
		ID:   expectedID,
		Text: text,
	}

	todo, err := todoStore.Get(context.Background(), id)
	require.NoError(t, err)

	assert.Equal(t, expectedTodo, todo)
}

func TestList_CannotCreateATodo(t *testing.T) {
	todoList := NewService(idgen.NewConstantGenerator("id"), NewReadOnlyStore(NewInMemoryStore()), nil)

	_, err := todoList.CreateTodo(context.Background(), "My first todo")
	require.Error(t, err)
}

func TestList_ListTodos(t *testing.T) {
	todoStore := NewInMemoryStore()

	todo := Todo{
		ID:   "id",
		Text: "Make the listing work",
	}
	require.NoError(t, todoStore.Store(context.Background(), todo))

	todoList := NewService(idgen.NewConstantGenerator("id"), todoStore, nil)

	todos, err := todoList.ListTodos(context.Background())
	require.NoError(t, err)

	expectedTodos := []Todo{todo}

	assert.Equal(t, expectedTodos, todos)
}

func TestList_MarkAsDone(t *testing.T) {
	todoStore := NewInMemoryStore()

	const id = "id"

	todo := Todo{
		ID:   id,
		Text: "Do me",
	}
	require.NoError(t, todoStore.Store(context.Background(), todo))

	events := &todoEventsStub{}
	todoList := NewService(nil, todoStore, events)

	err := todoList.MarkAsDone(context.Background(), id)
	require.NoError(t, err)

	expectedTodo := todo
	expectedTodo.Done = true

	actualTodo, err := todoStore.Get(context.Background(), todo.ID)
	require.NoError(t, err)

	assert.Equal(t, expectedTodo, actualTodo)

	expectedEvent := MarkedAsDone{
		ID: "id",
	}

	assert.Equal(t, expectedEvent, events.markedAsDone)
}

func TestList_CannotMarkANonExistingTodoDone(t *testing.T) {
	todoStore := NewInMemoryStore()

	events := &todoEventsStub{}
	todoList := NewService(nil, todoStore, events)

	const id = "id"

	err := todoList.MarkAsDone(context.Background(), id)
	require.Error(t, err)

	cause := errors.Cause(err)

	require.IsType(t, NotFoundError{}, cause)

	e := cause.(NotFoundError)
	assert.Equal(t, id, e.ID)
}

func TestList_StoringDoneTodoFails(t *testing.T) {
	inmemTodoStore := NewInMemoryStore()

	todo := Todo{
		ID:   "id",
		Text: "Do me",
	}
	require.NoError(t, inmemTodoStore.Store(context.Background(), todo))

	todoList := NewService(nil, NewReadOnlyStore(inmemTodoStore), &todoEventsStub{})

	err := todoList.MarkAsDone(context.Background(), "id")
	require.Error(t, err)
}

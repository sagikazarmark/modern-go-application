package todo

import (
	"context"
	"github.com/goph/idgen"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type todoStore struct {
	todos map[string]Todo
}

func (s *todoStore) Store(todo Todo) error {
	s.todos[todo.ID] = todo

	return nil
}

func (s *todoStore) All(ctx context.Context) ([]Todo, error) {
	todos := make([]Todo, len(s.todos))

	i := 0
	for _, todo := range s.todos {
		todos[i] = todo

		i++
	}

	return todos, nil
}

func TestTodoList_CreateTodo(t *testing.T) {
	todoStore := &todoStore{
		todos: make(map[string]Todo),
	}

	todoList := NewTodoList(idgen.NewConstantGenerator("id"), todoStore)

	req := CreateTodoRequest{
		Text: "My first todo",
	}

	resp, err := todoList.CreateTodo(context.Background(), req)
	require.NoError(t, err)

	expectedResponse := &CreateTodoResponse{
		ID: "id",
	}

	assert.Equal(t, expectedResponse, resp)

	todo := Todo{
		ID:   resp.ID,
		Text: req.Text,
	}

	assert.Equal(t, todo, todoStore.todos[todo.ID])
}

func TestTodoList_ListTodos(t *testing.T) {
	todoStore := &todoStore{
		todos: map[string]Todo{
			"id": {
				ID:   "id",
				Text: "Make the listing work",
			},
		},
	}

	todoList := NewTodoList(idgen.NewConstantGenerator("id"), todoStore)

	todos, err := todoList.ListTodos(context.Background())
	require.NoError(t, err)

	expectedTodos := &ListTodosResponse{
		Todos: []Todo{
			{
				ID:   "id",
				Text: "Make the listing work",
			},
		},
	}

	assert.Equal(t, expectedTodos, todos)
}

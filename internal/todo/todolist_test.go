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

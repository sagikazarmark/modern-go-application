package todo

import (
	"context"
	"github.com/goph/idgen"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTodoList_CreateTodo(t *testing.T) {
	todoStore := NewInmemoryTodoStore()
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
	todoStore := NewInmemoryTodoStore()
	todoStore.todos["id"] = Todo{
		ID:   "id",
		Text: "Make the listing work",
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

func TestTodoList_MarkAsDone(t *testing.T) {
	todo := Todo{
		ID:   "id",
		Text: "Make the listing work",
	}

	todoStore := NewInmemoryTodoStore()
	todoStore.todos["id"] = todo

	todoList := NewTodoList(nil, todoStore)

	req := MarkAsDoneRequest{
		ID: "id",
	}
	err := todoList.MarkAsDone(context.Background(), req)
	require.NoError(t, err)

	expectedTodo := todo
	expectedTodo.Done = true

	assert.Equal(t, expectedTodo, todoStore.todos["id"])
}

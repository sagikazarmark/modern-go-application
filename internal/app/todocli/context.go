package todocli

import (
	todov1 "github.com/sagikazarmark/todobackend-go-kit/api/todo/v1"
)

type context struct {
	client todov1.TodoListServiceClient
}

func (c *context) GetTodoClient() todov1.TodoListServiceClient {
	return c.client
}

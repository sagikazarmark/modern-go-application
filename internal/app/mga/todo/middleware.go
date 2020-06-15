package todo

import (
	"context"

	"github.com/sagikazarmark/todobackend-go-kit/todo"
)

// Middleware is a service middleware.
type Middleware func(todo.Service) todo.Service

// DefaultMiddleware helps implementing partial middleware.
type DefaultMiddleware struct {
	Service todo.Service
}

func (m DefaultMiddleware) AddItem(ctx context.Context, newItem todo.NewItem) (todo.Item, error) {
	return m.Service.AddItem(ctx, newItem)
}

func (m DefaultMiddleware) ListItems(ctx context.Context) ([]todo.Item, error) {
	return m.Service.ListItems(ctx)
}

func (m DefaultMiddleware) DeleteItems(ctx context.Context) error {
	return m.Service.DeleteItems(ctx)
}

func (m DefaultMiddleware) GetItem(ctx context.Context, id string) (todo.Item, error) {
	return m.Service.GetItem(ctx, id)
}

func (m DefaultMiddleware) UpdateItem(ctx context.Context, id string, itemUpdate todo.ItemUpdate) (todo.Item, error) {
	return m.Service.UpdateItem(ctx, id, itemUpdate)
}

func (m DefaultMiddleware) DeleteItem(ctx context.Context, id string) error {
	return m.Service.DeleteItem(ctx, id)
}

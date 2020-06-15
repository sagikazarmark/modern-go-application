package todoadapter

import (
	"context"

	"emperror.dev/errors"

	"github.com/sagikazarmark/todobackend-go-kit/todo"

	"github.com/sagikazarmark/modern-go-application/internal/app/mga/todo/todoadapter/ent"
	"github.com/sagikazarmark/modern-go-application/internal/app/mga/todo/todoadapter/ent/todoitem"
)

type entStore struct {
	client *ent.Client
}

// NewEntStore returns a new todo store backed by Ent ORM.
func NewEntStore(client *ent.Client) todo.Store {
	return entStore{
		client: client,
	}
}

func (s entStore) Store(ctx context.Context, todo todo.Item) error {
	existing, err := s.client.TodoItem.Query().Where(todoitem.UID(todo.ID)).First(ctx)
	if ent.IsNotFound(err) {
		_, err := s.client.TodoItem.Create().
			SetUID(todo.ID).
			SetTitle(todo.Title).
			SetCompleted(todo.Completed).
			SetOrder(todo.Order).
			Save(ctx)
		if err != nil {
			return err
		}

		return nil
	}
	if err != nil {
		return err
	}

	_, err = s.client.TodoItem.UpdateOneID(existing.ID).
		SetTitle(todo.Title).
		SetCompleted(todo.Completed).
		SetOrder(todo.Order).
		Save(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s entStore) GetAll(ctx context.Context) ([]todo.Item, error) {
	todoModels, err := s.client.TodoItem.Query().All(ctx)
	if err != nil {
		return nil, err
	}

	todos := make([]todo.Item, 0, len(todoModels))

	for _, todoModel := range todoModels {
		todos = append(todos, todo.Item{
			ID:        todoModel.UID,
			Title:     todoModel.Title,
			Completed: todoModel.Completed,
			Order:     todoModel.Order,
		})
	}

	return todos, nil
}

func (s entStore) GetOne(ctx context.Context, id string) (todo.Item, error) {
	todoModel, err := s.client.TodoItem.Query().Where(todoitem.UID(id)).First(ctx)
	if ent.IsNotFound(err) {
		return todo.Item{}, errors.WithStack(todo.NotFoundError{ID: id})
	}

	return todo.Item{
		ID:        todoModel.UID,
		Title:     todoModel.Title,
		Completed: todoModel.Completed,
		Order:     todoModel.Order,
	}, nil
}

func (s entStore) DeleteAll(ctx context.Context) error {
	_, err := s.client.TodoItem.Delete().Exec(ctx)

	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (s entStore) DeleteOne(ctx context.Context, id string) error {
	_, err := s.client.TodoItem.Delete().Where(todoitem.UID(id)).Exec(ctx)

	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

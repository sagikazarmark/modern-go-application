package todoadapter

import (
	"context"

	"emperror.dev/errors"

	"github.com/sagikazarmark/modern-go-application/internal/app/mga/todo"
	"github.com/sagikazarmark/modern-go-application/internal/app/mga/todo/todoadapter/ent"
	todop "github.com/sagikazarmark/modern-go-application/internal/app/mga/todo/todoadapter/ent/todo"
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

func (s entStore) Store(ctx context.Context, todo todo.Todo) error {
	existing, err := s.client.Todo.Query().Where(todop.UID(todo.ID)).First(ctx)
	if ent.IsNotFound(err) {
		_, err := s.client.Todo.Create().
			SetUID(todo.ID).
			SetText(todo.Text).
			SetDone(todo.Done).
			Save(ctx)
		if err != nil {
			return err
		}

		return nil
	}
	if err != nil {
		return err
	}

	_, err = s.client.Todo.UpdateOneID(existing.ID).
		SetText(todo.Text).
		SetDone(todo.Done).
		Save(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s entStore) All(ctx context.Context) ([]todo.Todo, error) {
	todoModels, err := s.client.Todo.Query().All(ctx)
	if err != nil {
		return nil, err
	}

	todos := make([]todo.Todo, 0, len(todoModels))

	for _, todoModel := range todoModels {
		todos = append(todos, todo.Todo{
			ID:   todoModel.UID,
			Text: todoModel.Text,
			Done: todoModel.Done,
		})
	}

	return todos, nil
}

func (s entStore) Get(ctx context.Context, id string) (todo.Todo, error) {
	todoModel, err := s.client.Todo.Query().Where(todop.UID(id)).First(ctx)
	if ent.IsNotFound(err) {
		return todo.Todo{}, errors.WithStack(todo.NotFoundError{ID: id})
	}

	return todo.Todo{
		ID:   todoModel.UID,
		Text: todoModel.Text,
		Done: todoModel.Done,
	}, nil
}

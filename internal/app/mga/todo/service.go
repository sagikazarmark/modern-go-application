package todo

import (
	"context"

	"emperror.dev/errors"
	"github.com/sagikazarmark/todobackend-go-kit/todo"
)

// +mga:event:dispatcher

// Events dispatches todo events.
type Events interface {
	// MarkedAsComplete dispatches a MarkedAsComplete event.
	MarkedAsComplete(ctx context.Context, event MarkedAsComplete) error
}

// +mga:event:handler

// MarkedAsComplete event is triggered when an item gets marked as complete.
type MarkedAsComplete struct {
	ID string
}

// EventMiddleware fires todo events.
func EventMiddleware(events Events) Middleware {
	return func(next todo.Service) todo.Service {
		return eventMiddleware{
			Service: DefaultMiddleware{Service: next},
			next:    next,

			events: events,
		}
	}
}

type eventMiddleware struct {
	todo.Service
	next todo.Service

	events Events
}

func (mw eventMiddleware) UpdateItem(ctx context.Context, id string, itemUpdate todo.ItemUpdate) (todo.Item, error) {
	var fireComplete bool
	if itemUpdate.Completed != nil && *itemUpdate.Completed {
		fireComplete = true
	}

	item, err := mw.next.UpdateItem(ctx, id, itemUpdate)
	if err != nil {
		return item, err
	}

	if fireComplete {
		event := MarkedAsComplete{
			ID: item.ID,
		}

		err = mw.events.MarkedAsComplete(ctx, event)
		if err != nil {
			// TODO: rollback item store here? retry?
			return item, errors.WithMessage(err, "mark item as complete")
		}
	}

	return item, nil
}

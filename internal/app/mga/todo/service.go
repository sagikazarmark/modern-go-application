package todo

import (
	"context"

	"emperror.dev/errors"
)

// +kit:endpoint:errorStrategy=service

// Service manages a todo list.
type Service interface {
	// AddItem adds a new item to the list.
	AddItem(ctx context.Context, newItem NewItem) (item Item, err error)

	// ListItems returns a list of items.
	ListItems(ctx context.Context) (items []Item, err error)

	// DeleteItems deletes all items from the list.
	DeleteItems(ctx context.Context) error

	// GetItem returns the details of an item.
	GetItem(ctx context.Context, id string) (item Item, err error)

	// UpdateItem updates an existing item.
	UpdateItem(ctx context.Context, id string, itemUpdate ItemUpdate) (item Item, err error)

	// DeleteItem deletes an item from the list.
	DeleteItem(ctx context.Context, id string) error
}

// Item is a note describing a task to be done.
type Item struct {
	ID        string
	Title     string
	Completed bool
	Order     int
}

// NewItem contains the details of a new Item.
type NewItem struct {
	Title string
	Order int
}

func (i NewItem) toItem(id string) Item {
	return Item{
		ID:    id,
		Title: i.Title,
		Order: i.Order,
	}
}

// ItemUpdate contains updates of an existing item.
type ItemUpdate struct {
	Title     *string
	Completed *bool
	Order     *int
}

func (i ItemUpdate) update(item Item) Item {
	if i.Title != nil {
		item.Title = *i.Title
	}

	if i.Completed != nil {
		item.Completed = *i.Completed
	}

	if i.Order != nil {
		item.Order = *i.Order
	}

	return item
}

// NewService returns a new Service.
func NewService(idgenerator IDGenerator, store Store, events Events) Service {
	return &service{
		idgenerator: idgenerator,
		store:       store,
		events:      events,
	}
}

type service struct {
	idgenerator IDGenerator
	store       Store
	events      Events
}

// IDGenerator generates a new ID.
type IDGenerator interface {
	// Generate generates a new ID.
	Generate() (string, error)
}

// Store persists items.
type Store interface {
	// Store stores an item.
	Store(ctx context.Context, item Item) error

	// All returns all items.
	All(ctx context.Context) ([]Item, error)

	// DeleteItems deletes all items in the store.
	DeleteAll(ctx context.Context) error

	// Get returns a single item by its ID.
	Get(ctx context.Context, id string) (Item, error)

	// DeleteOne deletes a single item by its ID.
	DeleteOne(ctx context.Context, id string) error
}

// NotFoundError is returned if an item cannot be found.
type NotFoundError struct {
	ID string
}

// Error implements the error interface.
func (NotFoundError) Error() string {
	return "item not found"
}

// Details returns error details.
func (e NotFoundError) Details() []interface{} {
	return []interface{}{"item_id", e.ID}
}

// NotFound tells a client that this error is related to a resource being not found.
// Can be used to translate the error to eg. status code.
func (NotFoundError) NotFound() bool {
	return true
}

// ServiceError tells the transport layer whether this error should be translated into the transport format
// or an internal error should be returned instead.
func (NotFoundError) ServiceError() bool {
	return true
}

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

type validationError struct {
	violations map[string][]string
}

func (validationError) Error() string {
	return "invalid item"
}

func (e validationError) Violations() map[string][]string {
	return e.violations
}

// Validation tells a client that this error is related to a resource being invalid.
// Can be used to translate the error to eg. status code.
func (validationError) Validation() bool {
	return true
}

// ServiceError tells the transport layer whether this error should be translated into the transport format
// or an internal error should be returned instead.
func (validationError) ServiceError() bool {
	return true
}

func (s service) AddItem(ctx context.Context, newItem NewItem) (Item, error) {
	id, err := s.idgenerator.Generate()
	if err != nil {
		return Item{}, err
	}

	if newItem.Title == "" {
		return Item{}, errors.WithStack(validationError{violations: map[string][]string{
			"title": {
				"title cannot be empty",
			},
		}})
	}

	item := newItem.toItem(id)

	err = s.store.Store(ctx, item)
	if err != nil {
		return Item{}, err
	}

	return item, nil
}

func (s service) ListItems(ctx context.Context) ([]Item, error) {
	return s.store.All(ctx)
}

func (s service) DeleteItems(ctx context.Context) error {
	return s.store.DeleteAll(ctx)
}

func (s service) GetItem(ctx context.Context, id string) (Item, error) {
	item, err := s.store.Get(ctx, id)
	if err != nil {
		return Item{}, errors.WithMessage(err, "get item")
	}

	return item, nil
}

func (s service) UpdateItem(ctx context.Context, id string, itemUpdate ItemUpdate) (Item, error) {
	item, err := s.store.Get(ctx, id)
	if err != nil {
		return Item{}, err
	}

	updatedItem := itemUpdate.update(item)

	if item == updatedItem {
		return item, nil
	}

	err = s.store.Store(ctx, updatedItem)
	if err != nil {
		return Item{}, errors.WithMessage(err, "update item")
	}

	if !item.Completed && updatedItem.Completed {
		event := MarkedAsComplete{
			ID: item.ID,
		}

		err = s.events.MarkedAsComplete(ctx, event)
		if err != nil {
			// TODO: rollback item store here? retry?
			return Item{}, errors.WithMessage(err, "mark item as complete")
		}
	}

	return updatedItem, nil
}

func (s service) DeleteItem(ctx context.Context, id string) error {
	err := s.store.DeleteOne(ctx, id)
	if err != nil {
		return errors.WithMessage(err, "delete item")
	}

	return nil
}

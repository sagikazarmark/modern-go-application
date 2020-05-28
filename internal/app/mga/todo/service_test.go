package todo

import (
	"context"
	"testing"

	"emperror.dev/errors"
	"github.com/go-bdd/gobdd"
	"github.com/goph/idgen"
	"github.com/goph/idgen/ulidgen"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type todoEventsStub struct {
	markedAsComplete MarkedAsComplete
}

func (s *todoEventsStub) MarkedAsComplete(ctx context.Context, event MarkedAsComplete) error {
	s.markedAsComplete = event

	return nil
}

func TestList_CreatesATodo(t *testing.T) {
	todoStore := NewInMemoryStore()

	const expectedID = "id"
	const text = "My first todo"

	todoList := NewService(idgen.NewConstantGenerator(expectedID), todoStore, nil)

	todo, err := todoList.AddItem(context.Background(), NewItem{Title: text})
	require.NoError(t, err)

	expectedTodo := Item{
		ID:    expectedID,
		Title: text,
	}

	assert.Equal(t, expectedTodo, todo)

	actualTodo, err := todoStore.Get(context.Background(), todo.ID)
	require.NoError(t, err)

	assert.Equal(t, expectedTodo, actualTodo)
}

func TestList_CannotCreateATodo(t *testing.T) {
	todoList := NewService(idgen.NewConstantGenerator("id"), NewReadOnlyStore(NewInMemoryStore()), nil)

	_, err := todoList.AddItem(context.Background(), NewItem{Title: "My first todo"})
	require.Error(t, err)
}

func TestList_ListTodos(t *testing.T) {
	todoStore := NewInMemoryStore()

	todo := Item{
		ID:    "id",
		Title: "Make the listing work",
	}
	require.NoError(t, todoStore.Store(context.Background(), todo))

	todoList := NewService(idgen.NewConstantGenerator("id"), todoStore, nil)

	todos, err := todoList.ListItems(context.Background())
	require.NoError(t, err)

	expectedTodos := []Item{todo}

	assert.Equal(t, expectedTodos, todos)
}

func TestList_MarkAsComplete(t *testing.T) {
	todoStore := NewInMemoryStore()

	const id = "id"

	todo := Item{
		ID:    id,
		Title: "Do me",
	}
	require.NoError(t, todoStore.Store(context.Background(), todo))

	events := &todoEventsStub{}
	todoList := NewService(nil, todoStore, events)

	err := todoList.MarkAsComplete(context.Background(), id)
	require.NoError(t, err)

	expectedTodo := todo
	expectedTodo.Completed = true

	actualTodo, err := todoStore.Get(context.Background(), todo.ID)
	require.NoError(t, err)

	assert.Equal(t, expectedTodo, actualTodo)

	expectedEvent := MarkedAsComplete{
		ID: "id",
	}

	assert.Equal(t, expectedEvent, events.markedAsComplete)
}

func TestList_CannotMarkANonExistingTodoComplete(t *testing.T) {
	todoStore := NewInMemoryStore()

	events := &todoEventsStub{}
	todoList := NewService(nil, todoStore, events)

	const id = "id"

	err := todoList.MarkAsComplete(context.Background(), id)
	require.Error(t, err)

	cause := errors.Cause(err)

	require.IsType(t, NotFoundError{}, cause)

	e := cause.(NotFoundError)
	assert.Equal(t, id, e.ID)
}

func TestList_StoringCompleteTodoFails(t *testing.T) {
	inmemTodoStore := NewInMemoryStore()

	todo := Item{
		ID:    "id",
		Title: "Do me",
	}
	require.NoError(t, inmemTodoStore.Store(context.Background(), todo))

	todoList := NewService(nil, NewReadOnlyStore(inmemTodoStore), &todoEventsStub{})

	err := todoList.MarkAsComplete(context.Background(), "id")
	require.Error(t, err)
}

type FeatureContext struct {
	Store   Store
	Service Service
}

func getFeatureContext(t gobdd.StepTest, ctx gobdd.Context) FeatureContext {
	v, err := ctx.Get("ctx")
	if err != nil {
		t.Fatal(err)
	}

	return v.(FeatureContext)
}

// nolint: gocognit
func TestList(t *testing.T) {
	suite := gobdd.NewSuite(t, gobdd.WithBeforeScenario(func(ctx gobdd.Context) {
		store := NewInMemoryStore()
		service := NewService(ulidgen.NewGenerator(), store, &todoEventsStub{})

		ctx.Set("ctx", FeatureContext{
			Store:   store,
			Service: service,
		})
	}))

	suite.AddStep(`(?:I|the user) adds? a new todo "(.*)" to the list`,
		func(t gobdd.StepTest, ctx gobdd.Context, text string) {
			fctx := getFeatureContext(t, ctx)

			todo, err := fctx.Service.AddItem(context.Background(), NewItem{Title: text})
			if err != nil {
				var cerr interface{ ServiceError() bool }

				if !errors.As(err, &cerr) || !cerr.ServiceError() {
					t.Fatal(err)
				}

				ctx.Set("error", err)

				return
			}

			ctx.Set("id", todo.ID)
		})

	suite.AddStep(`"(.+)" should be on the list`, func(t gobdd.StepTest, ctx gobdd.Context, text string) {
		if err, _ := ctx.GetError("error", nil); err != nil {
			t.Fatal(err)
		}

		fctx := getFeatureContext(t, ctx)

		id, _ := ctx.GetString("id")
		todo, err := fctx.Store.Get(context.Background(), id)
		if err != nil {
			t.Fatal(err)
		}

		if todo.Title != text {
			t.Errorf("cannot find %q todo entry", text)
		}
	})

	suite.AddStep(`it should fail with a validation error for the "(.+)" field saying that "(.+)"`,
		func(t gobdd.StepTest, ctx gobdd.Context, field string, violation string) {
			var err error
			{ // See https://github.com/go-bdd/gobdd/pull/95
				v, _ := ctx.GetError("error", nil)
				if v == nil {
					t.Fatal("a validation error was expected, but received none")
				}

				err = v.(error)
			}

			var verr interface {
				Validation() bool
				Violations() map[string][]string
			}

			if !errors.As(err, &verr) {
				t.Fatalf("a validation error was expected, the received error is not one: %s", err)
			}

			violations := verr.Violations()

			fieldViolations, ok := violations[field]
			if !ok || len(fieldViolations) == 0 {
				t.Fatalf("the returned validation error does not have violations for %q field", field)
			}

			if fieldViolations[0] != violation {
				t.Errorf("the %q field does not have a(n) %q violation", field, violation)
			}
		})

	suite.AddStep(`there is a todo "(.*)"`, func(t gobdd.StepTest, ctx gobdd.Context, text string) {
		fctx := getFeatureContext(t, ctx)

		const id = "todo"

		err := fctx.Store.Store(context.Background(), Item{
			ID:        id,
			Title:     text,
			Completed: false,
		})
		if err != nil {
			t.Fatal(err)
		}

		ctx.Set("id", id)
	})

	suite.AddStep(`(?:I|the user) marks? it as complete`, func(t gobdd.StepTest, ctx gobdd.Context) {
		fctx := getFeatureContext(t, ctx)

		id, _ := ctx.GetString("id")

		err := fctx.Service.MarkAsComplete(context.Background(), id)
		if err != nil {
			var cerr interface{ ServiceError() bool }

			if !errors.As(err, &cerr) || !cerr.ServiceError() {
				t.Fatal(err)
			}

			ctx.Set("error", err)

			return
		}
	})

	suite.AddStep(`it should be complete`, func(t gobdd.StepTest, ctx gobdd.Context) {
		if err, _ := ctx.GetError("error", nil); err != nil {
			t.Fatal(err)
		}

		fctx := getFeatureContext(t, ctx)

		id, _ := ctx.GetString("id")

		todo, err := fctx.Store.Get(context.Background(), id)
		if err != nil {
			t.Fatal(err)
		}

		if !todo.Completed {
			t.Error("todo is expected to be complete")
		}
	})

	suite.Run()
}

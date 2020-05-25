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
	markedAsDone MarkedAsDone
}

func (s *todoEventsStub) MarkedAsDone(ctx context.Context, event MarkedAsDone) error {
	s.markedAsDone = event

	return nil
}

func TestList_CreatesATodo(t *testing.T) {
	todoStore := NewInMemoryStore()

	const expectedID = "id"
	const text = "My first todo"

	todoList := NewService(idgen.NewConstantGenerator(expectedID), todoStore, nil)

	id, err := todoList.CreateTodo(context.Background(), text)
	require.NoError(t, err)

	assert.Equal(t, expectedID, id)

	expectedTodo := Todo{
		ID:   expectedID,
		Text: text,
	}

	todo, err := todoStore.Get(context.Background(), id)
	require.NoError(t, err)

	assert.Equal(t, expectedTodo, todo)
}

func TestList_CannotCreateATodo(t *testing.T) {
	todoList := NewService(idgen.NewConstantGenerator("id"), NewReadOnlyStore(NewInMemoryStore()), nil)

	_, err := todoList.CreateTodo(context.Background(), "My first todo")
	require.Error(t, err)
}

func TestList_ListTodos(t *testing.T) {
	todoStore := NewInMemoryStore()

	todo := Todo{
		ID:   "id",
		Text: "Make the listing work",
	}
	require.NoError(t, todoStore.Store(context.Background(), todo))

	todoList := NewService(idgen.NewConstantGenerator("id"), todoStore, nil)

	todos, err := todoList.ListTodos(context.Background())
	require.NoError(t, err)

	expectedTodos := []Todo{todo}

	assert.Equal(t, expectedTodos, todos)
}

func TestList_MarkAsDone(t *testing.T) {
	todoStore := NewInMemoryStore()

	const id = "id"

	todo := Todo{
		ID:   id,
		Text: "Do me",
	}
	require.NoError(t, todoStore.Store(context.Background(), todo))

	events := &todoEventsStub{}
	todoList := NewService(nil, todoStore, events)

	err := todoList.MarkAsDone(context.Background(), id)
	require.NoError(t, err)

	expectedTodo := todo
	expectedTodo.Done = true

	actualTodo, err := todoStore.Get(context.Background(), todo.ID)
	require.NoError(t, err)

	assert.Equal(t, expectedTodo, actualTodo)

	expectedEvent := MarkedAsDone{
		ID: "id",
	}

	assert.Equal(t, expectedEvent, events.markedAsDone)
}

func TestList_CannotMarkANonExistingTodoDone(t *testing.T) {
	todoStore := NewInMemoryStore()

	events := &todoEventsStub{}
	todoList := NewService(nil, todoStore, events)

	const id = "id"

	err := todoList.MarkAsDone(context.Background(), id)
	require.Error(t, err)

	cause := errors.Cause(err)

	require.IsType(t, NotFoundError{}, cause)

	e := cause.(NotFoundError)
	assert.Equal(t, id, e.ID)
}

func TestList_StoringDoneTodoFails(t *testing.T) {
	inmemTodoStore := NewInMemoryStore()

	todo := Todo{
		ID:   "id",
		Text: "Do me",
	}
	require.NoError(t, inmemTodoStore.Store(context.Background(), todo))

	todoList := NewService(nil, NewReadOnlyStore(inmemTodoStore), &todoEventsStub{})

	err := todoList.MarkAsDone(context.Background(), "id")
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

func TestList(t *testing.T) {
	suite := gobdd.NewSuite(t, gobdd.WithBeforeScenario(func(ctx gobdd.Context) {
		store := NewInMemoryStore()
		service := NewService(ulidgen.NewGenerator(), store, nil)

		ctx.Set("ctx", FeatureContext{
			Store:   store,
			Service: service,
		})
	}))

	suite.AddStep(`I add entry "(.*)"`, func(t gobdd.StepTest, ctx gobdd.Context, text string) {
		fctx := getFeatureContext(t, ctx)

		id, err := fctx.Service.CreateTodo(context.Background(), text)
		if err != nil {
			var cerr interface{ ServiceError() bool }

			if !errors.As(err, &cerr) || !cerr.ServiceError() {
				t.Fatal(err)
			}

			ctx.Set("error", err)

			return
		}

		ctx.Set("id", id)
	})

	suite.AddStep(`I should have a todo to "(.+)"`, func(t gobdd.StepTest, ctx gobdd.Context, text string) {
		if err, _ := ctx.GetError("error", nil); err != nil {
			t.Fatal(err)
		}

		fctx := getFeatureContext(t, ctx)

		id, _ := ctx.GetString("id")
		todo, err := fctx.Store.Get(context.Background(), id)
		if err != nil {
			t.Fatal(err)
		}

		if todo.Text != text {
			t.Errorf("cannot find %q todo entry", text)
		}

		if todo.Done {
			t.Errorf("%q should not be done", text)
		}
	})

	suite.AddStep(`I should see a validation error for the "(.+)" field saying that "(.+)"`,
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

	suite.Run()
}

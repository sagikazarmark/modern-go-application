package todo_test

import (
	"context"
	"testing"

	"emperror.dev/errors"
	"github.com/go-bdd/gobdd"
	"github.com/goph/idgen/ulidgen"
	"github.com/stretchr/testify/assert"

	"github.com/sagikazarmark/modern-go-application/internal/app/mga/todo"
)

func TestService(t *testing.T) {
	suite := gobdd.NewSuite(t)

	suite.AddStep(`an empty todo list`, givenAnEmptyTodoList)
	suite.AddStep(`(?:(?:I|the user)(?: also)? adds? )?(?:a new|an) item for "(.*)"`, addAnItem)
	suite.AddStep(`it should be (?:the only item )?on the list`, shouldBeOnTheList)
	suite.AddStep(`both items should be on the list`, allShouldBeOnTheList)
	suite.AddStep(`the list should be empty`, theListShouldBeEmpty)
	suite.AddStep(`it is marked as complete`, itemMarkedAsComplete)
	suite.AddStep(`it should be complete`, itemShouldBeComplete)
	suite.AddStep(`it is deleted`, deleteAnItem)
	suite.AddStep(`all items are deleted`, clearList)

	suite.Run()
}

type todoEventsStub struct {
	markedAsComplete todo.MarkedAsComplete
}

func (s *todoEventsStub) MarkedAsComplete(_ context.Context, event todo.MarkedAsComplete) error {
	s.markedAsComplete = event

	return nil
}

func getService(t gobdd.StepTest, ctx gobdd.Context) todo.Service {
	v, err := ctx.Get("service")
	if err != nil {
		t.Fatal(err)
	}

	return v.(todo.Service)
}

func givenAnEmptyTodoList(_ gobdd.StepTest, ctx gobdd.Context) {
	store := todo.NewInMemoryStore()
	service := todo.NewService(ulidgen.NewGenerator(), store, &todoEventsStub{})

	ctx.Set("service", service)
}

func addAnItem(t gobdd.StepTest, ctx gobdd.Context, title string) {
	service := getService(t, ctx)

	item, err := service.AddItem(context.Background(), todo.NewItem{Title: title})
	if err != nil {
		var cerr interface{ ServiceError() bool }

		if !errors.As(err, &cerr) || !cerr.ServiceError() {
			t.Fatal(err)
		}

		ctx.Set("error", err)

		return
	}

	ctx.Set("id", item.ID)
	ctx.Set("title", title)

	ids, _ := ctx.Get("ids", []string{})
	titles, _ := ctx.Get("titles", []string{})

	ctx.Set("ids", append(ids.([]string), item.ID))
	ctx.Set("titles", append(titles.([]string), title))
}

func shouldBeOnTheList(t gobdd.StepTest, ctx gobdd.Context) {
	if err, _ := ctx.GetError("error", nil); err != nil {
		t.Fatal(err)
	}

	service := getService(t, ctx)

	items, err := service.ListItems(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	title, _ := ctx.GetString("title", "")

	assert.Len(t, items, 1, "there should be one item on the list")
	assert.Equal(t, items[0].Title, title, "the item on the list should match the added item")
}

func allShouldBeOnTheList(t gobdd.StepTest, ctx gobdd.Context) {
	if err, _ := ctx.GetError("error", nil); err != nil {
		t.Fatal(err)
	}

	service := getService(t, ctx)

	items, err := service.ListItems(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	ids, _ := ctx.Get("ids", []string{})
	titles, _ := ctx.Get("titles", []string{})

	idMap := make(map[string]string)

	for i, id := range ids.([]string) {
		idMap[id] = titles.([]string)[i]
	}

	assert.Len(t, items, len(idMap))

	for _, item := range items {
		assert.Equal(t, idMap[item.ID], item.Title, "the item on the list should match the added item")
	}
}

func theListShouldBeEmpty(t gobdd.StepTest, ctx gobdd.Context) {
	service := getService(t, ctx)

	items, err := service.ListItems(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	assert.Len(t, items, 0, "the list should be empty")
}

func itemMarkedAsComplete(t gobdd.StepTest, ctx gobdd.Context) {
	id, _ := ctx.GetString("id")

	service := getService(t, ctx)

	completed := true

	_, err := service.UpdateItem(context.Background(), id, todo.ItemUpdate{Completed: &completed})
	if err != nil {
		t.Fatal(err)
	}
}

func itemShouldBeComplete(t gobdd.StepTest, ctx gobdd.Context) {
	id, _ := ctx.GetString("id")

	service := getService(t, ctx)

	item, err := service.GetItem(context.Background(), id)
	if err != nil {
		t.Fatal(err)
	}

	assert.True(t, item.Completed, "item should be complete")
}

func deleteAnItem(t gobdd.StepTest, ctx gobdd.Context) {
	id, _ := ctx.GetString("id")

	service := getService(t, ctx)

	err := service.DeleteItem(context.Background(), id)
	if err != nil {
		t.Fatal(err)
	}
}

func clearList(t gobdd.StepTest, ctx gobdd.Context) {
	service := getService(t, ctx)

	err := service.DeleteItems(context.Background())
	if err != nil {
		t.Fatal(err)
	}
}

package helloworlddriver

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testModel struct {
	Foo string
	Bar string `json:"bar"`
}

func TestJsonView_Render(t *testing.T) {
	view := &jsonView{}

	model := testModel{
		Foo: "foo",
		Bar: "bar",
	}

	var buf bytes.Buffer

	view.Render(&buf, model)

	assert.Equal(t, `{"Foo":"foo","bar":"bar"}`+"\n", buf.String())
}

package greetingdriver

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	err := view.Render(&buf, model)
	require.NoError(t, err)

	assert.Equal(t, `{"Foo":"foo","bar":"bar"}`+"\n", buf.String())
}

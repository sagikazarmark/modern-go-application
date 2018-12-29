package greetingdriver

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/goph/emperror"
	"github.com/pkg/errors"
	"github.com/sagikazarmark/modern-go-application/internal/greeting"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type helloWorldViewStub struct {
	output []byte
	err    error
}

func (v *helloWorldViewStub) Render(output io.Writer, model interface{}) error {
	if v.err != nil {
		return v.err
	}

	_, err := output.Write(v.output)
	if err != nil {
		return err
	}

	return nil
}

func TestHelloWorldWebOutput_Say(t *testing.T) {
	responseWriter := httptest.NewRecorder()
	body := `{"message":"Hello, World!"}`
	view := &helloWorldViewStub{
		output: []byte(body),
	}
	contentType := "application/json"
	output := newGreetingWebOutput(responseWriter, view, contentType, emperror.NewNoopHandler())

	output.Say(context.Background(), greeting.Hello{Message: "Hello, World!"})

	response := responseWriter.Result()
	defer response.Body.Close()

	responseBody, err := ioutil.ReadAll(response.Body)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, contentType, response.Header.Get("Content-Type"))
	assert.Equal(t, body, string(responseBody))
}

func TestHelloWorldWebOutput_Say_RenderError(t *testing.T) {
	responseWriter := httptest.NewRecorder()
	renderError := errors.New("error")
	view := &helloWorldViewStub{
		err: renderError,
	}
	contentType := "application/json"
	output := newGreetingWebOutput(responseWriter, view, contentType, emperror.NewNoopHandler())

	output.Say(context.Background(), greeting.Hello{Message: "Hello, World!"})

	response := responseWriter.Result()
	defer response.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, response.StatusCode)
}

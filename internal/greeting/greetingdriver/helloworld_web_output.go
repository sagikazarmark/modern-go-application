package greetingdriver

import (
	"bytes"
	"context"
	"io"
	"net/http"

	"github.com/goph/emperror"
	"github.com/sagikazarmark/modern-go-application/.gen/openapi/go"
	"github.com/sagikazarmark/modern-go-application/internal/greeting"
)

type helloWorldView interface {
	Render(output io.Writer, model interface{}) error
}

type helloWorldWebOutput struct {
	responseWriter http.ResponseWriter
	view           helloWorldView
	contentType    string

	errorHandler emperror.Handler
}

func newHelloWorldWebOutput(
	responseWriter http.ResponseWriter,
	view helloWorldView,
	contentType string,
	errorHandler emperror.Handler,
) *helloWorldWebOutput {
	return &helloWorldWebOutput{
		responseWriter: responseWriter,
		view:           view,
		contentType:    contentType,
		errorHandler:   errorHandler,
	}
}

func (o *helloWorldWebOutput) Say(ctx context.Context, hello greeting.Hello) {
	response := api.Hello{
		Message: hello.Message,
	}

	var buf bytes.Buffer

	err := o.view.Render(&buf, response)
	if err != nil {
		err = emperror.Wrap(err, "failed to render response")

		o.errorHandler.Handle(err)

		o.responseWriter.WriteHeader(http.StatusInternalServerError)

		return
	}

	o.responseWriter.Header().Set("Content-Type", o.contentType)
	o.responseWriter.WriteHeader(http.StatusOK)

	_, err = io.Copy(o.responseWriter, &buf)
	if err != nil {
		err = emperror.Wrap(err, "failed to send response")

		o.errorHandler.Handle(err)

		o.responseWriter.WriteHeader(http.StatusInternalServerError)

		return
	}
}

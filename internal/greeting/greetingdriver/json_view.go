package greetingdriver

import (
	"encoding/json"
	"io"

	"github.com/pkg/errors"
)

type jsonView struct{}

func (v *jsonView) Render(output io.Writer, model interface{}) error {
	encoder := json.NewEncoder(output)

	err := encoder.Encode(model)
	if err != nil {
		return errors.Wrap(err, "failed to render view")
	}

	return nil
}

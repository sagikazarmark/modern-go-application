package greetingdriver

import (
	"encoding/json"
	"net/http"

	"github.com/goph/emperror"
	"github.com/moogar0880/problems"
	"github.com/pkg/errors"

	api "github.com/sagikazarmark/modern-go-application/.gen/openapi/greeting/go"
	"github.com/sagikazarmark/modern-go-application/internal/greeting"
)

// HTTPController collects the greeting use cases and exposes them as HTTP handlers.
type HTTPController struct {
	greeter Greeter

	errorHandler emperror.Handler
}

// NewHTTPController returns a new HTTPController instance.
func NewHTTPController(greeter Greeter, errorHandler emperror.Handler) *HTTPController {
	return &HTTPController{
		greeter:      greeter,
		errorHandler: errorHandler,
	}
}

// SayHello says hello to someone.
func (c *HTTPController) SayHello(w http.ResponseWriter, r *http.Request) {
	req, err := decodeHTTPSayHelloRequest(r)
	if err != nil {
		c.errorHandler.Handle(err)

		err := encodeHTTPErrorWithCode(err, http.StatusBadRequest, w)
		if err != nil {
			c.errorHandler.Handle(err)
		}

		return
	}

	resp, err := c.greeter.SayHello(r.Context(), req)
	if err != nil {
		err := encodeHTTPErrorWithCode(err, http.StatusBadRequest, w)
		if err != nil {
			c.errorHandler.Handle(err)
		}

		return
	}

	err = encodeHTTPSayHelloResponse(w, resp)
	if err != nil {
		c.errorHandler.Handle(err)
	}
}

func decodeHTTPSayHelloRequest(r *http.Request) (greeting.HelloRequest, error) {
	var req greeting.HelloRequest

	var apiReq api.HelloRequest
	err := json.NewDecoder(r.Body).Decode(&apiReq)
	if err != nil {
		return req, errors.Wrap(err, "failed to decode request")
	}

	req = greeting.HelloRequest{
		Name: apiReq.Name,
	}

	return req, nil
}

func encodeHTTPSayHelloResponse(w http.ResponseWriter, resp *greeting.HelloResponse) error {
	apiResp := api.HelloResponse{
		Reply: resp.Reply,
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(apiResp)

	return errors.Wrap(err, "failed to send response")
}

func encodeHTTPErrorWithCode(err error, statusCode int, w http.ResponseWriter) error {
	problem := problems.NewDetailedProblem(statusCode, err.Error())

	w.Header().Set("Content-Type", problems.ProblemMediaType)
	err = json.NewEncoder(w).Encode(problem)

	return emperror.WrapWith(err, "failed to respond with error", "error", problem.Detail)
}

package greetingdriver

import (
	"encoding/json"
	"net/http"

	"github.com/goph/emperror"
	"github.com/moogar0880/problems"

	"github.com/sagikazarmark/modern-go-application/.gen/openapi/greeting/go"
	"github.com/sagikazarmark/modern-go-application/internal/greeting"
)

// HTTPController collects the greeting use cases and exposes them as HTTP handlers.
type HTTPController struct {
	helloService HelloService

	errorHandler emperror.Handler
}

// NewHTTPController returns a new HTTPController instance.
func NewHTTPController(helloService HelloService, errorHandler emperror.Handler) *HTTPController {
	return &HTTPController{
		helloService: helloService,
		errorHandler: errorHandler,
	}
}

// SayHello says hello to someone.
func (c *HTTPController) SayHello(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var apiReq api.HelloRequest

	if err := decoder.Decode(&apiReq); err != nil {
		c.errorHandler.Handle(err)

		http.Error(w, "invalid request", http.StatusBadRequest)

		return
	}

	req := greeting.HelloRequest{
		Greeting: apiReq.Greeting,
	}

	resp, err := c.helloService.SayHello(r.Context(), req)
	if err != nil {
		problem := problems.NewDetailedProblem(400, err.Error())

		w.Header().Set("Content-Type", problems.ProblemMediaType)
		err = json.NewEncoder(w).Encode(problem)
		if err != nil {
			c.errorHandler.Handle(emperror.WrapWith(err, "failed to respond with error", "error", problem.Detail))
		}

		return
	}

	apiResp := api.HelloResponse{
		Reply: resp.Reply,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(apiResp)
	if err != nil {
		c.errorHandler.Handle(emperror.WrapWith(err, "failed to send response"))
	}
}

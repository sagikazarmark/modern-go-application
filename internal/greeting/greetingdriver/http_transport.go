package greetingdriver

import (
	"context"
	"encoding/json"
	"net/http"

	kitoc "github.com/go-kit/kit/tracing/opencensus"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/goph/emperror"
	"github.com/gorilla/mux"
	"github.com/moogar0880/problems"
	"github.com/pkg/errors"
	api "github.com/sagikazarmark/modern-go-application/.gen/api/openapi/greeting/go"
	"github.com/sagikazarmark/modern-go-application/internal/greeting"
)

// MakeHTTPHandler mounts all of the service endpoints into an http.Handler.
func MakeHTTPHandler(greeter Greeter, errorHandler greeting.ErrorHandler) http.Handler {
	r := mux.NewRouter()

	errorEncoder := httptransport.ServerErrorEncoder(func(ctx context.Context, err error, w http.ResponseWriter) {
		// This replaces server error log
		errorHandler.Handle(err)

		err = encodeHTTPError(err, w)
		if err != nil {
			errorHandler.Handle(err)
		}
	})

	r.Methods("POST").Path("/sayHello").Handler(httptransport.NewServer(
		kitoc.TraceEndpoint("SayHello")(MakeSayHelloEndpoint(greeter)),
		decodeSayHelloHTTPRequest,
		encodeSayHelloHTTPResponse,
		errorEncoder,
	))

	return r
}

func decodeSayHelloHTTPRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req HelloRequest

	var apiReq api.HelloRequest
	err := json.NewDecoder(r.Body).Decode(&apiReq)
	if err != nil {
		return req, errors.Wrap(err, "failed to decode request")
	}

	req = HelloRequest{
		Name: apiReq.Name,
	}

	return req, nil
}

func encodeSayHelloHTTPResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	resp := response.(HelloResponse)

	apiResp := api.HelloResponse{
		Reply: resp.Reply,
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(apiResp)

	return errors.Wrap(err, "failed to send response")
}

func encodeHTTPError(err error, w http.ResponseWriter) error {
	problem := problems.NewDetailedProblem(http.StatusInternalServerError, err.Error())

	w.Header().Set("Content-Type", problems.ProblemMediaType)
	err = json.NewEncoder(w).Encode(problem)

	return emperror.WrapWith(err, "failed to respond with error", "error", problem.Detail)
}

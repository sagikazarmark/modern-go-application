package appkit

import (
	"context"
	"net/http"

	"emperror.dev/errors"
	"github.com/moogar0880/problems"
	kitxgrpc "github.com/sagikazarmark/kitx/transport/grpc"
	kitxhttp "github.com/sagikazarmark/kitx/transport/http"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewProblemConverter() kitxhttp.ProblemConverter {
	return kitxhttp.NewProblemConverter(kitxhttp.ProblemConverterConfig{
		Matchers: []kitxhttp.ProblemMatcher{
			kitxhttp.NewStatusProblemMatcher(http.StatusNotFound, kitxhttp.ErrorMatcherFunc(IsNotFoundError)),
			validationProblemConverter{},
		},
	})
}

type validationProblemConverter struct{}

func (v validationProblemConverter) NewProblem(_ context.Context, err error) problems.Problem {
	var verr interface {
		Violations() map[string][]string
	}

	if errors.As(err, &verr) {
		return NewValidationProblem(err.Error(), verr.Violations())
	}

	return problems.NewDetailedProblem(http.StatusUnprocessableEntity, err.Error())
}

func (v validationProblemConverter) MatchError(err error) bool {
	return IsValidationError(err)
}

func NewStatusConverter() kitxgrpc.StatusConverter {
	return kitxgrpc.NewStatusConverter(kitxgrpc.StatusConverterConfig{
		Matchers: []kitxgrpc.StatusMatcher{
			kitxgrpc.NewStatusCodeMatcher(codes.NotFound, kitxgrpc.ErrorMatcherFunc(IsNotFoundError)),
			validationStatusConverter{},
		},
	})
}

type validationStatusConverter struct{}

func (v validationStatusConverter) NewStatus(_ context.Context, err error) *status.Status {
	var verr interface {
		Violations() map[string][]string
	}

	if errors.As(err, &verr) {
		st := status.New(codes.InvalidArgument, err.Error())

		br := &errdetails.BadRequest{}

		for field, violations := range verr.Violations() {
			for _, violation := range violations {
				br.FieldViolations = append(br.FieldViolations, &errdetails.BadRequest_FieldViolation{
					Field:       field,
					Description: violation,
				})
			}
		}

		st, err := st.WithDetails(br)
		if err != nil {
			// If this errored, it will always error
			// here, so better panic so we can figure
			// out why than have this silently passing.
			panic(errors.Wrap(err, "unexpected error attaching metadata"))
		}

		return st
	}

	return status.New(codes.InvalidArgument, err.Error())
}

func (v validationStatusConverter) MatchError(err error) bool {
	return IsValidationError(err)
}

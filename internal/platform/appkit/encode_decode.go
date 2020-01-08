package appkit

import (
	"net/http"

	kitxgrpc "github.com/sagikazarmark/kitx/transport/grpc"
	kitxhttp "github.com/sagikazarmark/kitx/transport/http"
	"google.golang.org/grpc/codes"
)

func NewProblemConverter() kitxhttp.ProblemConverter {
	return kitxhttp.NewProblemConverter(kitxhttp.ProblemConverterConfig{
		Matchers: []kitxhttp.ProblemMatcher{
			kitxhttp.NewStatusProblemMatcher(http.StatusNotFound, kitxhttp.ErrorMatcherFunc(IsNotFoundError)),
		},
	})
}

func NewStatusConverter() kitxgrpc.StatusConverter {
	return kitxgrpc.NewStatusConverter(kitxgrpc.StatusConverterConfig{
		Matchers: []kitxgrpc.StatusMatcher{
			kitxgrpc.NewStatusCodeMatcher(codes.NotFound, kitxgrpc.ErrorMatcherFunc(IsNotFoundError)),
		},
	})
}

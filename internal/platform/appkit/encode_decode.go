package appkit

import (
	"net/http"

	kitxhttp "github.com/sagikazarmark/kitx/transport/http"
)

func NewProblemConverter() kitxhttp.ProblemConverter {
	return kitxhttp.NewProblemConverter(kitxhttp.ProblemConverterConfig{
		Matchers: []kitxhttp.ProblemMatcher{
			kitxhttp.NewStatusProblemMatcher(http.StatusNotFound, kitxhttp.ErrorMatcherFunc(IsNotFoundError)),
		},
	})
}

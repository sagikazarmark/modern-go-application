package appkit

import (
	"net/http"

	"github.com/moogar0880/problems"
)

// ValidationProblem describes an RFC-7807 problem.
type ValidationProblem struct {
	*problems.DefaultProblem

	Violations map[string][]string `json:"violations"`
}

// NewValidationProblem returns a problem with details and validation errors.
func NewValidationProblem(details string, violations map[string][]string) *ValidationProblem {
	problem := problems.NewDetailedProblem(http.StatusUnprocessableEntity, details)

	return &ValidationProblem{
		DefaultProblem: problem,
		Violations:     violations,
	}
}

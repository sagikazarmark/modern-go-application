package httpbin

import (
	"github.com/sagikazarmark/modern-go-application/internal/common"
)

// These interfaces are aliased so that the module code is separated from the rest of the application.
// If the module is moved out of the app, copy the aliased interfaces here.

// Logger is the fundamental interface for all log operations.
type Logger = common.Logger

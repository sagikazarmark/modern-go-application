package greetingdriver

import (
	"flag"
	"regexp"
	"testing"
)

func TestFunctional(t *testing.T) {
	if m := flag.Lookup("test.run").Value.String(); m == "" || !regexp.MustCompile(m).MatchString(t.Name()) {
		t.Skip("skipping as execution was not requested explicitly using go test -run")
	}

	t.Parallel()

	t.Run("HTTPController_SayHello", testSayHello)
}

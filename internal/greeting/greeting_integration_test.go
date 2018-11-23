package greeting_test

import (
	"flag"
	"regexp"
	"testing"
)

func TestIntegration(t *testing.T) {
	if m := flag.Lookup("test.run").Value.String(); m == "" || !regexp.MustCompile(m).MatchString(t.Name()) {
		t.Skip("skipping as execution was not requested explicitly using go test -run")
	}

	t.Parallel()

	t.Run("HelloWorld", testHelloWorld)
	t.Run("SayHello", testSayHello)
}

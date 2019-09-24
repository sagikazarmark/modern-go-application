package correlation

import (
	"context"
	"testing"
)

func TestMiddleware(t *testing.T) {
	t.Run("existing", func(t *testing.T) {
		cid := "1234"

		e := func(ctx context.Context, _ interface{}) (interface{}, error) {
			c := ctx.Value(correlationIDContextKey)
			if cid != c {
				t.Errorf(
					"context does not contain the expected encoded correlation ID\nexpected: %s\nactual:   %s",
					cid,
					c,
				)
			}

			return ctx, nil
		}
		e = Middleware()(e)

		ctx := context.WithValue(context.Background(), correlationIDContextKey, cid)

		e(ctx, nil)
	})

	t.Run("generated", func(t *testing.T) {
		cid := "BpLnfgDsc2WD8F2qNfHK5a84jjJkwzDk" // Result of seed(1)

		e := func(ctx context.Context, _ interface{}) (interface{}, error) {
			c := ctx.Value(correlationIDContextKey)
			if cid != c {
				t.Errorf(
					"context does not contain the expected encoded correlation ID\nexpected: %s\nactual:   %s",
					cid,
					c,
				)
			}

			return ctx, nil
		}
		e = Middleware()(e)

		ctx := context.Background()

		e(ctx, nil)
	})
}

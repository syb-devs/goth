package timer

import (
	"fmt"
	"time"

	"bitbucket.org/syb-devs/goth/app"
)

// New returns a timer middleware function
func New() app.Middleware {
	return func(h app.Handler) app.Handler {
		return app.HandlerFunc(
			func(ctx *app.Context) error {
				start := time.Now()
				err := h.Serve(ctx)
				elapsed := time.Now().Sub(start)
				ctx.Header().Set("Elapsed", elapsed.String())
				fmt.Printf("HTTP Headers:\n %v\n", ctx.Header())
				return err
			})
	}
}

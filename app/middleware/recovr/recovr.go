package recovr

import (
	"net/http"
	"runtime/debug"

	"bitbucket.org/syb-devs/goth/app"
)

// New returns a timer middleware function
func New() app.Middleware {
	return func(h app.Handler) app.Handler {
		return app.HandlerFunc(
			func(ctx *app.Context) error {
				defer func() {
					if err := recover(); err != nil {
						ctx.App.Log.Errorf("%s: %s", err, debug.Stack())
						http.Error(ctx, http.StatusText(500), 500)
					}
				}()

				return h.Serve(ctx)
			})
	}
}

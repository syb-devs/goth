package httptreemux

import (
	"net/http"

	"github.com/syb-devs/goth/app"

	"github.com/dimfeld/httptreemux"
)

// Muxer resolves requests to the corresponding HTTP handler
type Muxer struct {
	*httptreemux.TreeMux
	ctxGen app.CtxGenHTTP
}

// New returns a Muxer object
func New(ctxGen app.CtxGenHTTP) *Muxer {
	return &Muxer{
		ctxGen:  ctxGen,
		TreeMux: httptreemux.New(),
	}
}

// Handle registers an HTTP handler to a given verb / path combination
func (rt *Muxer) Handle(verb, path string, h app.Handler) {
	rt.TreeMux.Handle(verb, path, rt.wrapHandler(h))
}

func (rt *Muxer) wrapHandler(h app.Handler) httptreemux.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, urlParams map[string]string) {
		ctx := rt.ctxGen(w, r)
		ctx.URLParams = app.URLParams(urlParams)
		h.Serve(ctx)
	}
}

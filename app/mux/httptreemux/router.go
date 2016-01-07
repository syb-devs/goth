package httptreemux

import (
	"net/http"

	"bitbucket.org/syb-devs/goth/app"

	"github.com/dimfeld/httptreemux"
	"github.com/gorilla/context"
)

// Router resolves requests to the corresponding HTTP handler
type Router struct {
	*httptreemux.TreeMux
}

// New returns a Router object
func New() *Router {
	return &Router{
		TreeMux: httptreemux.New(),
	}
}

// Handle registers an HTTP handler to a given verb / path combination
func (rt *Router) Handle(verb, path string, h app.Handler) {
	rt.TreeMux.Handle(verb, path, rt.wrapHandler(h))
}

func (rt *Router) wrapHandler(h app.Handler) httptreemux.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, urlParams map[string]string) {
		context.Set(r, "url_params", app.URLParams(urlParams))
		h.ServeHTTP(w, r, nil)
	}
}

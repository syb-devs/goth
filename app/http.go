package app

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"net/textproto"
	"runtime"
	"strings"
)

// HandlerFunc represents a function that will serve an HTTP request
type HandlerFunc func(http.ResponseWriter, *http.Request, *Context) error

// ServeHTTP implements the Handler interface for the HandlerFunc type
func (hf HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request, ctx *Context) error {
	return hf(w, r, ctx)
}

// Handler interface extends the standard HTTP Handler injecting
// request context and acknowledging errors
type Handler interface {
	ServeHTTP(http.ResponseWriter, *http.Request, *Context) error
}

// URLParams contains variables from the URL path
type URLParams map[string]string

// ByName returns the value of the URL parameter with a given name
func (p URLParams) ByName(name string) string {
	return p[name]
}

// Router interface is used to bind a handler to a route
type Router interface {
	http.Handler
	Handle(verb, path string, h Handler)
}

// RouteProvider is an object that registers routes to handlers in a Router
type RouteProvider interface {
	RegisterRoutes(Router) error
}

func rootHandler(w http.ResponseWriter, r *http.Request, ctx *Context) error {
	return WriteJSON(w, map[string]string{"message": fmt.Sprintf("hello from %s", ctx.App.Name())})
}

// wrapHandler converts one of our custom Handlers to the Go HTTP Handler standard
// injecting the context and handling the error
func wrapHandler(f Handler, app *App) http.Handler {
	fmt.Println("wrapping handler...")
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := app.newContext(w, r)
		defer ctx.Close()

		err := f.ServeHTTP(w, r, ctx)
		if err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			stackBuff := make([]byte, 14096)
			runtime.Stack(stackBuff, false)
			fmt.Printf("Error serving %v: %v\n", r.URL, err)
			fmt.Printf("Stack dump %s\n", stackBuff)
		}
	}

	return http.HandlerFunc(fn)
}

func addContext(f Handler, app *App) Handler {
	return HandlerFunc(
		func(w http.ResponseWriter, r *http.Request, _ *Context) error {
			return f.ServeHTTP(w, r, app.newContext(w, r))
		})
}

// adaptHandler will ignore the context and the return error
// is used to convert a standard Go HTTP Handler to our augmented Handler interface
func adaptHandler(h http.Handler) Handler {
	fn := func(w http.ResponseWriter, r *http.Request, _ *Context) error {
		h.ServeHTTP(w, r)
		return nil
	}
	return HandlerFunc(fn)
}

// WriteJSON writes the given data marshaled as JSON, setting the appropiate
// content type header
func WriteJSON(w http.ResponseWriter, d interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(d)
}

// ReadJSON reads the request body and attempts to unmarshal it in the given destination
func ReadJSON(r *http.Request, dest interface{}) error {
	return json.NewDecoder(r.Body).Decode(dest)
}

type HTTPHeadersText string

func (t HTTPHeadersText) Parse() (http.Header, error) {
	headers := string(t)
	headers = strings.TrimRight(headers, "\r\n") + "\r\n\r\n"
	reader := bufio.NewReader(strings.NewReader(headers))
	tp := textproto.NewReader(reader)

	mimeHeader, err := tp.ReadMIMEHeader()
	if err != nil {
		return nil, err
	}

	return http.Header(mimeHeader), nil
}

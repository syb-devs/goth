package app

import (
	"net/http"
)

// Handler interface is implemented by the App handlers, which are Context aware
// and can return errors to be handled by the App
type Handler interface {
	Serve(*Context) error
}

// HandlerFunc represents a function that will serve an HTTP request
type HandlerFunc func(*Context) error

// Serve implements the Handler interface for the HandlerFunc type
func (hf HandlerFunc) Serve(ctx *Context) error {
	return hf(ctx)
}

// URLParams contains variables from the URL path
type URLParams map[string]string

// ByName returns the value of the URL parameter with a given name
func (p URLParams) ByName(name string) string {
	return p[name]
}

// Muxer interface is used to bind a handler to a route
type Muxer interface {
	http.Handler
	Handle(verb, path string, h Handler)
}

// Middleware is a function that wraps a handler and performs some tasks
// before and/or after calling the wrapped handler
type Middleware func(Handler) Handler

// MiddlewareChain interface is satisfied by objects that can build middleware chains
type MiddlewareChain interface {
	Append(...Middleware) MiddlewareChain
	Finally(Handler) Handler
}

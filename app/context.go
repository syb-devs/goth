package app

import (
	"io"
	"net/http"

	"bitbucket.org/syb-devs/goth/database"
	"bitbucket.org/syb-devs/goth/encoding"
	"bitbucket.org/syb-devs/goth/kv"
)

// CtxGenHTTP is a function that generates contexts from a pair
// of HTTP Request and ResponseWriter
type CtxGenHTTP func(w http.ResponseWriter, r *http.Request) *Context

// Context represents the isolated context for one request
type Context struct {
	App     *App
	Conn    database.Connection
	Request *http.Request
	http.ResponseWriter
	URLParams URLParams
	Codec     encoding.Codec
	UserID    string
	*kv.Store
}

// Close performs clean-up tasks for the Context
func (ctx *Context) Close() {
	ctx.Conn.Close()
}

// WriteString writes the given string in the Context's ResponseWriter
func (ctx *Context) WriteString(s string) (int, error) {
	return io.WriteString(ctx.ResponseWriter, s)
}

// URLParam returns the value for the requested URL parameter
func (ctx *Context) URLParam(name string) string {
	return ctx.URLParams.ByName(name)
}

// Decode decodes data from the context request into the destination type
func (ctx *Context) Decode(dest interface{}) error {
	return ctx.Codec.Decode(ctx.Request.Body, dest)
}

// Encode encodes the given data to the context response
func (ctx *Context) Encode(data interface{}) error {
	return ctx.Codec.Encode(ctx, data)
}

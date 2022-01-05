package app

import (
	"net/http"

	"github.com/syb-devs/goth/database"
)

// Context represents the isolated context for one request
type Context struct {
	App       *App
	Conn      database.Connection
	Request   *http.Request
	URLParams URLParams
	UserID    string
}

// Close performs clean-up tasks for the Context
func (ctx *Context) Close() {
	ctx.App.Close()
	ctx.Conn.Close()
}

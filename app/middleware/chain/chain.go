package chain

import (
	"github.com/syb-devs/goth/app"
)

// Chain is an ordered collection of Middlewares
type Chain struct {
	mws []app.Middleware
}

// New allocates a new Middleware Chain with the given Middlewares and returns it
func New(mws ...app.Middleware) Chain {
	return Chain{
		mws: mws,
	}
}

// Append adds another middleware to the Chain
func (c Chain) Append(mw ...app.Middleware) app.MiddlewareChain {
	return c
}

// Finally returns the given Handler wrapped inside of the chain middlewares
func (c Chain) Finally(h app.Handler) app.Handler {
	if h == nil {
		panic("Chain.Finally requires a valid app.Handler, nil given")
	}

	for i := len(c.mws) - 1; i >= 0; i-- {
		mw := c.mws[i]
		h = mw(h)
	}
	return h
}

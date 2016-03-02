package main

import (
	"bitbucket.org/syb-devs/goth/app"
)

func rootHandler(ctx *app.Context) error {
	return ctx.Encode(struct {
		Message string `json:"message"`
	}{"Goth Example API v0.3.0"})
}

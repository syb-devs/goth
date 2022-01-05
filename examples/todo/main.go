package main

import (
	"fmt"
	"os"

	"github.com/syb-devs/goth/app"
	"github.com/syb-devs/goth/app/middleware/buffer"
	"github.com/syb-devs/goth/app/middleware/chain"
	"github.com/syb-devs/goth/app/middleware/recovr"
	"github.com/syb-devs/goth/app/middleware/timer"
	"github.com/syb-devs/goth/app/mux/httptreemux"
	"github.com/syb-devs/goth/database"
	"github.com/syb-devs/goth/database/driver/mongodb"
	"github.com/syb-devs/goth/log"
	"github.com/syb-devs/goth/rest"
	"github.com/syb-devs/goth/user"

	"github.com/syb-devs/dockerlink"

	"github.com/rs/cors"
)

func main() {
	myApp := app.NewApp("Goth example App v0.1.0")

	// HTTP setup
	mux := httptreemux.New(myApp.NewContextHTTP)
	myApp.SetMuxer(mux)

	corsOpts := cors.Options{
		Debug:          false,
		AllowedHeaders: []string{"*"},
	}
	handler := cors.New(corsOpts).Handler(mux)
	myApp.SetHandler(handler)

	mainChain := chain.New(
		buffer.New(),
		recovr.New(),
		timer.New(),
		errMiddleware,
	)
	myApp.AddChain(mainChain, "main")
	myApp.AddChain(mainChain, "pub")

	myApp.Handle("GET", "/", myApp.WrapHandlerFunc(rootHandler, "main"))

	ps := database.ConnectionParams{
		"url":      getMongoURI(),
		"database": "goth",
	}
	conn, err := mongodb.NewConnection(ps, database.NewResourceMap())
	if err != nil {
		panic(err)
	}
	myApp.DB.Connection = conn
	myApp.DB.ResourceMap = conn.Map()
	myApp.DB.Repository = mongodb.NewRepository(conn)

	myApp.DB.RegisterResource(Todo{}, "todos", "")
	myApp.DB.RegisterResource(User{}, "users", "")
	myApp.DB.RegisterResource(Profile{}, "profiles", "")

	user.RegisterType(&User{}, "username")

	rest.Register(myApp, &Todo{}, "todos")
	rest.Register(myApp, &User{}, "users")
	rest.Register(myApp, &Profile{}, "profiles")

	createData(myApp)
	myApp.Run()
}

func getMongoURI() string {
	if uri := os.Getenv("MONGO_URL"); uri != "" {
		return uri
	}
	if link, err := dockerlink.GetLink("mongodb", 27017, "tcp"); err == nil {
		return fmt.Sprintf("%s:%d", link.Address, link.Port)
	}
	panic("mongodb connection not found, use MONGO_URL env var or a docker link with mongodb name")
}

func errMiddleware(h app.Handler) app.Handler {
	f := func(ctx *app.Context) error {
		err := h.Serve(ctx)
		if err != nil {
			log.Errorf(
				"error serving %s %s: %v",
				ctx.Request.Method,
				ctx.Request.URL.String(),
				err)
		}
		return nil
	}
	return app.HandlerFunc(f)
}

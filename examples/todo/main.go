package main

import (
	"fmt"
	"os"

	"bitbucket.org/syb-devs/goth/app"
	"bitbucket.org/syb-devs/goth/app/middleware/buffer"
	"bitbucket.org/syb-devs/goth/app/middleware/chain"
	"bitbucket.org/syb-devs/goth/app/middleware/cors"
	"bitbucket.org/syb-devs/goth/app/middleware/recovr"
	"bitbucket.org/syb-devs/goth/app/middleware/timer"
	"bitbucket.org/syb-devs/goth/app/mux/httptreemux"
	"bitbucket.org/syb-devs/goth/database"
	"bitbucket.org/syb-devs/goth/database/driver/mongodb"
	"bitbucket.org/syb-devs/goth/rest"
	// _ "bitbucket.org/syb-devs/goth/user"

	"github.com/syb-devs/dockerlink"
)

func main() {
	myApp := app.NewApp("Goth example App v0.1.0")

	// HTTP setup
	myApp.SetMuxer(httptreemux.New(myApp.NewContextHTTP))

	corsOpts := cors.Options{
		Debug:          false,
		AllowedHeaders: []string{"*"},
	}
	mainChain := chain.New(
		buffer.New(),
		recovr.New(),
		timer.New(),
		cors.New(corsOpts).Handler,
	)
	myApp.AddChain(mainChain, "main")

	myApp.Handle("GET", "/", myApp.WrapHandlerFunc(rootHandler, "main"))
	myApp.Handle("GET", "/hello", myApp.WrapHandlerFunc(helloJSONHandler, "main"))
	myApp.Handle("GET", "/debug/vars", myApp.WrapHandlerFunc(expvarHandler, "main"))
	// myApp.Handle("GET", "/users", myApp.WrapHandlerFunc(listUsers, "main"))

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

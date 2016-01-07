package main

import (
	"fmt"
	"os"

	"bitbucket.org/syb-devs/goth/app"
	"bitbucket.org/syb-devs/goth/app/mux/httptreemux"
	"bitbucket.org/syb-devs/goth/database"
	"bitbucket.org/syb-devs/goth/database/driver/mongodb"
	_ "bitbucket.org/syb-devs/goth/user"

	// jwt "github.com/auth0/go-jwt-middleware"
	"github.com/rs/cors"
	"github.com/syb-devs/dockerlink"
)

func main() {
	app := app.NewApp("Goth example app v0.1")

	ps := database.ConnectionParams{
		"url":      getMongoURI(),
		"database": "stock_holmes",
	}
	conn, err := mongodb.NewConnection(ps, database.NewResourceMap())
	if err != nil {
		panic(err)
	}

	app.DB.Connection = conn
	app.DB.ResourceMap = conn.Map()
	app.DB.Repository = mongodb.NewRepository(conn)

	corsOpts := cors.Options{
		Debug:          false,
		AllowedHeaders: []string{"*"},
	}
	// jwtOpts := jwt.Options{
	// 	Debug:               false,
	// 	ValidationKeyGetter: user.JWTKeyFunc,
	// }

	// app.Use(jwt.New(jwtOpts).Handler, "jwt")
	app.UseGlobal(cors.New(corsOpts).Handler)

	// Setup the Websockets server
	// wss := core.NewWSServer(app)
	// wss.LoadEventHandlers()
	//
	// sm := core.SkipMiddleware
	// app.Handle("GET", "/ws", sm(wss, "jwt"))

	app.SetRouter(httptreemux.New())
	app.Run()
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

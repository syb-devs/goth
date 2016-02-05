package main

import (
	"expvar"
	"fmt"

	"bitbucket.org/syb-devs/goth/app"
	"bitbucket.org/syb-devs/goth/database"
	"bitbucket.org/syb-devs/goth/database/driver/mongodb"
	"bitbucket.org/syb-devs/goth/log"
	// _ "bitbucket.org/syb-devs/goth/user"
)

func rootHandler(ctx *app.Context) error {
	ctx.Header().Set("Foo", "bar")
	ctx.WriteString("rootHandler\n")
	return nil
}

func helloJSONHandler(ctx *app.Context) error {
	// panic("mi abuela fuma")

	return ctx.Encode(struct {
		Message string `json:"message"`
	}{"hello!"})
}

func expvarHandler(ctx *app.Context) error {
	w := ctx
	ctx.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprintf(w, "{\n")
	first := true
	expvar.Do(func(kv expvar.KeyValue) {
		if !first {
			fmt.Fprintf(w, ",\n")
		}
		first = false
		fmt.Fprintf(w, "%q: %s", kv.Key, kv.Value)
	})
	fmt.Fprintf(w, "\n}\n")
	return nil
}

func listUsers(ctx *app.Context) error {
	users := []*User{}
	err := ctx.App.DB.FindMany(&users, database.NewQ(nil))
	if err != nil {
		return err
	}
	for _, user := range users {
		err = ctx.App.DB.Repository.(*mongodb.Repository).FetchRelated(user, "todos", "profile")
		if err != nil {
			log.Error(err)
			return err
		}
	}
	return ctx.Encode(users)
}

package main

import (
	"bitbucket.org/syb-devs/goth/app"
	"bitbucket.org/syb-devs/goth/database"
	"bitbucket.org/syb-devs/goth/log"

	"gopkg.in/mgo.v2/bson"
)

func createData(a *app.App) error {
	users := []*User{}
	err := a.DB.FindMany(&users, database.NewQ(nil))
	if err != nil {
		return err
	}
	if len(users) > 0 {
		log.Debug("user and todos already created")
		return nil
	}

	log.Debug("creating user and todos...")

	todos := []Todo{
		{Title: "Learn Go", Done: true},
		{Title: "Learn Rust"},
		{Title: "Learn DSP"},
		{Title: "Learn C++"},
	}

	u := &User{}
	u.Name = "John"

	for _, todo := range todos {
		err = a.DB.Insert(&todo)
		if err != nil {
			return err
		}
		log.Debugf("created todo: %+v", todo)
		u.TodoIDs = append(u.TodoIDs, todo.GetID().(bson.ObjectId))
	}

	p := &Profile{
		Twitter:  "@JohnDoe",
		Facebook: "https://www.facebook.com/johndoe",
		Linkedin: "https://www.linkedin.com/in/john-doe",
	}
	err = a.DB.Insert(p)
	if err != nil {
		return err
	}
	log.Debugf("created profile: %+v", p)
	u.ProfileID = &p.ID

	err = a.DB.Insert(u)
	if err != nil {
		return err
	}
	log.Debugf("created user: %+v", u)
	return nil
}

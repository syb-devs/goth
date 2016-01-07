package user

import (
	"bitbucket.org/syb-devs/goth/database"
	"bitbucket.org/syb-devs/goth/database/driver/mongodb"
	"github.com/syb-devs/gotools/auth"
)

// User represents a user of the app
type User struct {
	mongodb.Resource `bson:",inline"`
	auth.Auth        `bson:",inline"`
	FirstName        string `bson:"firstName" json:"firstName"`
	LastName         string `bson:"lastName" json:"lastName"`
}

type userDispatcher struct{}

func (d *userDispatcher) New() database.Resource         { return &User{} }
func (d *userDispatcher) NewList() database.ResourceList { return &[]*User{} }

// QueryByUsername returns a user, querying by username, or error if not found
func QueryByUsername(username string) database.Dict {
	return database.Dict{"username": username}
}

package user

import (
	"github.com/syb-devs/goth/app"
	"github.com/syb-devs/goth/database"
	"github.com/syb-devs/goth/rest"
)

func init() {
	app.RegisterModule(&module{BaseModule: app.NewBaseModule("user")})
}

type module struct {
	*app.BaseModule
}

// Bootstrap performs initialization tasks, such as registering resources and HTTP routes
func (m *module) Bootstrap(app *app.App, level int) error {
	if level != 5 {
		return nil
	}
	app.RegisterResource(&User{}, database.ResourceConfig{
		Name:            "user",
		ColName:         "users",
		ArchiveOnDelete: true,
	})
	rest.RegisterResource(app, rest.ResourceConfig{
		Name:    "user",
		URLName: "users",
		Handler: newUserHandler(),
	})
	return nil
}

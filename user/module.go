package user

import (
	"bitbucket.org/syb-devs/goth/app"
)

func init() {
	app.RegisterModule(&module{BaseModule: app.NewBaseModule("goth.user")})
}

type module struct {
	*app.BaseModule
}

// Bootstrap performs initialization tasks, such as registering resources and HTTP routes
func (m *module) Bootstrap(a *app.App, level int) error {
	if level != 5 {
		return nil
	}
	a.Handle("POST", "/users/register", a.WrapHandlerFunc(register, "pub"))
	a.Handle("POST", "/users/sessions", a.WrapHandlerFunc(login, "pub"))

	return nil
}

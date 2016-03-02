package user

import (
	"errors"

	"bitbucket.org/syb-devs/goth/app"
	"bitbucket.org/syb-devs/goth/auth/jwt"
	"bitbucket.org/syb-devs/goth/database"
)

var (
	// ErrEmptyUserPass happens when no username and/or password is given for a user
	ErrEmptyUserPass = errors.New("username and/or password not set")

	// ErrInvalidUserPass happens when no  valid username and/or password is given for a user
	ErrInvalidUserPass = errors.New("username and/or password is not valid")
)

type retJWT struct {
	Token string `json:"token"`
}

func register(ctx *app.Context) error {
	user := newUser(ctx)
	err := ctx.Decode(user)
	if err != nil {
		return err
	}
	user.SetPassword(user.GetPassword())

	if err = ctx.App.DB.Insert(user); err != nil {
		return err
	}
	return ctx.Encode(user)
}

func login(ctx *app.Context) error {
	loginData := &struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}
	err := ctx.Decode(loginData)
	if err != nil {
		return err
	}
	if loginData.Username == "" || len(loginData.Password) == 0 {
		return ErrEmptyUserPass
	}
	user := newUser(ctx).(Interface)
	err = ctx.App.DB.FindOne(user, database.NewQ(
		database.Dict{usernameDBField: loginData.Username}))
	if err != nil {
		return err
	}
	if err = user.CheckAuth(loginData.Username, []byte(loginData.Password)); err != nil {
		return ErrInvalidUserPass
	}
	token, err := jwt.New(user, nil)
	if err != nil {
		return err
	}
	return ctx.Encode(retJWT{token})
}

func newUser(ctx *app.Context) Interface {
	return ctx.App.DB.CreateResource(userType).(Interface)
}

package user

import (
	"errors"
	"net/http"
	"time"

	"bitbucket.org/syb-devs/goth/app"
	"bitbucket.org/syb-devs/goth/database"
	"bitbucket.org/syb-devs/goth/rest"
)

var (
	ErrEmptyUserPass   = errors.New("username and/or password not set")
	ErrInvalidUserPass = errors.New("username and/or password is not valid")
)

var jwtExp = 24 * time.Hour

type userHandler struct {
	rest.ResourceHandler
	database.ResourceCodec
	database.ResourceValidator
}

func newUserHandler() *userHandler {
	baseHandler := rest.NewDefResourceHandler(&userDispatcher{})
	return &userHandler{
		ResourceHandler:   baseHandler,
		ResourceValidator: baseHandler,
		ResourceCodec:     baseHandler,
	}
}

func (h *userHandler) RegisterRoutes(r app.Router) error {
	sm := app.SkipMiddleware

	r.Handle("POST", "/users/register", sm(app.HandlerFunc(h.register), "jwt"))
	r.Handle("POST", "/users/sessions", sm(app.HandlerFunc(h.login), "jwt"))
	return nil
}

type retJWT struct {
	Token string `json:"token"`
}

func (h *userHandler) register(w http.ResponseWriter, r *http.Request, ctx *app.Context) error {
	rUser, err := h.userFromRequest(r)
	if err != nil {
		return err
	}

	user := &User{}
	user.SetUserName(rUser.Username)
	user.GeneratePassword([]byte(rUser.Password))

	if err = ctx.App.DB.Insert(user); err != nil {
		return err
	}
	jwt, err := newJWT(user, jwtExp)
	if err != nil {
		return err
	}
	return h.Encode(w, retJWT{jwt})
}

func (h *userHandler) login(w http.ResponseWriter, r *http.Request, ctx *app.Context) error {
	rUser, err := h.userFromRequest(r)
	if err != nil {
		return err
	}
	user := &User{}
	query := database.NewQ(QueryLogin(rUser.Username))
	err = ctx.App.DB.FindOne(user, query)
	if err != nil {
		return err
	}
	if err = user.Check(rUser.Username, []byte(rUser.Password)); err != nil {
		return ErrInvalidUserPass
	}
	jwt, err := newJWT(user, jwtExp)
	if err != nil {
		return err
	}
	return h.Encode(w, retJWT{jwt})
}

type reqUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *userHandler) userFromRequest(r *http.Request) (*reqUser, error) {
	ru := &reqUser{}
	if err := h.Decode(r.Body, ru); err != nil {
		return nil, err
	}
	if ru.Username == "" || len(ru.Password) == 0 {
		return nil, ErrEmptyUserPass
	}
	return ru, nil
}

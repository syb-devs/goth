package user

import (
	"log"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/syb-devs/goth/app"
)

func init() {
	app.RegisterWSBindFunc(func(svr *app.WSServer) error {
		svr.Bind("auth", checkWSJWT)
		return nil
	})
}

// TODO: retrieve from env, config in database...
var jwtSecret = []byte("jander_klander")

func newJWT(user *User, exp time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims["exp"] = time.Now().Add(exp).Unix()
	token.Claims["user_id"] = user.ID.Hex()

	return token.SignedString([]byte(jwtSecret))
}

// JWTKeyFunc is used to get the key used to sign the JSON Web Tokens
func JWTKeyFunc(t *jwt.Token) (interface{}, error) {
	return jwtSecret, nil
}

func checkWSJWT(ws *app.WSConn, e *app.WSEvent) error {
	data := &struct {
		Token string
	}{}

	err := e.DecodeData(data)
	if err != nil {
		return err
	}

	token, err := jwt.Parse(data.Token, JWTKeyFunc)
	log.Printf("parsed token: %+v", token)

	if err != nil {
		return err
	}

	ws.UserID = token.Claims["user_id"].(string)
	ws.IsAuth = true

	return nil
}

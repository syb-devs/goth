package user

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

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

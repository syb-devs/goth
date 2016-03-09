package jwt

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var (
	// Secret is the key used to sign the JWT tokens
	Secret []byte
	// Duration used for the token expiration
	Duration time.Duration
)

type user interface {
	GetID() interface{}
}

// New returns an encoded and signed JWT Token string with the given claims
func New(user user, claims map[string]interface{}) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims["exp"] = time.Now().Add(Duration).Unix()
	token.Claims["userId"] = user.GetID()
	if claims != nil {
		for claim, val := range claims {
			token.Claims[claim] = val
		}
	}
	return token.SignedString(Secret)
}

// GetKeyFunc is used to get the key used to sign the JSON Web Tokens
func GetKeyFunc(t *jwt.Token) (interface{}, error) {
	return Secret, nil
}

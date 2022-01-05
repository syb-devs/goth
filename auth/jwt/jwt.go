package jwt

import (
	"errors"
	"reflect"
	"time"

	"github.com/syb-devs/goth/app"
	"github.com/syb-devs/goth/database"

	jwt "gopkg.in/dgrijalva/jwt-go.v2"
)

var (
	// Secret is the key used to sign the JWT tokens
	Secret []byte
	// Duration used for the token expiration
	Duration = 7 * 24 * time.Hour

	// ErrInvalidUserID is returned when a user ID is not found
	ErrInvalidUserID = errors.New("no valid user ID found on request context store")
)

const userIDField = "userId"

type user interface {
	GetID() interface{}
}

// New returns an encoded and signed JWT Token string with the given claims
func New(user user, claims map[string]interface{}) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims["exp"] = time.Now().Add(Duration).Unix()
	token.Claims[userIDField] = user.GetID()
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

// TokenProcessorOptions is a struct for specifying configuration options for the middleware.
type TokenProcessorOptions struct {
	UserType interface{}
}

// NewTokenProcessorFunc allocates and returns a TokenProcessor function
// which extracts the user ID from theparsed token, retrieves the user
// from database and stores it in the context
func NewTokenProcessorFunc(options ...TokenProcessorOptions) func(ctx *app.Context, token *jwt.Token) error {
	var opts TokenProcessorOptions
	if len(options) == 0 {
		opts = TokenProcessorOptions{}
	} else {
		opts = options[0]
	}
	userType := reflect.TypeOf(opts.UserType)

	return func(ctx *app.Context, token *jwt.Token) error {
		userID := token.Claims[userIDField]
		if userID == nil {
			return ErrInvalidUserID
		}
		user := ctx.App.DB.CreateResource(userType).(database.Resource)
		err := ctx.App.DB.Get(userID, user)
		if err != nil {
			return err
		}
		ctx.User = user.(app.User)
		return nil
	}
}

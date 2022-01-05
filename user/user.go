package user

import (
	"errors"
	"reflect"

	"golang.org/x/crypto/bcrypt"

	"github.com/syb-devs/goth/database"
)

var (
	// ErrInvalidUserName is returned when an empty username is given for login or signup
	ErrInvalidUserName = errors.New("invalid username")
	// ErrInvalidPassword is returned when an empty password is given for login or signup
	ErrInvalidPassword = errors.New("invalid password")

	userType        reflect.Type
	usernameDBField string
)

// RegisterType registers the type used in the user handlers and the corresponding
// query for retrieving it from the database
func RegisterType(user Interface, usernameField string) {
	userType = reflect.TypeOf(user)
	usernameDBField = usernameField
}

// Interface must be satisfied by the registered user type
type Interface interface {
	database.Resource
	SetPassword([]byte) error
	GetPassword() []byte
	CheckAuth(username string, pass []byte) error
}

// User represents a user of the app
type User struct {
	Username string   `bson:"username" json:"username"`
	Password password `bson:"password" json:"password"`
}

// CheckPassword checks password for the User
func (u *User) CheckPassword(password []byte) error {
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), password); err != nil {
		return ErrInvalidPassword
	}
	return nil
}

// CheckAuth checks the username and password for the user
func (u *User) CheckAuth(username string, password []byte) error {
	if u.Username != username {
		return ErrInvalidUserName
	}
	return u.CheckPassword(password)
}

type password []byte

func (p password) MarshalJSON() ([]byte, error) {
	str := "\"******\""
	return []byte(str), nil
}

func (p password) UnmarshalJSON(data []byte) error {
	p = password(data)
	return nil
}

// SetPassword encrypts the plain password and sets it to the User
func (u *User) SetPassword(plain []byte) error {
	hash, err := bcrypt.GenerateFromPassword(plain, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = password(hash)
	return nil
}

// GetPassword returns the password of the User
func (u *User) GetPassword() []byte {
	return []byte(u.Password)
}

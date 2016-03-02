package app

// User interface identifies a user of the app
type User interface {
	GetID() interface{}
	GetIDString() string
	GetEmail() string
}

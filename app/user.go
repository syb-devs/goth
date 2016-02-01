package app

// User interface identifies a user of the app
type User interface {
	GetID() string
	GetEmail() string
}

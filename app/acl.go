package app

import (
	db "github.com/syb-devs/goth/database"
)

// ResourceAccessChecker interface is used for abstracting access control over resources
type ResourceAccessChecker interface {
	UserAllowed(ctx *Context, userID db.ResourceOwner, res db.Resource, action string) (bool, error)
}

// AccessChecker is a default implementation of ResourceAccessChecker interface
type AccessChecker struct{}

// UserAllowed checks if the user is allowed to perform the given action over the given resource
func (a AccessChecker) UserAllowed(ctx *Context, user db.ResourceOwner, res db.Resource, action string) (bool, error) {
	//TODO(zareone) check if user is admin
	return res.BelongsTo(user), nil
}

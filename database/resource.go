package database

import (
	"errors"
	"reflect"
)

const (
	ResourceActionCreate = iota
	ResourceActionRetrieve
	ResourceActionUpdate
	ResourceActionDelete
	ResourceActionList
)

var (
	ErrSliceExpected = errors.New("expected slice of database.Resource")
)

// Resource represents an entity that can be persisted to database.
type Resource interface {
	GetID() interface{}
	SetID(interface{})
	GetOwnerID() interface{}
	SetOwnerID(interface{})
	BelongsTo(ResourceOwner) bool
}

// ResourceOwner represents the owner of a Resource
type ResourceOwner interface {
	GetID() interface{}
}

// ResourceList is an slice of Resource.
type ResourceList interface{}

// ResourceValidator interface is used to check the validity of a given resource
// level is used when the validation is dependent of the current action (creating, updating...)
type ResourceValidator interface {
	Validate(level int, res interface{}) error
}

// DummyResourceValidator implements the ResourceValidator interface
// but does not actually validate anything
type DummyResourceValidator struct{}

// Validate implements the ResourceValidator interface, but it's a dummy function
func (v *DummyResourceValidator) Validate(level int, res interface{}) error {
	return nil
}

// Toucher is the interface implemented by an object that can update its
// modification timestamps.
type Toucher interface {
	Touch()
}

// SoftDeletable is the interface implemented by an object that can be marked
// for logical deletion.
type SoftDeletable interface {
	MarkDeleted()
}

// ResourceConfig has the needed info to Register a resource in the App
type ResourceConfig struct {
	Name            string
	ColName         string
	ArchiveOnDelete bool
}

// CheckResourceList checks that a ResourceList is a slice of types that conform
// to the Resource interface
func CheckResourceList(l ResourceList) error {
	//TODO(zareone) check that the list:
	// - Is actually a slice or pointer to one
	// - The type of the slice implements the Resource interface
	return nil
}

// AsResourceList returns a slice of Resource from a slice
// of any type that implements the Resource interface
func AsResourceList(list interface{}) ([]Resource, error) {
	var ret []Resource
	v := reflect.ValueOf(list)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Slice {
		return nil, ErrSliceExpected
	}
	for i := 0; i < v.Len(); i++ {
		ret = append(ret, v.Index(i).Interface().(Resource))
	}
	return ret, nil
}

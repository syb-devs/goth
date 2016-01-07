package database

import (
	"encoding/json"
	"io"
	"net/http"
)

const (
	ResourceActionCreate = iota
	ResourceActionRetrieve
	ResourceActionUpdate
	ResourceActionDelete
	ResourceActionList
)

// Resource represents an entity that can be persisted to database.
type Resource interface {
	GetID() interface{}
	SetID(interface{}) error
}

// ResourceList is an slice of Resource.
type ResourceList interface{}

// ResourceConstructor interface is used to instance new resources and resource lists
type ResourceConstructor interface {
	New() Resource
	NewList() ResourceList
}

// ResourceCodec interface is used to encode and decode resources between the App
// and the transport layers (HTTP, websockets)
type ResourceCodec interface {
	Encode(io.Writer, interface{}) error
	Decode(io.Reader, interface{}) error
}

// JSONResourceCodec encodes and decodes resources to and from JSON data
type JSONResourceCodec struct{}

// Encode marshals a resource (r) as JSON data in the given writer (w)
func (c *JSONResourceCodec) Encode(w io.Writer, res interface{}) error {
	if rw, ok := w.(http.ResponseWriter); ok {
		rw.Header().Set("Content-Type", "application/json")
	}
	return json.NewEncoder(w).Encode(res)
}

// Decode unmarshalls JSON data read from the Reader (r) into res
func (c *JSONResourceCodec) Decode(r io.Reader, res interface{}) error {
	return json.NewDecoder(r).Decode(res)
}

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

package database

import (
	"errors"
	"reflect"
	"sync"
)

var (
	// ErrTypeNotRegistered is returned when asking data for a type which has not been registered in the map
	ErrTypeNotRegistered = errors.New("type not registered")
	// ErrNilResource is returned when a nil pointer is given instead of a valid Resource
	ErrNilResource = errors.New("model type must be addressable (not nil)")
	// ErrInvalidColName is returned when registering a Resource without a valid collection name (empty string)
	ErrInvalidColName = errors.New("a collection name must be given")
	// ErrResourceNotStruct is returned when a non-struct data type is given in place of a Resource
	ErrResourceNotStruct = errors.New("resources must be of type struct")
	// ErrRelationshipNotFound is returned when a specified relationship can not be found
	ErrRelationshipNotFound = errors.New("relationship not defined")
)

// ResourceMap maintains a relationship of Resource types and their respective table/collections for persistence
type ResourceMap struct {
	sync.Mutex
	cols  map[reflect.Type]ResourceMapItem
	types map[string]reflect.Type
	rels  map[reflect.Type][]Relationship
}

// NewResourceMap returns a new ResourceMap
func NewResourceMap() *ResourceMap {
	return &ResourceMap{
		cols:  make(map[reflect.Type]ResourceMapItem, 0),
		types: make(map[string]reflect.Type, 0),
		rels:  make(map[reflect.Type][]Relationship, 0),
	}
}

// RegisterResource registers a Resource type to a table / collection
func (m *ResourceMap) RegisterResource(t interface{}, col string, deletedCol string) error {
	m.Lock()
	defer m.Unlock()

	tt, err := m.toType(t)
	if err != nil {
		return err
	}
	if col == "" {
		return ErrInvalidColName
	}
	m.cols[tt] = ResourceMapItem{col, deletedCol}
	m.types[tt.String()] = tt
	m.rels[tt] = getRelationships(t)
	return nil
}

// CreateResource allocates and returns object of the given type
func (m *ResourceMap) CreateResource(t reflect.Type) interface{} {
	return reflect.New(t).Elem().Interface()
}

// CreateResourceList allocates and returns a list of objects of the given type
func (m *ResourceMap) CreateResourceList(t reflect.Type) interface{} {
	slice := reflect.MakeSlice(reflect.SliceOf(t), 0, 0)

	// Create a pointer to a slice value and set it to the slice
	ptr := reflect.New(slice.Type())
	ptr.Elem().Set(slice)
	return ptr.Interface()
}

// TypeFromString returns the reflect.Type for a registered type from its full name (package.Type)
func (m *ResourceMap) TypeFromString(name string) (t reflect.Type, ok bool) {
	t, ok = m.types[name]
	return
}

// Relationships returns the relationships defined for the given resource
func (m *ResourceMap) Relationships(t interface{}) []Relationship {
	tt, err := m.toType(t)
	if err != nil {
		panic(err)
	}
	return m.rels[tt]
}

// Relationship returns one relationship by its name
func (m *ResourceMap) Relationship(t interface{}, name string) (Relationship, error) {
	rels := m.Relationships(t)
	for _, rel := range rels {
		if rel.Name == name {
			return rel, nil
		}
	}
	return Relationship{}, ErrRelationshipNotFound
}

// ColFor returns the name of the collection / table registered for the given type
func (m *ResourceMap) ColFor(t interface{}) (string, error) {
	colData, err := m.findByType(t)
	if err != nil {
		return "", err
	}
	return colData.ColName, nil
}

// DeletedColFor returns the name of the collection / table registered
// for deleted resources of the given type
func (m *ResourceMap) DeletedColFor(t interface{}) (string, error) {
	colData, err := m.findByType(t)
	if err != nil {
		return "", err
	}
	return colData.DeletedColName, nil
}

func (m *ResourceMap) findByType(t interface{}) (ResourceMapItem, error) {
	var ret ResourceMapItem

	tt, err := m.toType(t)
	if err != nil {
		return ret, err
	}
	if c, ok := m.cols[tt]; ok {
		return c, nil
	}
	return ret, ErrTypeNotRegistered
}

func (m *ResourceMap) toType(i interface{}) (reflect.Type, error) {
	if i == nil {
		return nil, ErrNilResource
	}
	t := reflect.TypeOf(i)
	// If pointer or slice, get it's element
	for t.Kind() == reflect.Ptr || t.Kind() == reflect.Slice {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return nil, ErrResourceNotStruct
	}
	return t, nil
}

// ResourceMapItem has information about the collections for a given resource type
type ResourceMapItem struct {
	ColName        string
	DeletedColName string
}

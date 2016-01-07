package database

//TODO(zareone) Move to mongodb package, as the resource map is an implementation detail

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
)

var (
	ErrTypeNotRegistered = errors.New("type not registered")
	ErrNilResource       = errors.New("resource must be addressable (not nil)")
)

// ResourceMap maintains a relationship of Resource types and their respective table/collections for persistence
type ResourceMap struct {
	sync.Mutex
	cols map[reflect.Type]ResourceMapItem
}

// NewResourceMap returns a new ResourceMap
func NewResourceMap() *ResourceMap {
	return &ResourceMap{cols: make(map[reflect.Type]ResourceMapItem, 0)}
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
		return errors.New("a collection name must be given")
	}

	m.cols[tt] = ResourceMapItem{col, deletedCol}
	return nil
}

// ColFor returns the name of the collection / table registered for the given type
func (m *ResourceMap) ColFor(t interface{}) (string, error) {
	if t == nil {
		return "", ErrNilResource
	}
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
	t := reflect.TypeOf(i)
	// If pointer or slice, get it's element
	for t.Kind() == reflect.Ptr || t.Kind() == reflect.Slice {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("type should be struct: %v", reflect.TypeOf(i))
	}
	return t, nil
}

// ResourceMapItem has information about the collections for a given resource type
type ResourceMapItem struct {
	ColName        string
	DeletedColName string
}

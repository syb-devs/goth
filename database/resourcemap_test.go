package database_test

import (
	"reflect"
	"testing"

	"bitbucket.org/syb-devs/goth/database"
)

type User struct {
	name string
	age  int
}

func TestResourceMap(t *testing.T) {
	var tests = []struct {
		resource         interface{}
		colName          string
		deletedColName   string
		typeName         string
		expectedResource interface{}
	}{
		{User{}, "users", "archived_users", "database_test.User", &User{}},
		{&User{}, "users", "archived_users", "database_test.User", &User{}},
	}

	for _, test := range tests {
		rmap := database.NewResourceMap()
		err := rmap.RegisterResource(test.resource, test.colName, test.deletedColName)
		if err != nil {
			t.Errorf("RegisterResource: %v", err)
		}
		typ, _ := rmap.TypeFromString(test.typeName)
		resource := rmap.CreateResource(typ)
		if !sameType(test.expectedResource, resource) {
			t.Errorf("%T and %T have differing types", test.resource, resource)
		}
		col, err := rmap.ColFor(test.resource)
		if err != nil {
			t.Errorf("ColFor: %v", err)
		}
		if col != test.colName {
			t.Errorf("collection name, expected %s, got %s", test.colName, col)
		}
		deletedCol, err := rmap.DeletedColFor(resource)
		if err != nil {
			t.Errorf("DeletedColFor: %v", err)
		}
		if deletedCol != test.deletedColName {
			t.Errorf("deleted collection name, expected %s, got %s", test.deletedColName, deletedCol)
		}
	}
}

func TestRegisterResourceErrs(t *testing.T) {
	var tests = []struct {
		resource interface{}
		colName  string
		err      error
	}{
		{nil, "foo", database.ErrNilResource},
		{"something", "foo", database.ErrResourceNotStruct},
		{User{}, "", database.ErrInvalidColName},
	}

	for _, test := range tests {
		rmap := database.NewResourceMap()
		err := rmap.RegisterResource(test.resource, test.colName, "")
		if err != test.err {
			t.Errorf("expecting error %v, got %+v", test.err, err)
		}
	}
}

func TestRelationship(t *testing.T) {
	var tests = []struct {
		resource interface{}
		expected []database.Relationship
	}{}

	for _, test := range tests {
		rmap := database.NewResourceMap()
		err := rmap.RegisterResource(test.resource, "foo", "")
		if err != nil {
			t.Errorf("error %v", err)
			t.Fail()
		}
	}

}

func sameType(a, b interface{}) bool {
	return reflect.TypeOf(a) == reflect.TypeOf(b)
}

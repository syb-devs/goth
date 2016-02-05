package database

import (
	"reflect"
	"strings"
)

// These are the kind of relatinships available
const (
	HasZeroOne = iota
	HasOne
	HasMany
	BelongsToOne
	BelongsToMany

	RelTag = "rel"
)

// Relationship represents how two databse Resources relate to each other
type Relationship struct {
	Kind        int
	Name        string
	FieldName   string
	TargetField string
}

// RelationshipTargeter is used to extract the relationship target fields
// without using the reflect-based default implementation
type RelationshipTargeter interface {
	RelationshipTarget(relName string) (fieldref interface{})
}

func getRelationships(res interface{}) []Relationship {
	if res == nil {
		panic(ErrNilResource)
	}
	t := reflect.TypeOf(res)
	// If pointer, get it's element
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		panic(ErrResourceNotStruct)
	}
	rels := []Relationship{}
	numFields := t.NumField()
	for i := 0; i < numFields; i++ {
		f := t.Field(i)
		rel := relationshipFromField(f)
		if rel != nil {
			rels = append(rels, *rel)
		}
	}
	return rels
}

func relationshipFromField(f reflect.StructField) *Relationship {
	tag := f.Tag.Get(RelTag)
	if tag == "" {
		return nil
	}
	tagParts := strings.Split(tag, ",")
	if len(tagParts) < 2 {
		panic("model tag needs at least 2 comma separated values (name, target field)")
	}
	return &Relationship{
		Kind:        relKindFromField(f),
		Name:        tagParts[0],
		FieldName:   f.Name,
		TargetField: tagParts[1],
	}
}

func relKindFromField(f reflect.StructField) int {
	switch f.Type.Kind() {
	case reflect.Slice:
		return HasMany
	case reflect.Ptr:
		return HasZeroOne
	default:
		return HasOne
	}
}

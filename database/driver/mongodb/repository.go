package mongodb

import (
	"reflect"

	"bitbucket.org/syb-devs/goth/database"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var _ database.Repository = (*Repository)(nil)

// Repository implements the Repository interface for the MongoDB database
type Repository struct {
	Conn        *Conn
	IDGenerator func() interface{}
}

// NewRepository returns a MongoDB repository using the given connection
func NewRepository(conn *Conn) *Repository {
	return &Repository{
		Conn:        conn,
		IDGenerator: func() interface{} { return bson.NewObjectId() },
	}
}

// Insert inserts the resource in the corresponding collection
func (r *Repository) Insert(d database.Resource) error {
	colName, err := r.Conn.Map().ColFor(d)
	if err != nil {
		return err
	}
	r.touch(d)
	if r.IDGenerator != nil && d.GetID().(bson.ObjectId).Hex() == "" {
		d.SetID(r.IDGenerator())
	}
	return r.Conn.C(colName).Insert(d)
}

// Update updates the resource in the database
func (r *Repository) Update(d database.Resource) error {
	colName, err := r.Conn.Map().ColFor(d)
	if err != nil {
		return err
	}
	r.touch(d)
	return r.Conn.C(colName).UpdateId(d.GetID(), d)
}

func (r *Repository) touch(d database.Resource) {
	if t, ok := d.(database.Toucher); ok {
		t.Touch()
	}
}

// Delete deletes the resource from the database
func (r *Repository) Delete(d database.Resource) error {
	colName, err := r.Conn.Map().ColFor(d)
	if err != nil {
		return err
	}

	delColName, err := r.Conn.Map().DeletedColFor(d)
	if err != nil {
		return err
	}
	if delColName == "" {
		return r.Conn.C(colName).RemoveId(d.GetID())
	}

	if td, ok := d.(database.SoftDeletable); ok {
		// Logic delete
		td.MarkDeleted()
	}
	err = r.Conn.C(delColName).Insert(d)
	if err != nil {
		return err
	}
	return r.Conn.C(colName).RemoveId(d.GetID())
}

// Get retrieves a resource from the database and stores in the given type
func (r *Repository) Get(ID interface{}, dest database.Resource) error {
	if idstr, ok := ID.(string); ok {
		ID = bson.ObjectIdHex(idstr)
	}
	return r.FindOne(dest, database.NewQ(bson.M{"_id": ID}))
}

// FindOne runs the given query, retrieving a single resource and stores in the given type
func (r *Repository) FindOne(dest database.Resource, query database.Query) error {
	it, err := r.query(dest, query)
	if err != nil {
		return err
	}
	return it.One(dest)
}

// FindMany runs the given query, retrieving all matching resources and stores in the given type slice
func (r *Repository) FindMany(dest database.ResourceList, query database.Query) error {
	it, err := r.query(dest, query)
	if err != nil {
		return err
	}
	return it.All(dest)
}

func (r *Repository) query(dest interface{}, query database.Query) (*mgo.Query, error) {
	colName, err := r.Conn.Map().ColFor(dest)

	if err != nil {
		return nil, err
	}
	it := r.Conn.C(colName).Find(query.Where)
	if len(query.Sort) > 0 {
		it.Sort(query.Sort...)
	}
	// Negative numbers are accepted
	if query.Limit != 0 {
		it.Limit(query.Limit)
	}
	// Negative numbers are accepted
	if query.Skip != 0 {
		it.Skip(query.Skip)
	}
	return it, nil
}

// FetchRelated fetchs resouces related to the given resource
func (r *Repository) FetchRelated(source database.Resource, relations ...string) error {
	for _, relation := range relations {
		rel, err := r.Conn.Map().Relationship(source, relation)
		if err != nil {
			return err
		}
		switch rel.Kind {
		case database.HasOne, database.HasZeroOne:
			ID := idFromField(rel, source)
			if ID == nil {
				continue
			}
			dest := getTargetField(rel, source)
			err = r.Get(ID, dest.(database.Resource))
			if err != nil {
				return err
			}
		case database.HasMany:
			IDs := fieldByName(source, rel.FieldName).([]bson.ObjectId)
			if len(IDs) == 0 {
				continue
			}
			where := bson.M{"_id": bson.M{"$in": IDs}}
			dest := getTargetField(rel, source)
			err = r.FindMany(dest, database.NewQ(where))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func idFromField(rel database.Relationship, res database.Resource) *bson.ObjectId {
	rawID := fieldByName(res, rel.FieldName)
	switch rel.Kind {
	case database.HasOne:
		val := rawID.(bson.ObjectId)
		return &val
	case database.HasZeroOne:
		return rawID.(*bson.ObjectId)
	default:
		return nil
	}
}

// getTargetField returns the struct field for storing related resources
func getTargetField(rel database.Relationship, res database.Resource) interface{} {
	if targeter, ok := res.(database.RelationshipTargeter); ok {
		// The Resource provides a GetRelTarget method that gives us the target field
		return targeter.RelationshipTarget(rel.Name)
	}

	v := reflect.ValueOf(res)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	f := v.FieldByName(rel.TargetField)
	if f.Kind() != reflect.Ptr {
		return f.Addr().Interface()
	}
	if !f.IsNil() {
		return f.Interface()
	}

	zval := reflect.New(f.Type().Elem())
	f.Set(zval)
	return zval.Interface()
}

func fieldByName(i interface{}, field string) interface{} {
	t := reflect.ValueOf(i)
	// If pointer or slice, get it's element
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.FieldByName(field).Interface()
}

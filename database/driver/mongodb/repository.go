package mongodb

import (
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
	if r.IDGenerator != nil {
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
		return nil
	}
	if delColName == "" {
		return r.Conn.C(colName).RemoveId(d.GetID())
	}

	if td, ok := d.(database.SoftDeletable); ok {
		// Logic delete
		td.MarkDeleted()
	}
	return r.Conn.C(delColName).Insert(d)
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

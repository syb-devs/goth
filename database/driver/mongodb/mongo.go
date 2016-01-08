package mongodb

import (
	"encoding/gob"

	"bitbucket.org/syb-devs/goth/database"
	"github.com/syb-devs/gotools/time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func init() {
	gob.Register(bson.ObjectId(""))
}

// Compile-time interface check
var _ database.Connection = (*Conn)(nil)
var _ database.Resource = (*Resource)(nil)

// IsNotFound checks if the given error is a MongoDB not found errorr
func IsNotFound(err error) bool {
	return err == mgo.ErrNotFound
}

// Conn represents a MongoDB connection
type Conn struct {
	*mgo.Database
	resMap *database.ResourceMap
}

// NewConnection creates a MongoDB connection with the given parameters
func NewConnection(ps database.ConnectionParams, resMap *database.ResourceMap) (*Conn, error) {
	c := &Conn{
		resMap: resMap,
	}
	return c, c.Connect(ps)
}

// Connect tries to connect to the MongoDB database
func (c *Conn) Connect(ps database.ConnectionParams) error {
	sess, err := mgo.Dial(ps["url"].(string))
	if err != nil {
		return err
	}
	sess.SetMode(mgo.Monotonic, true)
	c.Database = sess.DB(ps["database"].(string))
	return nil
}

// Copy creates a copy of the connection
func (c *Conn) Copy() database.Connection {
	dbName := c.Database.Name
	sess := c.Database.Session.Copy()

	return &Conn{
		Database: sess.DB(dbName),
		resMap:   c.resMap,
	}
}

// Close ends the current connection with the MongoDB databse
func (c *Conn) Close() error {
	c.Database.Session.Close()
	return nil
}

// Map returs the resource map associated with the connection
func (c *Conn) Map() *database.ResourceMap {
	return c.resMap
}

// Resource implements the Resource interface for MongoDB
type Resource struct {
	ID                bson.ObjectId `json:"id" bson:"_id"`
	time.DeleteStamps `json:",inline" bson:",inline"`
}

// SetID sets the ID for the resource
func (r *Resource) SetID(ID interface{}) error {
	r.ID = ID.(bson.ObjectId)
	return nil
}

// GetID returns the ID for the resoruce
func (r *Resource) GetID() interface{} {
	return r.ID
}

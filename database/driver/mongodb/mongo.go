package mongodb

import (
	"encoding/gob"

	"bitbucket.org/syb-devs/goth/database"
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
	ID                bson.ObjectId `bson:"_id" json:"id"`
	OwnerID           bson.ObjectId `bson:"ownerId,omitempty" json:"ownerId"`
	database.DeleteTS `json:",inline" bson:",inline"`
}

// SetID sets the ID for the resource
func (r *Resource) SetID(ID interface{}) {
	r.ID = mongoID(ID)
}

func mongoID(ID interface{}) bson.ObjectId {
	switch ID := ID.(type) {
	case string:
		return bson.ObjectIdHex(ID)
	case bson.ObjectId:
		return ID
	default:
		panic("invalid type for bson.ObjectId")
	}
}

// GetID returns the ID for the resoruce
func (r *Resource) GetID() interface{} {
	return r.ID
}

// GetIDString returns a string representation of the Resource ID
func (r *Resource) GetIDString() string {
	return r.ID.Hex()
}

// SetOwnerID sets the Owner ID of the resource
func (r *Resource) SetOwnerID(ID interface{}) {
	r.OwnerID = mongoID(ID)
}

// GetOwnerID returns the Owner ID of the resource
func (r *Resource) GetOwnerID() interface{} {
	return r.OwnerID
}

// BelongsTo checks if the resource belongs to (is ownbed by) another resource
func (r *Resource) BelongsTo(ow database.ResourceOwner) bool {
	mongoID, ok := ow.GetID().(bson.ObjectId)
	if !ok {
		return false
	}
	return r.GetOwnerID().(bson.ObjectId) == mongoID
}

package database

// ConnectionParams is a map which stores database connection parameters
type ConnectionParams map[string]interface{}

// Connection is an interface for handling database connections
type Connection interface {
	Connect(ConnectionParams) error
	Close() error
	Copy() Connection
}

// Repository is an interface implemented by database repositories
type Repository interface {
	Insert(d Resource) error
	Update(d Resource) error
	Delete(d Resource) error
	Get(ID interface{}, dest Resource) error
	FindOne(dest Resource, q Query) error
	FindMany(dest ResourceList, q Query) error
}

// Query represents a generic query for a database
type Query struct {
	Where interface{}
	Limit int
	Skip  int
	Sort  []string
}

// NewQuery allocates and returns a new Query object
func NewQuery(where interface{}, limit, skip int, sort ...string) Query {
	if where == nil {
		where = Dict{}
	}
	return Query{
		Where: where,
		Limit: limit,
		Skip:  skip,
		Sort:  sort,
	}
}

// NewQ allocates and returns a new Query object
func NewQ(where interface{}, sort ...string) Query {
	return NewQuery(where, 0, 0, sort...)
}

// Dict is an alias to map[string]interface{}
type Dict map[string]interface{}

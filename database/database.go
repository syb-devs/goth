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

// Dict is an alias to map[string]interface{}
type Dict map[string]interface{}

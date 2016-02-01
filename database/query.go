package database

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

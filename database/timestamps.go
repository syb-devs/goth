package database

import "time"

// TS is a group of timestamps for creation and update times
type TS struct {
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}

// Touch updates the UpdatedAt timestamp to the current time
func (ts *TS) Touch() {
	if ts.CreatedAt.IsZero() {
		ts.CreatedAt = time.Now()
	}
	ts.UpdatedAt = time.Now()
}

// DeleteTS adds a timestamp for storing time of delete
type DeleteTS struct {
	TS        `bson:",inline" json:",inline"`
	DeletedAt time.Time `bson:"deletedAt,omitempty" json:"-"`
}

// MarkDeleted updates the DeletedAt timestamp
func (ts *DeleteTS) MarkDeleted() {
	if ts.DeletedAt.IsZero() {
		ts.DeletedAt = time.Now()
	}
}

package mongodb_test

import (
	"testing"

	"github.com/syb-devs/goth/database/driver/mongodb"
)

func TestResourceBelongsTo(t *testing.T) {
	tests := []struct {
		ownerID         string
		resourceOwnerID string
		expected        bool
	}{
		{"507f191e810c19729de860ea", "507f191e810c19729de860ea", true},
		{"507f191e810c19729de860ea", "507f1f77bcf86cd799439011", false},
	}

	for i, test := range tests {
		u := &mongodb.Resource{}
		u.SetID(test.ownerID)

		r := &mongodb.Resource{}
		r.SetOwnerID(test.resourceOwnerID)
		actual := r.BelongsTo(u)
		if actual != test.expected {
			t.Errorf("test #%d: expecting %v got %v", i+1, test.expected, actual)
		}
	}
}

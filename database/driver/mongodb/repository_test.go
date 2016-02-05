package mongodb

import (
	"testing"

	"bitbucket.org/syb-devs/goth/database"

	"gopkg.in/mgo.v2/bson"
)

func BenchmarkGetTargetFieldReflect(b *testing.B) {
	benchmarkGetTargetFieldInterface(b, &User{})
}

func BenchmarkGetTargetFieldInterface(b *testing.B) {
	benchmarkGetTargetFieldInterface(b, &User2{})
}

func benchmarkGetTargetFieldInterface(b *testing.B, u database.Resource) {
	dbmap := initMap()
	rels := dbmap.Relationships(u)
	if len(rels) == 0 {
		b.Fatalf("no relationships for type %T", u)
	}
	for n := 0; n < b.N; n++ {
		for _, rel := range rels {
			getTargetField(rel, u)
		}
	}
}

func initMap() *database.ResourceMap {
	dbmap := database.NewResourceMap()
	dbmap.RegisterResource(User{}, "users", "")
	dbmap.RegisterResource(User2{}, "users", "")
	dbmap.RegisterResource(Todo{}, "todos", "")
	dbmap.RegisterResource(Profile{}, "profiles", "")

	return dbmap
}

type User struct {
	Resource  `bson:",inline" json:",inline"`
	Name      string          `bson:"name" json:"name"`
	TodoIDs   []bson.ObjectId `bson:"todoIds" json:"-" rel:"todos,Todos"`
	Todos     []*Todo         `bson:"-" json:"todos"`
	ProfileID *bson.ObjectId  `bson:"profileId" json:"-" rel:"profile,Profile"`
	Profile   *Profile        `bson:"-" json:"profile"`
}

type Profile struct {
	Resource `bson:",inline" json:",inline"`
	Twitter  string `bson:"twitter" json:"twitter"`
	Facebook string `bson:"facebook" json:"facebook"`
	Linkedin string `bson:"linkedin" json:"linkedin"`
}

type Todo struct {
	Resource `bson:",inline" json:",inline"`
	Title    string `bson:"title" json:"title"`
	Done     bool   `bson:"done" json:"done"`
}

type User2 User

func (u *User2) RelationshipTarget(rel string) interface{} {
	switch rel {
	case "todos":
		return &u.Todos
	case "profile":
		if u.Profile == nil {
			u.Profile = &Profile{}
		}
		return u.Profile
	default:
		return nil
	}
}

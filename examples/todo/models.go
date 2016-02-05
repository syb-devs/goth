package main

import (
	"bitbucket.org/syb-devs/goth/database/driver/mongodb"

	"gopkg.in/mgo.v2/bson"
)

type User struct {
	mongodb.Resource `bson:",inline" json:",inline"`
	Name             string          `bson:"name" json:"name"`
	TodoIDs          []bson.ObjectId `bson:"todoIds" json:"-" rel:"todos,Todos"`
	Todos            []*Todo         `bson:"-" json:"todos"`
	ProfileID        *bson.ObjectId  `bson:"profileId" json:"-" rel:"profile,Profile"`
	Profile          *Profile        `bson:"-" json:"profile"`
}

type Profile struct {
	mongodb.Resource `bson:",inline" json:",inline"`
	Twitter          string `bson:"twitter" json:"twitter"`
	Facebook         string `bson:"facebook" json:"facebook"`
	Linkedin         string `bson:"linkedin" json:"linkedin"`
}

type Todo struct {
	mongodb.Resource `bson:",inline" json:",inline"`
	Title            string `bson:"title" json:"title"`
	Done             bool   `bson:"done" json:"done"`
}

package models

import "gopkg.in/mgo.v2/bson"

//Shop Structure
type Shop struct {
	ID       bson.ObjectId `json:"_id" bson:"_id"`
	Name     string        `json:"name" bson:"name"`
	Category string        `json:"category" bson:"category"`
	OwnerID  string        `json:"owner" bson:"owner"`
}

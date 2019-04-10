package model

import "gopkg.in/mgo.v2/bson"

// User ...
type User struct {
	ID        bson.ObjectId `bson:"_id" json:"id"`
	Name      string        `bson:"name" json:"name"`
	Interests []string      `bson:"interests" json:"interests"`
}

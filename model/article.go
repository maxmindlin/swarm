package model

import "gopkg.in/mgo.v2/bson"

// Article organizes the different components of an article
type Article struct {
	ID          bson.ObjectId `bson:"_id"`
	Title       string        `bson:"title" json:"title"`
	Description string        `bson:"description" json:"description"`
	URL         string        `bson:"url" json:"url"`
	Keywords    []string      `bson:"keywords" json:"keywords"`
}

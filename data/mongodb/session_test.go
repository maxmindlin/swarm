package mongodb

import (
	"testing"

	"github.com/maxmindlin/swarm/model"

	"gopkg.in/mgo.v2/bson"
)

func TestSession(t *testing.T) {
	testObj := model.Article{
		ID:          bson.NewObjectId(),
		Title:       "Foo Title",
		Description: "This is fake news",
		URL:         "https://foxnews.com",
		Keywords:    []string{"foo", "blah"},
	}

	// Test creating a new session client
	s, err := NewSession()
	if err != nil {
		t.Errorf("Failed to create a new MongoDB session client")
	}
	defer s.Close()

	// Test getting a collection
	_, err = s.GetCollection("swarm_test", "foo")
	if err != nil {
		t.Errorf("Failed to get a test collection from MongoDB instance")
	}

	// Test inserting a document
	err = s.Insert("swarm_test", "foo", testObj)
	if err != nil {
		t.Errorf("Failed to write dumy data to MongoDB instance")
	}

	// Test dropping the dummy database
	err = s.DropDatabase("swarm_test")
	if err != nil {
		t.Errorf("Failed to drop the MongoDB test db")
	}
}

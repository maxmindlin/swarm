package mongodb

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Session ...
type Session struct {
	Client *mongo.Client
}

// NewSession generates a new session object, with a mongo client included
func NewSession() (*Session, error) {
	ctx, close := context.WithTimeout(context.Background(), 10*time.Second)
	defer close()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}
	return &Session{client}, err
}

// GetCollection returns a reference to a mongo collection object
// Collection objects are used for data insertion
func (s *Session) GetCollection(dbName string, collName string) (*mongo.Collection, error) {
	if s.Client != nil {
		return s.Client.Database(dbName).Collection(collName), nil
	}
	return nil, errors.New("Mongo client not provided")
}

// Close disconnects the mongo client from the db instance
func (s *Session) Close() {
	if s.Client != nil {
		s.Client.Disconnect(context.TODO())
	}
}

// DropDatabase drops a db of name dbName from the mongo instance
func (s *Session) DropDatabase(dbName string) error {
	if s.Client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return s.Client.Database(dbName).Drop(ctx)
	}
	return errors.New("No client declared to drop database from")
}

// Insert allows quick single document insertion into MongoDB
func (s *Session) Insert(dbName string, collName string, document interface{}) error {
	if s.Client == nil {
		return errors.New("No client declared to pass data thru")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	c, err := s.GetCollection(dbName, collName)
	if err != nil {
		return err
	}
	_, err = c.InsertOne(ctx, document)
	if err != nil {
		return err
	}
	return nil
}

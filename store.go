package main

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type Store interface {
	SaveUser(*User) error
	RetrieveUser() error
}

var ctx = context.Background()

type MongoStore struct {
	db   *mongo.Database
	coll string
}

func (ms *MongoStore) SaveUser(user *User) error {
	_, err := ms.db.Collection(ms.coll).InsertOne(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

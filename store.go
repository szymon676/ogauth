package main

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Store interface {
	SaveUser(*User) error
	RetrieveUser(email string) (*User, error)
}

var ctx = context.TODO()

type MongoStore struct {
	db   *mongo.Database
	coll string
}

func (ms *MongoStore) SaveUser(user *User) error {
	if _, err := ms.RetrieveUser(user.Email); err == nil {
		return errors.New("user already exists")
	}
	_, err := ms.db.Collection(ms.coll).InsertOne(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (ms *MongoStore) RetrieveUser(email string) (*User, error) {
	filter := bson.M{"email": email}
	result := ms.db.Collection(ms.coll).FindOne(ctx, filter)

	var user *User
	err := result.Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("user not found")
	} else if err != nil {
		return nil, err
	}

	return user, nil
}

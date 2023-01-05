package database

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

// This file is used to wrap MongoDB models 
// it allows to keep all Mongo objects in one place

// MongoCollection represents a mongogb implementation of [model.Collection]
type MongoCollection struct {
	Collection *mongo.Collection
}

func (c MongoCollection) Insert(docs ...interface{}) error{
	return nil
}

func (c MongoCollection) CountDocuments(ctx context.Context, docs ...interface{}) (int64, error){
	return c.CountDocuments(ctx, docs)
}

// GetMongoCollection returns a mongo collection with the [name]
// Wraps a Mongo Collection with our implementation of it
func GetMongoCollection(mC *mongo.Collection) MongoCollection {
	c := MongoCollection{mC}

	return c
}

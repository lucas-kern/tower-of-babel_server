package database

import (
	"context"
	"fmt"
	"log"
	"time"
	
	// "github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Database represents the current Database connection
type Database struct {
	client 	*mongo.Client
	cancel 	CancelFunc
	url     string
}

// DatastoreSession is a session connection to the Database
type DatastoreSession interface { //TODO: not yet used
	Close()
	Copy() DatastoreSession
	DB(name string)
}

// Datastore represents a store for the data (Database session)
type Datastore struct {
	client *mongo.Client
	name    string
}

func Connect(host string) (*Database, error){
	// localhost := "localhost:27017"
	mongohost := "mongodb://" + host

	fmt.Print("Connecting to Host: ", host + "\n")

	client, err := mongo.NewClient(options.Client().ApplyURI(mongohost))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() // TODO ensure that calls can still be made to DB with this deferred. Might need to store it in the struct to pass it around
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &Database{client: client}, nil
}

func (d *Database) Close() {
	if d.client != nil {
		d.client.Disconnect(context.TODO())
	}
	d.client = nil
	log.Println("Database closed")
}


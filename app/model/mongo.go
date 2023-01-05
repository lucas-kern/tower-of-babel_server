package model

// import "time"
import "context"


// Collection represents a row of data
type Collection interface {
	Insert(docs ...interface{}) error
	CountDocuments(ctx context.Context, docs ...interface{}) (int64, error)

	// Update performs an [update] on [selector] in the collection
	// Update(selector interface{}, update interface{}) error

	//Remove will remove the document(s) that matches the [selector]
	// Remove(selector interface{}) error

	//Will create an index key on a document
	// EnsureIndex(key []string, time time.Duration) error
}
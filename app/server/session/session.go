package session

import (
	"log"
	"net/http"

	"fmt"

	"github.com/Don-V/RSB-backend/app/backend/database"
	"github.com/Don-V/RSB-backend/app/backend/operations/transaction"
	"github.com/Don-V/RSB-backend/app/model"
	"github.com/kidstuff/mongostore"
	"github.com/gorilla/sessions"
)

// Store represents a datastore for the session
// It is supposed to be a wrapper for gorilla's [sessions.Store]
type Store interface {
	Get(r *http.Request, name string) (Session, error)
	New(r *http.Request, name string) (Session, error)
	Save(r *http.Request, w http.ResponseWriter, s Session) error
}

// internalStore is the store that is used by the server
// It uses the [internalSession]
// It uses [mongostore.MongoStore] for [MaxAge]
type internalStore struct {
	store sessions.Store
}

func (st internalStore) Get(r *http.Request, name string) (Session, error) {
	sess, err := st.store.Get(r, name)
	sess.Options.HttpOnly = true
	return &internalSession{
		session: sess,
	}, err
}

func (st internalStore) New(r *http.Request, name string) (Session, error) {
	sess, err := st.store.New(r, name)
	return &internalSession{
		session: sess,
	}, err
}

func (st internalStore) Save(r *http.Request, w http.ResponseWriter, s Session) error {
	return st.store.Save(r, w, s.(internalSession).session)
}

func (st internalStore) MaxAge(age int) {
	if mst, ok := st.store.(*mongostore.MongoStore); ok {
		mst.MaxAge(age)
	}
}

// Session represents a session for a request
// It is supposed to be a wrapper for gorilla's [sessions.Session]
type Session interface {
	Name() string
	Save(r *http.Request, w http.ResponseWriter) error
	Store() Store
	Set(key string, value interface{})
	Get(key string) (interface{}, bool)
	Assign(result transaction.Result) error
	IsAuth() bool
	IsNew() bool
	ID() string
	Clear()
	Destroy()
	All() map[interface{}]interface{}
}

// internalStore is the session used by the server
// It wraps gorilla's [sessions.Session]
type internalSession struct {
	session *sessions.Session
}

func (s internalSession) Name() string {
	return s.session.Name()
}

func (s internalSession) Save(r *http.Request, w http.ResponseWriter) error {
	return s.session.Save(r, w)
}

func (s internalSession) Store() Store {
	return internalStore{
		store: s.session.Store(),
	}
}

func (s internalSession) Set(key string, value interface{}) {
	s.session.Values[key] = value
}

func (s internalSession) Get(key string) (res interface{}, ok bool) {
	res, ok = s.session.Values[key]
	return //named parameters
}

func (s internalSession) Destroy() {
	s.session.Options.MaxAge = -1
	s.session.Values = nil //make(map[interface{}]interface{})
}

func (s internalSession) Clear() {
	s.session.Values = make(map[interface{}]interface{})
}

// AuthSetup setups a session to have the needed user data from an authentication transaction
// As of right now, this is an operation
func (s internalSession) Assign(result transaction.Result) error {
	data, ok := result.GetData().(map[string]interface{})
	if !ok {
		err := fmt.Errorf("From Session: Invalid transaction type in SessionAuthSetup")
		log.Println(err)
		return err
	}

	s.Set("userid", data["userid"].(string)) //Need for casting? // should be used for user-session cross refernce
	s.Set("email", data["email"].(string))
	s.Set("username", data["username"].(string))

	return nil
}

// IsAuth returns true if this is a valid session that has been setup
func (s internalSession) IsAuth() bool {
	_, ok := s.Get("userid")
	email, _ := s.Get("email")
	username, _ := s.Get("username")
	return ok && email != nil && username != nil //TODO: needs testing
}

func (s internalSession) IsNew() bool {
	return s.session.IsNew
}

func (s internalSession) ID() string {
	return s.session.ID
}

//TODO: remove
func (s internalSession) All() map[interface{}]interface{} {
	return s.session.Values
}

// NewSessionStore returns a new session store for managing sessions
// TODO: move this to a more specific file/package
func NewSessionStore(keyPairs ...[]byte) Store {
	return internalStore{
		store: sessions.NewCookieStore(keyPairs...),
	}
}

// NewMongoSessionStore returns a new session store for managing sessions
// The store is implemted in a MongoDB database
// TODO: move this to a more specific file/package
func NewMongoSessionStore(c model.Collection, keyPairs ...[]byte) Store {
	expire := 60 * 60 //seconds
	coll, ok := c.(database.MongoCollection)
	if !ok {
		log.Fatalf("Collection passed in is not MongoCollection{}. It was %T", c) //TODO: return err, panic -> recover? vs Fatalf
	}
	store := &internalStore{
		store: mongostore.NewMongoStore(coll.Collection, expire, true, keyPairs...),
	}
	store.MaxAge(expire)
	return store
}
package session

import (
	"log"
	"net/http"

	"github.com/lucas-kern/tower-of-babel_server/app/server/database"
	"github.com/lucas-kern/tower-of-babel_server/app/model"

	"github.com/lucas-kern/mongostore"
	"github.com/gorilla/sessions"
)

// Session is used to create and manage sessions 
// it uses a store to store the user sessions once they are created

// TODO: test that this all works Then create users and test that when a user is created they are given only one session at a time

// Store represents a datastore for the session
// It is a wrapper for gorilla's [sessions.Store] interface
type Store interface {
	Get(r *http.Request, name string) (Session, error)
	New(r *http.Request, name string) (Session, error)
	Save(r *http.Request, w http.ResponseWriter, s Session) error
}

// sessionStore is the store that is used by the server
// It uses the [internalSession]
// It uses [mongostore.MongoStore] for [MaxAge]
type sessionStore struct {
	store sessions.Store
}

func (st sessionStore) Get(r *http.Request, name string) (Session, error) {
	sess, err := st.store.Get(r, name)
	sess.Options.HttpOnly = true
	return &internalSession{
		session: sess,
	}, err
}

func (st sessionStore) New(r *http.Request, name string) (Session, error) {
	sess, err := st.store.New(r, name)
	return &internalSession{
		session: sess,
	}, err
}

func (st sessionStore) Save(r *http.Request, w http.ResponseWriter, s Session) error {
	return st.store.Save(r, w, s.(internalSession).session)
}

func (st sessionStore) MaxAge(age int) {
	if mst, ok := st.store.(*mongostore.MongoStore); ok {
		mst.MaxAge(age)
	}
}

// Session represents a session for a request
// It is a wrapper for gorilla's [sessions.Session]
type Session interface {
	Name() string
	Save(r *http.Request, w http.ResponseWriter) error
	Store() Store
	Set(key string, value interface{})
	Get(key string) (interface{}, bool)
	Values()  map[interface{}]interface{}
	IsAuth() bool
	IsNew() bool
	ID() string
	Clear()
	Destroy()
}

// sessionStore is the session used by the server
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
	return sessionStore{
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
// func (s internalSession) Assign(result transaction.Result) error {
// 	data, ok := result.GetData().(map[string]interface{})
// 	if !ok {
// 		err := fmt.Errorf("From Session: Invalid transaction type in SessionAuthSetup")
// 		log.Println(err)
// 		return err
// 	}

// 	s.Set("userid", data["userid"].(string)) //Need for casting? // should be used for user-session cross refernce
// 	s.Set("email", data["email"].(string))
// 	s.Set("username", data["username"].(string))

// 	return nil
// }

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
func (s internalSession) Values() map[interface{}]interface{} {
	return s.session.Values
}

// NewSessionStore returns a new session store for managing sessions
// TODO: move this to a more specific file/package
func NewSessionStore(keyPairs ...[]byte) Store {
	return sessionStore{
		store: sessions.NewCookieStore(keyPairs...),
	}
}

// NewMongoSessionStore returns a new session store for managing sessions
// The store is implemented in a MongoDB database
// TODO: move this to a more specific file/package
func NewMongoSessionStore(c model.Collection, keyPairs ...[]byte) Store {
	expire := 60 * 60 //seconds
	coll, ok := c.(database.MongoCollection)
	if !ok {
		log.Fatalf("Collection passed in is not MongoCollection{}. It was %T", c) //TODO: return err, panic -> recover? vs Fatalf
	}
	store := &sessionStore{
		//MAJOR ISSUE find out why it won't update to my new mongostore project
		store: mongostore.NewMongoStore(coll.Collection, expire, true, keyPairs...),
	}
	store.MaxAge(expire)
	return store
}
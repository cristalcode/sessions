package sessions

import (
	"net/http"

	"github.com/gorilla/sessions"
)

var store *sessions.CookieStore
var storeName string

//Session saves an user Session in memcached.
type Session struct {
	Value Unique
	token string
}

//Unique is needs GetID function and is used as Value in Session struct.
type Unique interface {
	GetID() string
}

//InitSession creates a new Session Store.
func InitSession(secret, storeN, domain string) {
	store = sessions.NewCookieStore([]byte(secret))
	store.Options.Domain = domain
	storeName = storeN
}

//NewSession returns a new Session
func NewSession(u Unique) Session {
	return Session{
		Value: u,
		token: generateToken(),
	}
}

//Save set values of session.
func (s *Session) Save(r *http.Request, w http.ResponseWriter, key string) error {
	session, err := store.Get(r, storeName)
	if err != nil {
		return err
	}
	session.Values[key] = s.token
	err = setCacheSession(s)
	if err != nil {
		return err
	}
	return session.Save(r, w)
}

func (s *Session) update() error {
	return setCacheSession(s)
}

//GetID returns value ID.
func (s *Session) GetID(r *http.Request, w http.ResponseWriter, key string) (ID string, err error) {
	err = s.Get(r, w, key)
	if err == nil {
		ID = s.Value.GetID()
	}
	return
}

//Get returns information of session.
func (s *Session) Get(r *http.Request, w http.ResponseWriter, key string) error {
	session, err := store.Get(r, storeName)
	if err != nil {
		return err
	}
	var ok bool
	s.token, ok = session.Values[key].(string)
	if !ok {
		return ErrNoSession
	}
	return getCacheSession(s)
}

//Delete removes session.
func (s *Session) Delete(r *http.Request, w http.ResponseWriter, key string) error {
	session, err := store.Get(r, storeName)
	if err != nil {
		return err
	}
	delete(session.Values, key)
	err = session.Save(r, w)
	if err != nil {
		return err
	}
	return deleteCacheSession(s)
}

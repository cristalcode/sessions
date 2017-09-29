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
	req   *http.Request
	rw    http.ResponseWriter
}

//Unique needs GetID function and is used as Value in Session struct.
type Unique interface {
	GetID() string
}

//InitSession creates a new gorilla CookieStore.
func InitSession(secret, storeN, domain string) {
	store = sessions.NewCookieStore([]byte(secret))
	store.Options.Domain = domain
	storeName = storeN
}

//NewSession returns a new Session using request and response for futher functions.
func NewSession(r *http.Request, w http.ResponseWriter) *Session {
	return &Session{
		req: r,
		rw:  w,
	}
}

//Save stores token into CookieStore and memcached.
func (s *Session) Save(u Unique, key string) error {
	s.token = newToken()
	session, err := store.Get(s.req, storeName)
	if err != nil {
		return err
	}
	session.Values[key] = s.token
	s.Value = u
	err = setCacheSession(s)
	if err != nil {
		return err
	}
	return session.Save(s.req, s.rw)
}

//Update saves new session value into memcached storage.
func (s *Session) Update(u Unique, key string) error {
	session, err := store.Get(s.req, storeName)
	if err != nil {
		return err
	}
	var ok bool
	s.token, ok = session.Values[key].(string)
	if !ok {
		return ErrNoSession
	}
	s.Value = u
	return setCacheSession(s)
}

//Get uses CookieStore to get saved token and gets session value from memcached storage.
func (s *Session) Get(key string) error {
	session, err := store.Get(s.req, storeName)
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

//GetID returns value ID.
func (s *Session) GetID(key string) (id string, err error) {
	err = s.Get(key)
	if err != nil {
		id = s.Value.GetID()
	}
	return
}

//Delete removes token of CookieStore and memcached value.
func (s *Session) Delete(key string) error {
	session, err := store.Get(s.req, storeName)
	if err != nil {
		return err
	}
	delete(session.Values, key)
	err = session.Save(s.req, s.rw)
	if err != nil {
		return err
	}
	return deleteCacheSession(s)
}

package sessions

import (
	ex "errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"gitlab.com/sociallabs/quickrest/errors"
)

var store *sessions.CookieStore
var storeName string

//Session saves an user Session in memcached.
type Session struct {
	Value Unique
	Token string `json:"-"`
}

//Unique is needs GetID function and is used as Value in Session struct.
type Unique interface {
	GetID() string
}

//InitSession creates a new Session Store.
func InitSession(secret, storeN string) {
	store = sessions.NewCookieStore([]byte(secret))
	storeName = storeN
}

//NewSession returns a new Session
func NewSession(u Unique, token string) Session {
	return Session{
		Value: u,
		Token: token,
	}
}

//Save set values of session.
func (s *Session) Save(c *gin.Context) errors.Message {
	session, err := store.Get(c.Request, storeName)
	if err != nil {
		return errors.NewMessage(http.StatusInternalServerError, err)
	}
	session.Values["token"] = s.Token
	ex := setCacheSession(s)
	if ex != errors.NoError {
		return ex
	}
	err = session.Save(c.Request, c.Writer)
	if err != nil {
		return errors.NewMessage(http.StatusInternalServerError, err)
	}
	return errors.NoError
}

func (s *Session) update() errors.Message {
	return setCacheSession(s)
}

//GetID returns value ID.
func (s *Session) GetID(c *gin.Context) (ID string, ex errors.Message) {
	ex = s.Get(c)
	if ex == errors.NoError {
		ID = s.Value.GetID()
	}
	return
}

//Get returns information of session.
func (s *Session) Get(c *gin.Context) errors.Message {
	session, err := store.Get(c.Request, storeName)
	if err != nil {
		return errors.NewMessage(http.StatusInternalServerError, err)
	}
	var ok bool
	s.Token, ok = session.Values["token"].(string)
	if !ok {
		return errors.NewMessage(http.StatusUnauthorized, ex.New("no Session"))
	}
	ex := getCacheSession(s)
	if ex != errors.NoError {
		return ex
	}
	err = session.Save(c.Request, c.Writer)
	if err != nil {
		return errors.NewMessage(http.StatusInternalServerError, err)
	}
	return errors.NoError
}

//Delete removes session.
func (s *Session) Delete(c *gin.Context) errors.Message {
	session, err := store.Get(c.Request, storeName)
	if err != nil {
		return errors.NewMessage(http.StatusInternalServerError, err)
	}
	delete(session.Values, "token")
	err = session.Save(c.Request, c.Writer)
	if err != nil {
		return errors.NewMessage(http.StatusInternalServerError, err)
	}
	return deleteCacheSession(s)
}

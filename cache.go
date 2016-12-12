package sessions

import (
	"bytes"
	"encoding/gob"
	ex "errors"
	"net/http"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"gitlab.com/sociallabs/quickrest/errors"
)

var (
	cache *memcache.Client
)

//InitCache connects to memcached host
func InitCache(cacheHost string, uniques ...Unique) {
	cache = memcache.New(cacheHost)
	cache.Timeout = 5 * time.Second
	if cache == nil {
		panic("Nil Memcached client in host: " + cacheHost)
	}
	for _, u := range uniques {
		gob.Register(u)
	}
}

func setCacheSession(s *Session) errors.Message {
	buff := new(bytes.Buffer)
	enc := gob.NewEncoder(buff)
	err := enc.Encode(s)
	if err != nil {
		return errors.NewMessage(http.StatusInternalServerError, err)
	}
	err = cache.Set(&memcache.Item{Key: s.Token, Value: buff.Bytes()})
	if err != nil {
		return errors.NewMessage(http.StatusInternalServerError, err)
	}
	return errors.NoError
}

func getCacheSession(s *Session) errors.Message {
	item, err := cache.Get(s.Token)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			return errors.NewMessage(http.StatusUnauthorized, ex.New("no session"))
		}
		return errors.NewMessage(http.StatusInternalServerError, err)
	}
	buff := bytes.NewBuffer(item.Value)
	dec := gob.NewDecoder(buff)
	err = dec.Decode(s)
	if err != nil {
		return errors.NewMessage(http.StatusInternalServerError, err)
	}
	return errors.NoError
}

func deleteCacheSession(s *Session) errors.Message {
	err := cache.Delete(s.Token)
	if err != nil {
		return errors.NewMessage(http.StatusInternalServerError, err)
	}
	return errors.NoError
}

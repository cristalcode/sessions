package sessions

import (
	"bytes"
	"encoding/gob"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

var (
	cache *memcache.Client
)

//InitCache connects to memcached host
func InitCache(cacheHost string, uniques ...Unique) {
	cache = memcache.New(cacheHost)
	cache.Timeout = 5 * time.Second
	if cache == nil {
		panic("memcached is not running on: " + cacheHost)
	}
	for _, u := range uniques {
		gob.Register(u)
	}
}

func setCacheSession(s *Session) error {
	buff := new(bytes.Buffer)
	enc := gob.NewEncoder(buff)
	err := enc.Encode(s)
	if err != nil {
		return err
	}
	return cache.Set(&memcache.Item{Key: s.token, Value: buff.Bytes()})
}

func getCacheSession(s *Session) error {
	item, err := cache.Get(s.token)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			return ErrNoSession
		}
		return err
	}
	buff := bytes.NewBuffer(item.Value)
	dec := gob.NewDecoder(buff)
	return dec.Decode(s)
}

func deleteCacheSession(s *Session) error {
	return cache.Delete(s.token)
}

package sessions

import "errors"

//ErrNoSession is returned when token is not stored in memcached or gorilla cookie store.
var ErrNoSession = errors.New("token is not stored")

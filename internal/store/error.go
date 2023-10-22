package store

import "errors"

var ErrKeyAlreadyExists = errors.New("key already exists")
var ErrKeyNotFound = errors.New("key not found")

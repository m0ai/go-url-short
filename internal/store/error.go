package store

type ErrKeyAlreadyExists struct {
	Err error
}

type ErrKeyNotFound struct {
	Err error
}

func (e ErrKeyAlreadyExists) Error() string {
	return "key already exists"
}

func (e ErrKeyNotFound) Error() string {
	return "key Not Found"
}

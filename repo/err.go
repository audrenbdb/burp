package repo

import (
	"errors"
	"fmt"
)

var ErrNotFound = errors.New("resource not found")

type Err struct {
	err error
}

func (e Err) Error() string { return e.err.Error() }
func (e Err) Unwrap() error { return e.err }

func Error(msg string) Err {
	return Err{
		err: errors.New(msg),
	}
}

func Errorf(msg string, args ...any) Err {
	return Err{
		err: fmt.Errorf(msg, args...),
	}
}

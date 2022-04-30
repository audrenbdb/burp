package burp

import (
	"errors"
	"fmt"
)

var (
	ErrBeerNameTooLong       = Error("name exceed 15 character")
	ErrBeerNameMissing       = Error("name is missing")
	ErrBeerCreateDateMissing = Error("creation date is missing")
	ErrBeerUpdateDateMissing = Error("update date is missing")

	ErrIDEmpty = Error("id cannot be empty")

	ErrCurrencyNotSupported = Error("currency not supported")
)

type Err struct {
	err error
}

func Error(str string) error {
	return Err{err: errors.New(str)}
}

func Errorf(str string, args ...any) error {
	return Err{err: fmt.Errorf(str, args...)}
}

func (e Err) Error() string { return e.err.Error() }
func (e Err) Unwrap() error { return e.err }

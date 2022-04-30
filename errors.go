package burp

import "fmt"

type ErrNotFound struct {
	ResourceName string
}

type ErrBadRequest struct {
	Hint string
}

type ErrUnauthorized struct {
	Hint string
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("%s not found", e.ResourceName)
}

func (e ErrBadRequest) Error() string {
	return fmt.Sprintf("bad request: %s", e.Hint)
}

func (e ErrUnauthorized) Error() string {
	return fmt.Sprintf("unauthorized: %s", e.Hint)
}

var (
	ErrBeerNameTooShort  = ErrBadRequest{Hint: "beer name should have at least 2 characters"}
	ErrBeerNameTooLong   = ErrBadRequest{Hint: "beer name should have maximum 15 characters"}
	ErrBeerRatingInvalid = ErrBadRequest{Hint: "beer rating should be between 0 (very bad) and 5 (very good)"}
	ErrBeerNotFound      = ErrNotFound{ResourceName: "beer"}
)

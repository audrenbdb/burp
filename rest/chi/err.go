package chi

import (
	"burp"
	"burp/repo"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

type HandlerWithErr func(w http.ResponseWriter, r *http.Request) error

type apiError struct {
	Code         int       `json:"code"`
	ErrorMessage string    `json:"error"`
	Time         time.Time `json:"time"`
}

func (err apiError) Error() string { return err.ErrorMessage }

// Handle centralizes handlers error handling
func Handle(fn HandlerWithErr) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := fn(w, r)
		if err == nil {
			return
		}

		var unmarshalTypeError *json.UnmarshalTypeError
		var parseTimeError *time.ParseError

		now := time.Now()
		apiErr := apiError{Time: now}

		switch {
		case errors.As(err, &unmarshalTypeError):
			apiErr.Code = http.StatusBadRequest
			apiErr.ErrorMessage = fmt.Sprintf("corrupted %s type", unmarshalTypeError.Field)
		case errors.As(err, &parseTimeError):
			apiErr.Code = http.StatusBadRequest
			apiErr.ErrorMessage = fmt.Sprintf("corrupted time value: %s", parseTimeError.Value)
		case errors.As(err, &burp.Err{}):
			apiErr.Code = http.StatusBadRequest
			apiErr.ErrorMessage = err.Error()
		case errors.As(err, &apiErr):
			apiErr.Time = now
		case errors.Is(err, repo.ErrNotFound):
			apiErr.Code = http.StatusNotFound
			apiErr.ErrorMessage = err.Error()
		default:
			log.Printf("Internal error:\n%q\n", err.Error())
			apiErr.Code = http.StatusInternalServerError
			apiErr.ErrorMessage = "internal error"
		}

		jsonB, err := json.Marshal(apiErr)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("unable to produce a valid api error"))
		}

		w.WriteHeader(apiErr.Code)
		w.Write(jsonB)
	}
}

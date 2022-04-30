package burp

import (
	"encoding/json"
	"errors"
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestUnmarshalBeer(t *testing.T) {
	t.Run("A beer with name too short cannot be unmarshalled", func(t *testing.T) {
		beerPayload := []byte(`{"alcoholLevel":5}`)
		var b Beer
		err := json.Unmarshal(beerPayload, &b)
		if !errors.Is(err, ErrBeerNameTooShort) {
			t.Errorf("got error: %v, want: %v", err, ErrBeerNameTooShort)
		}
	})
	t.Run("A beer with name too long cannot be unmarshalled", func(t *testing.T) {
		beerPayload := []byte(`{"name":"This is an example of a very long beer name"}`)
		var b Beer
		err := json.Unmarshal(beerPayload, &b)
		if !errors.Is(err, ErrBeerNameTooLong) {
			t.Errorf("got error: %v, want: %v", err, ErrBeerNameTooLong)
		}
	})
	t.Run("A beer with invalid rating cannot be unmarshalled", func(t *testing.T) {
		beerPayload := []byte(`{"name":"Guinness","rating":-1}`)
		var b Beer
		err := json.Unmarshal(beerPayload, &b)
		if !errors.Is(err, ErrBeerRatingInvalid) {
			t.Errorf("got error: %v, want: %v", err, ErrBeerRatingInvalid)
		}
	})
	t.Run("A beer with valid fields should be unmarshalled properly", func(t *testing.T) {
		beerPayload := []byte(`{"name":"Guinness","rating":4}`)
		var b Beer
		json.Unmarshal(beerPayload, &b)
		if diff := cmp.Diff(b, Beer{Name: "Guinness", Rating: VeryGood}); diff != "" {
			t.Errorf("got/want beer: %s", diff)
		}
	})
}

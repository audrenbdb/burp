package chi

import (
	"burp"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"net/http"
	"time"
)

func DeleteBeer(remover BeerRemover) HandlerWithErr {
	return func(w http.ResponseWriter, r *http.Request) error {
		p := chi.URLParam(r, "id")
		id, err := uuid.Parse(p)
		if err != nil {
			return apiError{
				Code:         http.StatusBadRequest,
				ErrorMessage: fmt.Sprintf("invalid id %q: %s", p, err),
			}
		}

		err = remover.RemoveBeer(r.Context(), burp.ID{UUID: id})
		if err != nil {
			return err
		}

		w.WriteHeader(http.StatusNoContent)

		return nil
	}
}

func GetBeer(selector BeerSelector) HandlerWithErr {
	return func(w http.ResponseWriter, r *http.Request) error {
		p := chi.URLParam(r, "id")
		id, err := uuid.Parse(p)
		if err != nil {
			return apiError{
				Code:         http.StatusBadRequest,
				ErrorMessage: fmt.Sprintf("invalid id %q: %s", p, err),
			}
		}

		beer, err := selector.SelectBeer(r.Context(), burp.ID{UUID: id})
		if err != nil {
			return err
		}

		return json.NewEncoder(w).Encode(&beer)
	}
}

func PutBeer(saver BeerSaver) HandlerWithErr {
	return func(w http.ResponseWriter, r *http.Request) error {
		var beer burp.Beer

		if err := json.NewDecoder(r.Body).Decode(&beer); err != nil {
			return err
		}

		if id := chi.URLParam(r, "id"); id != beer.ID.String() {
			return apiError{
				Code:         http.StatusBadRequest,
				ErrorMessage: fmt.Sprintf("resource ID %q not found in request body", id),
			}
		}

		if err := beer.Validate(); err != nil {
			return err
		}

		if err := saver.SaveBeer(r.Context(), &beer); err != nil {
			return err
		}

		w.WriteHeader(http.StatusAccepted)
		return json.NewEncoder(w).Encode(beer)
	}
}

func PostBeer(saver BeerSaver) HandlerWithErr {
	type fields struct {
		Name  string     `json:"name"`
		Price burp.Price `json:"price"`
	}

	return func(w http.ResponseWriter, r *http.Request) error {
		var beerFields fields
		now := time.Now().UTC()

		err := json.NewDecoder(r.Body).Decode(&beerFields)
		if err != nil {
			return err
		}

		beer := burp.Beer{
			ID:        burp.ID{UUID: uuid.New()},
			CreatedAt: now,
			UpdatedAt: now,

			Name: beerFields.Name,
			Price: burp.Price{
				Currency: beerFields.Price.Currency,
				Amount:   beerFields.Price.Amount,
			},
		}

		err = beer.Validate()
		if err != nil {
			return err
		}

		err = saver.SaveBeer(r.Context(), &beer)
		if err != nil {
			return err
		}

		w.WriteHeader(http.StatusCreated)
		return json.NewEncoder(w).Encode(beer)
	}
}

package httpreq

import (
	"burp"
	"burp/router"
	"context"
	"encoding/json"
	"net/http"
	"path"
	"strconv"
	"strings"
)

//go:generate mockgen -source $GOFILE -destination ../mock/$GOFILE -package mock -mock_names Service=Service

func NewHandler(service Service) http.Handler {
	r := router.New()
	r.Get("/v1/beers", handleGetBeers(service))

	r.Put("/v1/beers/", handlePutBeer(service))
	r.Get("/v1/beers/", handleGetOneBeer(service))
	r.Delete("/v1/beers/", handleDeleteBeer(service))
	return r
}

type Service interface {
	BeerSaver
	BeerDeleter
	BeerGetter
	BeerSearcher
}

type (
	BeerSaver interface {
		SaveBeer(ctx context.Context, beer burp.Beer) error
	}
	BeerDeleter interface {
		DeleteBeer(ctx context.Context, name string) error
	}
	BeerGetter interface {
		GetBeerByName(ctx context.Context, name string) (burp.Beer, error)
	}
	BeerSearcher interface {
		SearchBeers(ctx context.Context, filter burp.BeerFilter) ([]burp.Beer, error)
	}
)

// handler is a custom http.HandlerFunc with an error attached
type handler func(w http.ResponseWriter, r *http.Request) error

func handlePutBeer(saver BeerSaver) handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		var b burp.Beer
		err := json.NewDecoder(r.Body).Decode(&b)
		if err != nil {
			return err
		}
		if p := path.Base(r.URL.Path); p != b.Name {
			return burp.ErrBadRequest{Hint: "beer name differs in URL and JSON body"}
		}
		err = saver.SaveBeer(r.Context(), b)
		if err != nil {
			return err
		}
		w.WriteHeader(http.StatusAccepted)
		return json.NewEncoder(w).Encode(b)
	}
}

func handleDeleteBeer(deleter BeerDeleter) handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		name := path.Base(r.URL.Path)
		err := deleter.DeleteBeer(r.Context(), name)
		if err != nil {
			return err
		}
		w.WriteHeader(http.StatusNoContent)
		return nil
	}
}

func handleGetOneBeer(getter BeerGetter) handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		name := path.Base(r.URL.Path)
		beer, err := getter.GetBeerByName(r.Context(), name)
		if err != nil {
			return err
		}
		return json.NewEncoder(w).Encode(beer)
	}
}

func handleGetBeers(searcher BeerSearcher) handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		f, err := beerFilterFromRequest(r)
		if err != nil {
			return err
		}
		beers, err := searcher.SearchBeers(r.Context(), f)
		if err != nil {
			return err
		}
		return json.NewEncoder(w).Encode(beers)
	}
}

func beerFilterFromRequest(r *http.Request) (filter burp.BeerFilter, err error) {
	query := r.URL.Query()
	if name := query.Get("name"); name != "" {
		filter.NameContains = &name
	}
	if ratingsParam := query.Get("ratings"); ratingsParam != "" {
		filter.Ratings, err = parseRatingsParam(ratingsParam)
		if err != nil {
			return filter, err
		}
	}
	return filter, nil
}

func parseRatingsParam(param string) (*[]burp.Rating, error) {
	var ratings []burp.Rating
	for _, s := range strings.Split(param, ",") {
		rating, err := parseRatingParam(s)
		if err != nil {
			return nil, err
		}
		ratings = append(ratings, rating)
	}
	return &ratings, nil
}

func parseRatingParam(r string) (burp.Rating, error) {
	rating, err := strconv.Atoi(r)
	if err != nil {
		return 0, burp.ErrBeerRatingInvalid
	}
	return burp.Rating(rating), nil
}

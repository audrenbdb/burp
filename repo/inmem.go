package repo

import (
	"burp"
	"context"
	"strings"
)

type inMemRepo struct {
	// map of beer name : burp.Beer
	beers map[string]burp.Beer
}

func NewInMemory() burp.Repo {
	return &inMemRepo{
		beers: map[string]burp.Beer{},
	}
}

func (r *inMemRepo) SaveBeer(ctx context.Context, beer burp.Beer) error {
	r.beers[beer.Name] = beer
	return nil
}

func (r *inMemRepo) GetBeerByName(ctx context.Context, name string) (burp.Beer, error) {
	beer, found := r.beers[name]
	if !found {
		return beer, burp.ErrBeerNotFound
	}
	return beer, nil
}

func (r *inMemRepo) DeleteBeer(ctx context.Context, name string) error {
	delete(r.beers, name)
	return nil
}

func (r *inMemRepo) GetBeers(ctx context.Context, filter burp.BeerFilter) ([]burp.Beer, error) {
	beers := make([]burp.Beer, 0)
	if filter.NameContains == nil {
		empty := ""
		filter.NameContains = &empty
	}
	for name, beer := range r.beers {
		if strings.Contains(strings.ToLower(name), strings.ToLower(*filter.NameContains)) {
			beers = append(beers, beer)
		}
	}
	return beers, nil
}

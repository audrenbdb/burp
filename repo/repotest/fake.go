package repotest

import (
	"burp"
	"burp/repo"
	"context"
)

var FakeRepo = &fakeRepo{
	beers: make(map[burp.ID]*burp.Beer),
}

type fakeRepo struct {
	beers map[burp.ID]*burp.Beer
}

func (f *fakeRepo) SaveBeer(ctx context.Context, beer *burp.Beer) error {
	f.beers[beer.ID] = beer
	return nil
}

func (f *fakeRepo) SelectBeer(ctx context.Context, id burp.ID) (*burp.Beer, error) {
	beer, ok := f.beers[id]
	if !ok {
		return nil, repo.ErrNotFound
	}
	return beer, nil
}

func (f *fakeRepo) RemoveBeer(ctx context.Context, id burp.ID) error {
	delete(f.beers, id)
	return nil
}

package repotest

import (
	"burp"
	"burp/repo"
	"context"
)

type Repo struct {
	burp.BeerSaver
	burp.BeerSelector
	burp.BeerRemover
}

func (r Repo) SaveBeer(ctx context.Context, beer *burp.Beer) error {
	if r.BeerSaver != nil {
		return r.BeerSaver.SaveBeer(ctx, beer)
	}
	return repo.Errorf("SaveBeer(ctx, %+v) is unimplemented", beer)
}

func (r Repo) SelectBeer(ctx context.Context, id burp.ID) (*burp.Beer, error) {
	if r.BeerSelector != nil {
		return r.BeerSelector.SelectBeer(ctx, id)
	}
	return nil, repo.Errorf("SelectBeer(ctx, %+v) is unimplemented", id)
}

func (r Repo) RemoveBeer(ctx context.Context, id burp.ID) error {
	if r.BeerRemover != nil {
		return r.BeerRemover.RemoveBeer(ctx, id)
	}
	return repo.Errorf("RemoveBeer(ctx, %+v) is unimplemented", id)
}

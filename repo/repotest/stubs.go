package repotest

import (
	"burp"
	"burp/burptest"
	"burp/repo"
	"context"
)

type (
	beerSaverStub    struct{ Err error }
	beerRemoverStub  struct{ Err error }
	beerSelectorStub struct {
		Beer *burp.Beer
		Err  error
	}
)

var BeerSaverErrStub = beerSaverStub{Err: repo.Error(burptest.RandString(20))}
var BeerRemoverErrStub = beerRemoverStub{Err: repo.Error(burptest.RandString(20))}
var BeerSelectorNotFoundStub = beerSelectorStub{Err: repo.ErrNotFound}
var BeerSelectorStub = beerSelectorStub{Beer: burptest.RandBeer()}

func (s beerSaverStub) SaveBeer(ctx context.Context, b *burp.Beer) error   { return s.Err }
func (s beerRemoverStub) RemoveBeer(ctx context.Context, id burp.ID) error { return s.Err }
func (b beerSelectorStub) SelectBeer(ctx context.Context, id burp.ID) (*burp.Beer, error) {
	return b.Beer, b.Err
}

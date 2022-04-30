package repotest

import (
	"burp"
	"context"
)

type (
	BeerSaverSpy    struct{ BeerSaved *burp.Beer }
	BeerRemoverSpy  struct{ RemovedID burp.ID }
	BeerSelectorSpy struct {
		SelectedID burp.ID
		Beer       *burp.Beer
	}
)

func (s *BeerSaverSpy) SaveBeer(ctx context.Context, b *burp.Beer) error {
	s.BeerSaved = b
	return nil
}

func (b *BeerRemoverSpy) RemoveBeer(ctx context.Context, id burp.ID) error {
	b.RemovedID = id
	return nil
}

func (b *BeerSelectorSpy) SelectBeer(ctx context.Context, id burp.ID) (*burp.Beer, error) {
	b.SelectedID = id
	return b.Beer, nil
}

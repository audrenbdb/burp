package burp

import (
	"context"
	"fmt"
)

type BeerRepo interface {
	BeerSaver
	BeerSelector
	BeerRemover
}

type BeerSaver interface {
	SaveBeer(ctx context.Context, beer *Beer) error
}

type BeerSelector interface {
	SelectBeer(ctx context.Context, id ID) (*Beer, error)
}

type BeerRemover interface {
	RemoveBeer(ctx context.Context, id ID) error
}

type Brewer struct {
	BeerRepo BeerRepo
}

func (b *Brewer) SaveBeer(ctx context.Context, beer *Beer) error {
	if err := b.BeerRepo.SaveBeer(ctx, beer); err != nil {
		return fmt.Errorf("unable to save beer %+v: %w", beer, err)
	}

	return nil
}

func (b *Brewer) RemoveBeer(ctx context.Context, id ID) error {
	if err := b.BeerRepo.RemoveBeer(ctx, id); err != nil {
		return fmt.Errorf("unable to remove beer %+v: %w", id, err)
	}

	return nil
}

func (b *Brewer) SelectBeer(ctx context.Context, id ID) (*Beer, error) {
	beer, err := b.BeerRepo.SelectBeer(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("unable to select beer with id %q: %w", id, err)
	}

	return beer, nil
}

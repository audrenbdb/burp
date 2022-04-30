package burp

import "context"

type Repo interface {
	SaveBeer(ctx context.Context, beer Beer) error
	GetBeerByName(ctx context.Context, name string) (Beer, error)
	GetBeers(ctx context.Context, filter BeerFilter) ([]Beer, error)
	DeleteBeer(ctx context.Context, name string) error
}

type service struct {
	repo Repo
}

func NewService(r Repo) *service {
	return &service{repo: r}
}

func (s *service) SaveBeer(ctx context.Context, beer Beer) error {
	return s.repo.SaveBeer(ctx, beer)
}

func (s *service) DeleteBeer(ctx context.Context, name string) error {
	return s.repo.DeleteBeer(ctx, name)
}

func (s *service) GetBeerByName(ctx context.Context, name string) (Beer, error) {
	return s.repo.GetBeerByName(ctx, name)
}

func (s *service) SearchBeers(ctx context.Context, filter BeerFilter) ([]Beer, error) {
	return s.repo.GetBeers(ctx, filter)
}

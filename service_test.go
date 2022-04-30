package burp_test

import (
	"burp"
	"burp/repo"
	"context"
	"errors"
	"github.com/google/go-cmp/cmp"
	"golang.org/x/exp/slices"
	"testing"
)

func TestSaveBeer(t *testing.T) {
	repo := repo.NewInMemory()
	ctx := context.Background()
	service := burp.NewService(repo)

	t.Run("A beer saved should be found with its name", func(t *testing.T) {
		guinness := burp.Beer{Name: "Guinness"}
		service.SaveBeer(ctx, guinness)
		beer, _ := repo.GetBeerByName(ctx, guinness.Name)
		if diff := cmp.Diff(beer, guinness); diff != "" {
			t.Errorf("got/want beer: %s", diff)
		}
	})
}

func TestDeleteBeer(t *testing.T) {
	repo := repo.NewInMemory()
	ctx := context.Background()
	service := burp.NewService(repo)

	t.Run("A beer deleted should not be found afterward", func(t *testing.T) {
		guinness := burp.Beer{Name: "Guinness"}
		repo.SaveBeer(ctx, guinness)
		service.DeleteBeer(ctx, guinness.Name)
		_, err := repo.GetBeerByName(ctx, guinness.Name)
		if !errors.Is(err, burp.ErrBeerNotFound) {
			t.Errorf("got err: %v, wanted: %v", err, burp.ErrBeerNotFound)
		}
	})
}

func TestGetBeerByName(t *testing.T) {
	ctx := context.Background()

	heineken := burp.Beer{Name: "Heineken", Rating: burp.VeryBad}
	repo := repo.NewInMemory()
	repo.SaveBeer(ctx, heineken)
	s := burp.NewService(repo)

	b, _ := s.GetBeerByName(ctx, heineken.Name)
	if diff := cmp.Diff(heineken, b); diff != "" {
		t.Errorf("got/want beer: %s", diff)
	}
}

func TestSearchBeers(t *testing.T) {
	ctx := context.Background()

	heineken := burp.Beer{Name: "Heineken", Rating: burp.VeryBad}
	guinness := burp.Beer{Name: "Guinness", Rating: burp.Average}

	repo := repo.NewInMemory()
	repo.SaveBeer(ctx, heineken)
	repo.SaveBeer(ctx, guinness)

	s := burp.NewService(repo)

	search := "heinek"
	beers, _ := s.SearchBeers(ctx, burp.BeerFilter{
		NameContains: &search,
	})
	if !slices.Contains(beers, heineken) {
		t.Errorf("beers should contain heineken, got: %v", beers)
	}
}

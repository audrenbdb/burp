package burp_test

import (
	"burp"
	"burp/burptest"
	"burp/repo/repotest"
	"context"
	"errors"
	"testing"
)

func TestSaveBeer(t *testing.T) {
	beer := burptest.RandBeer()
	spy := &repotest.BeerSaverSpy{}
	repo := repotest.Repo{BeerSaver: spy}
	brewer := &burp.Brewer{BeerRepo: repo}

	err := brewer.SaveBeer(context.Background(), beer)
	if err != nil {
		t.Errorf("SaveBeer(ctx, %+v) returned unexpected error:\ngot %v want nil", beer, err)
	}

	if beer != spy.BeerSaved {
		t.Errorf("SaveBeer(ctx, %+v) did not save beer properly in repo:\ngot %+v", beer, spy.BeerSaved)
	}
}

func TestSaveBeerOnRepoFailure(t *testing.T) {
	beer := burptest.RandBeer()
	stub := repotest.BeerSaverErrStub
	repo := repotest.Repo{BeerSaver: stub}
	brewer := &burp.Brewer{BeerRepo: repo}

	err := brewer.SaveBeer(context.Background(), beer)
	if !errors.Is(err, stub.Err) {
		t.Errorf("SaveBeer(ctx, %+v) returned unexpected error:\ngot %v want %v", beer, err, stub.Err)
	}
}

func TestRemoveBeerOnRepoFailure(t *testing.T) {
	beer := burptest.RandBeer()
	stub := repotest.BeerRemoverErrStub
	repo := repotest.Repo{BeerRemover: stub}
	brewer := &burp.Brewer{BeerRepo: repo}

	err := brewer.RemoveBeer(context.Background(), beer.ID)
	if !errors.Is(err, stub.Err) {
		t.Errorf("RemoveBeer(ctx, %q) returned unexpected error:\ngot %v want %v", beer.ID, err, stub.Err)
	}
}

func TestRemoveBeer(t *testing.T) {
	beer := burptest.RandBeer()
	spy := &repotest.BeerRemoverSpy{}
	repo := repotest.Repo{BeerRemover: spy}
	brewer := &burp.Brewer{
		BeerRepo: repo,
	}

	err := brewer.RemoveBeer(context.Background(), beer.ID)
	if err != nil {
		t.Errorf("RemoveBeer(ctx, %+v) returned unexpected error:\ngot %v want nil", beer.ID, err)
	}

	if spy.RemovedID != beer.ID {
		t.Errorf("RemoveBeer(ctx, %+v) has not removed beer from repository", beer.ID)
	}
}

func TestSelectBeerThatDoesNotExist(t *testing.T) {
	beer := burptest.RandBeer()
	stub := repotest.BeerSelectorNotFoundStub
	repo := repotest.Repo{BeerSelector: stub}
	brewer := &burp.Brewer{BeerRepo: repo}

	_, err := brewer.SelectBeer(context.Background(), beer.ID)
	if !errors.Is(err, stub.Err) {
		t.Errorf("SelectBeer(ctx, %q) returned unexpected error:\ngot %v want %v", beer.ID, err, stub.Err)
	}
}

func TestSelectBeer(t *testing.T) {
	stub := repotest.BeerSelectorStub
	repo := repotest.Repo{BeerSelector: stub}
	brewer := burp.Brewer{BeerRepo: repo}

	got, err := brewer.SelectBeer(context.Background(), stub.Beer.ID)
	if err != nil {
		t.Errorf("SelectBeer(ctx, %+v) returned unexpected error:\ngot %v want nil", stub.Beer.ID, err)
	}

	if got != stub.Beer {
		t.Errorf("SelectBeer(ctx, %+v) returned unexpected beer:\ngot %+v want %+v", stub.Beer.ID, got, stub.Beer)
	}
}

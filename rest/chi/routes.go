package chi

import (
	"burp"
	"context"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type App interface {
	BeerSaver
	BeerRemover
	BeerSelector
}

type BeerSaver interface {
	SaveBeer(ctx context.Context, beer *burp.Beer) error
}

type BeerSelector interface {
	SelectBeer(ctx context.Context, id burp.ID) (*burp.Beer, error)
}

type BeerRemover interface {
	RemoveBeer(ctx context.Context, id burp.ID) error
}

func Handler(app App) http.Handler {
	r := chi.NewRouter()

	r.Post("/api/v1/beers", Handle(PostBeer(app)))
	r.Put("/api/v1/beers/{id}", Handle(PutBeer(app)))
	r.Get("/api/v1/beers/{id}", Handle(GetBeer(app)))
	r.Delete("/api/v1/beers/{id}", Handle(DeleteBeer(app)))

	return r
}

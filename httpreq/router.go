package httpreq

import (
	"burp"
	"context"
	"net/http"
	"sort"
)

//go:generate mockgen -source $GOFILE -destination ../mock/$GOFILE -package mock -mock_names Service=Service

func NewHandler(service Service) http.Handler {
	r := router{}
	r.Get("/v1/beers", handleGetBeers(service))

	r.Put("/v1/beers/", handlePutBeer(service))
	r.Get("/v1/beers/", handleGetOneBeer(service))
	r.Delete("/v1/beers/", handleDeleteBeer(service))

	return r.Mux()
}

type Service interface {
	BeerSaver
	BeerDeleter
	BeerGetter
	BeerSearcher
}

type (
	// custom handler with error attached to centralize error handling
	handler func(w http.ResponseWriter, r *http.Request) error

	BeerSaver interface {
		SaveBeer(ctx context.Context, beer burp.Beer) error
	}
	BeerDeleter interface {
		DeleteBeer(ctx context.Context, name string) error
	}
	BeerGetter interface {
		GetBeerByName(ctx context.Context, name string) (burp.Beer, error)
	}
	BeerSearcher interface {
		SearchBeers(ctx context.Context, filter burp.BeerFilter) ([]burp.Beer, error)
	}
)

type router struct {
	endpoints []endpoint
}

type endpoint struct {
	pattern  string
	handlers handlers
}

// map of method:handler
type handlers map[string]handler

func (r *router) Mux() *http.ServeMux {
	mux := http.NewServeMux()
	endpoints := r.endpoints
	sort.Slice(endpoints, func(i, j int) bool {
		return len(endpoints[i].pattern) > len(endpoints[j].pattern)
	})
	for _, ep := range endpoints {
		mux.Handle(ep.pattern, ep.handlers)
	}
	return mux
}

func (r *router) Get(pattern string, handler handler) {
	r.setEndpoint(http.MethodGet, pattern, handler)
}

func (r *router) Post(pattern string, handler handler) {
	r.setEndpoint(http.MethodPost, pattern, handler)
}

func (r *router) Put(pattern string, handler handler) {
	r.setEndpoint(http.MethodPut, pattern, handler)
}

func (r *router) Delete(pattern string, handler handler) {
	r.setEndpoint(http.MethodDelete, pattern, handler)
}

func (r *router) setEndpoint(method, pattern string, h handler) {
	if r.endpoints == nil {
		r.endpoints = []endpoint{}
	}
	ep, ok := r.matchingEndpoint(pattern)
	if !ok {
		ep = endpoint{pattern: pattern, handlers: map[string]handler{}}
		r.endpoints = append(r.endpoints, ep)
	}
	ep.handlers[method] = h
}

func (r router) matchingEndpoint(pattern string) (endpoint, bool) {
	for _, ep := range r.endpoints {
		if ep.pattern == pattern {
			return ep, true
		}
	}
	return endpoint{}, false
}

func (h handlers) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	next, ok := h[r.Method]
	if !ok {
		http.Error(w, "method not supported", http.StatusMethodNotAllowed)
		return
	}
	if err := next(w, r); err != nil {
		handleError(err, w)
	}
}

func handleError(err error, w http.ResponseWriter) {
	switch err.(type) {
	case burp.ErrUnauthorized:
		http.Error(w, err.Error(), http.StatusUnauthorized)
	case burp.ErrBadRequest:
		http.Error(w, err.Error(), http.StatusBadRequest)
	case burp.ErrNotFound:
		http.Error(w, err.Error(), http.StatusNotFound)
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

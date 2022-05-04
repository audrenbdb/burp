package router

import (
	"burp"
	"fmt"
	"log"
	"net/http"
	"sync"
)

type router struct {
	muxMutex  sync.Mutex
	mux       *http.ServeMux
	endpoints []endpoint
}

// New creates a new router based on http.ServeMux
func New() *router {
	return &router{endpoints: []endpoint{}}
}

type endpoint struct {
	pattern       string
	methodHandler handlers
}

// map of method:handler
type handlers map[string]handler

// custom handler with error attached to centralize error handling
type handler = func(w http.ResponseWriter, r *http.Request) error

func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if r.mux == nil {
		r.initMux()
	}
	r.mux.ServeHTTP(w, req)
}

func (r *router) initMux() {
	r.muxMutex.Lock()
	defer r.muxMutex.Unlock()
	if r.mux != nil {
		return
	}
	r.mux = http.NewServeMux()
	endpoints := r.endpoints
	for _, ep := range endpoints {
		r.mux.Handle(ep.pattern, ep.methodHandler)
	}
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
	ep := r.findEndpoint(pattern)
	if ep == nil {
		newEndpoint := endpoint{
			pattern:       pattern,
			methodHandler: map[string]handler{},
		}
		r.endpoints = append(r.endpoints, newEndpoint)
		ep = &newEndpoint
	}
	_, alreadyRegistered := ep.methodHandler[method]
	if alreadyRegistered {
		m := fmt.Sprintf("method %s already registered with pattern %s", method, pattern)
		log.Fatal(m)
	}
	ep.methodHandler[method] = h
}

func (r router) findEndpoint(pattern string) *endpoint {
	for _, ep := range r.endpoints {
		if ep.pattern == pattern {
			return &ep
		}
	}
	return nil
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

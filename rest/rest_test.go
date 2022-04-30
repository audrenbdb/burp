package rest_test

import (
	"burp"
	"burp/repo/repotest"
	"burp/rest/chi"
	"context"
	"net/http"
	"testing"
	"time"
)

const addr = "localhost:8080"

var (
	ctx        = context.Background()
	repository = repotest.FakeRepo
	client     = http.DefaultClient
)

func TestMain(m *testing.M) {
	testContext, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	ctx = testContext

	// list of all routers/handlers to e2e test against
	handlers := []http.Handler{
		chi.Handler(&burp.Brewer{
			BeerRepo: repository,
		}),
	}

	for _, handler := range handlers {
		server := &http.Server{
			Addr:    addr,
			Handler: handler,
		}

		go server.ListenAndServe()

		m.Run()

		server.Shutdown(ctx)
	}
}

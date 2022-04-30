package rest_test

import (
	"burp"
	"burp/burptest"
	"burp/repo"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestDeleteBeer(t *testing.T) {
	beer := burptest.RandBeer()
	endpoint := fmt.Sprintf("http://%s/api/v1/beers/%s", addr, beer.ID.String())

	repository.SaveBeer(ctx, beer)

	response := sendReq(t, http.MethodDelete, endpoint, http.NoBody)

	if response.status != http.StatusNoContent {
		t.Errorf("DELETE beer at endpoint %q returned status %d, want %d",
			endpoint,
			response.status,
			http.StatusNoContent,
		)
	}

	_, err := repository.SelectBeer(ctx, beer.ID)
	if !errors.Is(err, repo.ErrNotFound) {
		t.Errorf("Selecting beer from repository after its delete request returned error %s, want %s", err, repo.ErrNotFound)
	}
}

func TestDeleteBeerWithInvalidID(t *testing.T) {
	id := burptest.RandString(10)
	endpoint := fmt.Sprintf("http://%s/api/v1/beers/%s", addr, id)

	response := sendReq(t, http.MethodDelete, endpoint, http.NoBody)

	if response.status != http.StatusBadRequest {
		t.Errorf("DELETE beer at endpoint %q returned status %d, want %d",
			endpoint,
			response.status,
			http.StatusBadRequest,
		)
	}

	want := fmt.Sprintf("invalid id \\\"%s\\\"", id)
	if !strings.Contains(string(response.body), want) {
		t.Errorf(
			"DELETE beer at endpoint %q\nreturned body: %s\nwant body: %s",
			endpoint,
			string(response.body),
			want,
		)
	}
}

func TestGetBeer(t *testing.T) {
	beer := burptest.RandBeer()
	endpoint := fmt.Sprintf("http://%s/api/v1/beers/%s", addr, beer.ID.String())

	repository.SaveBeer(ctx, beer)

	response := sendReq(t, http.MethodGet, endpoint, http.NoBody)

	if response.status != http.StatusOK {
		t.Errorf("GET beer at endpoint %q returned status %d, want %d",
			endpoint,
			response.status,
			http.StatusOK,
		)
	}

	var got burp.Beer
	json.Unmarshal(response.body, &got)

	if diff := cmp.Diff(beer, &got); diff != "" {
		t.Errorf("RandBeer in response body should match from one saved in repository, (-want/+got):\n%s", diff)
	}
}

func TestGetBeerWithInvalidID(t *testing.T) {
	id := burptest.RandString(10)
	endpoint := fmt.Sprintf("http://%s/api/v1/beers/%s", addr, id)

	response := sendReq(t, http.MethodGet, endpoint, http.NoBody)

	if response.status != http.StatusBadRequest {
		t.Errorf("GET beer at endpoint %q returned status %d, want %d",
			endpoint,
			response.status,
			http.StatusBadRequest,
		)
	}

	want := fmt.Sprintf("invalid id \\\"%s\\\"", id)
	if !strings.Contains(string(response.body), want) {
		t.Errorf(
			"Get beer at endpoint %q\nreturned body: %s\nwant body: %s",
			endpoint,
			string(response.body),
			want,
		)
	}
}

func TestGetBeerNotFound(t *testing.T) {
	beer := burptest.RandBeer()
	endpoint := fmt.Sprintf("http://%s/api/v1/beers/%s", addr, beer.ID)

	response := sendReq(t, http.MethodGet, endpoint, http.NoBody)

	if response.status != http.StatusNotFound {
		t.Errorf("GET beer at endpoint %q returned status %d, want %d",
			endpoint,
			response.status,
			http.StatusNotFound,
		)
	}

	want := "resource not found"
	if !strings.Contains(string(response.body), want) {
		t.Errorf(
			"Get beer at endpoint %q\nreturned body: %s\nwant body: %s",
			endpoint,
			string(response.body),
			want,
		)
	}
}

func TestPutBeer(t *testing.T) {
	beer := burptest.RandBeer()
	endpoint := fmt.Sprintf("http://%s/api/v1/beers/%s", addr, beer.ID.String())
	jsonB, err := json.Marshal(beer)
	if err != nil {
		t.Fatalf("Marshalling beer %+v returned unexpected error: %s", beer, jsonB)
	}

	response := sendReq(t, http.MethodPut, endpoint, bytes.NewReader(jsonB))

	if response.status != http.StatusAccepted {
		t.Errorf("PUT beer json %s at endpoint %q returned status %d, want %d",
			string(jsonB),
			endpoint,
			response.status,
			http.StatusAccepted,
		)
	}

	got, _ := repository.SelectBeer(ctx, beer.ID)
	if diff := cmp.Diff(beer, got); diff != "" {
		t.Errorf("RandBeer found in repository with id %q should match from one sent in PUT request, (-want/+got):\n%s", beer.ID, diff)
	}
}

func TestPutBeerWithMismatchingID(t *testing.T) {
	beer := burptest.RandBeer()
	urlID := uuid.NewString()
	endpoint := fmt.Sprintf("http://%s/api/v1/beers/%s", addr, urlID)
	jsonB, err := json.Marshal(beer)
	if err != nil {
		t.Fatalf("Marshalling beer %+v returned unexpected error: %s", beer, jsonB)
	}

	response := sendReq(t, http.MethodPut, endpoint, bytes.NewReader(jsonB))

	if response.status != http.StatusBadRequest {
		t.Errorf(
			"PUT beer json %s at endpoint %q returned\nstatus: %d, want: %d, \nbody: %s",
			string(jsonB),
			endpoint,
			response.status,
			http.StatusBadRequest,
			string(response.body),
		)
	}

	want := fmt.Sprintf("resource ID \\\"%s\\\" not found in request body", urlID)
	if !strings.Contains(string(response.body), want) {
		t.Errorf(
			"PUT beer json %s\nat endpoint %q\nreturned body: %s\nwant body: %s",
			string(jsonB),
			endpoint,
			string(response.body),
			want,
		)
	}
}

func TestPutBeerWithInvalidFields(t *testing.T) {
	beer := burptest.RandBeer()
	endpoint := fmt.Sprintf("http://%s/api/v1/beers/%s", addr, beer.ID.String())

	tests := []struct {
		name string

		key   string
		value any

		want string
	}{
		{
			name: "IDCorrupted",

			key:   "id",
			value: -1,

			want: "corrupted id type",
		},
		{
			name: "NameTooLong",

			key:   "name",
			value: burptest.RandString(16),

			want: burp.ErrBeerNameTooLong.Error(),
		},
		{
			name:  "CreatedAtEmpty",
			key:   "createdAt",
			value: nil,

			want: burp.ErrBeerCreateDateMissing.Error(),
		},
		{
			name:  "CreatedAtCorrupted",
			key:   "createdAt",
			value: -1,

			want: "corrupted time value: -1",
		},
		{
			name:  "UpdatedAtEmpty",
			key:   "updatedAt",
			value: nil,

			want: burp.ErrBeerUpdateDateMissing.Error(),
		},
		{
			name:  "UpdatedAtCorrupted",
			key:   "updatedAt",
			value: -1,

			want: "corrupted time value: -1",
		},
		{
			name: "NameEmpty",

			key:   "name",
			value: "",

			want: burp.ErrBeerNameMissing.Error(),
		},
		{
			name: "NameCorrupted",

			key:   "name",
			value: 0,

			want: "corrupted name type",
		},
		{
			name: "PriceCurrencyInvalid",

			key: "price",
			value: map[string]any{
				"currency": "Pesos",
				"amount":   beer.Price.Amount,
			},

			want: burp.ErrCurrencyNotSupported.Error(),
		},
		{
			name: "PriceCurrencyCorrupted",

			key: "price",
			value: map[string]any{
				"currency": 1,
				"amount":   beer.Price.Amount,
			},

			want: "corrupted price.currency type",
		},
		{
			name: "PriceAmountCorrupted",

			key: "price",
			value: map[string]any{
				"currency": beer.Price.Currency,
				"amount":   -1,
			},

			want: "corrupted price.amount type",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fields := map[string]any{
				"id":        beer.ID,
				"createdAt": beer.CreatedAt,
				"updatedAt": beer.UpdatedAt,
				"name":      beer.Name,
				"price": map[string]any{
					"currency": beer.Price.Currency,
					"amount":   beer.Price.Amount,
				},
			}

			fields[test.key] = test.value

			jsonB, err := json.Marshal(fields)
			if err != nil {
				t.Fatalf("Marshalling fields %v returned unexpected error: %s", fields, jsonB)
			}

			response := sendReq(t, http.MethodPut, endpoint, bytes.NewReader(jsonB))

			if response.status != http.StatusBadRequest {
				t.Errorf(
					"PUT beer json %s at endpoint %q returned\nstatus: %d, want: %d, \nbody: %s",
					string(jsonB),
					endpoint,
					response.status,
					http.StatusBadRequest,
					string(response.body),
				)
			}

			if !strings.Contains(string(response.body), test.want) {
				t.Errorf(
					"PUT beer json %s\nat endpoint %q\nreturned body: %s\nwant body: %s",
					string(jsonB),
					endpoint,
					string(response.body),
					test.want,
				)
			}
		})
	}
}

func TestPostBeer(t *testing.T) {
	beer := burptest.RandBeer()
	endpoint := "http://" + addr + "/api/v1/beers"
	fields := map[string]any{
		"name": beer.Name,
		"price": map[string]any{
			"currency": beer.Price.Currency,
			"amount":   beer.Price.Amount,
		},
	}

	jsonB, err := json.Marshal(fields)
	if err != nil {
		t.Fatalf("Marshalling fields %v returned unexpected error: %s", fields, jsonB)
	}

	response := sendReq(t, http.MethodPost, endpoint, bytes.NewReader(jsonB))

	if response.status != http.StatusCreated {
		t.Errorf("POST beer json %s at endpoint %q returned status %d, want %d",
			string(jsonB),
			endpoint,
			response.status,
			http.StatusCreated,
		)
	}

	err = json.Unmarshal(response.body, beer)
	if err != nil {
		t.Fatalf("Unmarshalling response body %s into a burp.RandBeer returned error %s", string(response.body), err)
	}

	got, _ := repository.SelectBeer(ctx, beer.ID)
	if diff := cmp.Diff(beer, got); diff != "" {
		t.Errorf("RandBeer found in repository with id %q should match from one received in POST response body, (-want/+got):\n%s", beer.ID, diff)
	}
}

func TestPostBeerWithInvalidFields(t *testing.T) {
	beer := burptest.RandBeer()
	endpoint := "http://" + addr + "/api/v1/beers"

	tests := []struct {
		name string

		key   string
		value any

		want string
	}{
		{
			name: "NameTooLong",

			key:   "name",
			value: burptest.RandString(16),

			want: burp.ErrBeerNameTooLong.Error(),
		},
		{
			name: "NameEmpty",

			key:   "name",
			value: "",

			want: burp.ErrBeerNameMissing.Error(),
		},
		{
			name: "NameCorrupted",

			key:   "name",
			value: 0,

			want: "corrupted name type",
		},
		{
			name: "PriceCurrencyInvalid",

			key: "price",
			value: map[string]any{
				"currency": "Pesos",
				"amount":   beer.Price.Amount,
			},

			want: burp.ErrCurrencyNotSupported.Error(),
		},
		{
			name: "PriceCurrencyCorrupted",

			key: "price",
			value: map[string]any{
				"currency": 1,
				"amount":   beer.Price.Amount,
			},

			want: "corrupted price.currency type",
		},
		{
			name: "PriceAmountCorrupted",

			key: "price",
			value: map[string]any{
				"currency": beer.Price.Currency,
				"amount":   -1,
			},

			want: "corrupted price.amount type",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fields := map[string]any{
				"name": beer.Name,
				"price": map[string]any{
					"currency": beer.Price.Currency,
					"amount":   beer.Price.Amount,
				},
			}

			fields[test.key] = test.value

			jsonB, err := json.Marshal(fields)
			if err != nil {
				t.Fatalf("Marshalling fields %v returned unexpected error: %s", fields, jsonB)
			}

			response := sendReq(t, http.MethodPost, endpoint, bytes.NewReader(jsonB))

			if response.status != http.StatusBadRequest {
				t.Errorf(
					"POST beer json %s at endpoint %q returned status %d, want %d",
					string(jsonB),
					endpoint,
					response.status,
					http.StatusBadRequest,
				)
			}

			if !strings.Contains(string(response.body), test.want) {
				t.Errorf(
					"POST beer json %s\nat endpoint %q\nreturned body: %s\nwant body: %s",
					string(jsonB),
					endpoint,
					string(response.body),
					test.want,
				)
			}
		})
	}
}

type resp struct {
	status int
	body   []byte
}

func sendReq(t *testing.T, method string, url string, reader io.Reader) resp {
	r, err := http.NewRequest(method, url, reader)
	if err != nil {
		t.Fatalf("creating an HTTP request with method %q and URL %q failed: %s", method, url, err)
	}

	response, err := client.Do(r)
	if err != nil {
		t.Fatalf("sending request with Do() failed: %s", err)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("reading http.Response.Body failed: %s", err)
	}

	return resp{
		status: response.StatusCode,
		body:   body,
	}
}

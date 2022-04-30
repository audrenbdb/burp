package httpreq_test

import (
	"burp"
	"burp/httpreq"
	"burp/mock"
	"bytes"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPutBeer(t *testing.T) {
	t.Run("When beer has different name in url and in JSON body then it should return bad request", func(t *testing.T) {
		guinness := burp.Beer{Name: "Guinness", Rating: burp.Average}
		payload, _ := json.Marshal(guinness)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPut, "/v1/beers/Heineken", bytes.NewBuffer(payload))

		h := httpreq.NewHandler(nil)
		h.ServeHTTP(w, r)
		if w.Code != http.StatusBadRequest {
			t.Errorf("got: %d, want: %d", w.Code, http.StatusBadRequest)
		}
	})

	t.Run("When saving beer is unauthorized then it should return unauthorized", func(t *testing.T) {
		guinness := burp.Beer{Name: "Guinness", Rating: burp.Average}
		payload, _ := json.Marshal(guinness)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPut, "/v1/beers/"+guinness.Name, bytes.NewBuffer(payload))

		ctrl := gomock.NewController(t)
		s := mock.NewService(ctrl)
		s.EXPECT().
			SaveBeer(r.Context(), guinness).
			Return(burp.ErrUnauthorized{})
		h := httpreq.NewHandler(s)
		h.ServeHTTP(w, r)
		if w.Code != http.StatusUnauthorized {
			t.Errorf("got: %d, want: %d", w.Code, http.StatusUnauthorized)
		}
	})
	t.Run("When saving beer is successful then it should return beer in body response as JSON", func(t *testing.T) {
		heineken := burp.Beer{Name: "Heineken", Rating: burp.VeryBad}
		payload, _ := json.Marshal(heineken)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPut, "/v1/beers/"+heineken.Name, bytes.NewBuffer(payload))

		ctrl := gomock.NewController(t)
		s := mock.NewService(ctrl)
		s.EXPECT().
			SaveBeer(r.Context(), heineken).
			Return(nil)
		h := httpreq.NewHandler(s)
		h.ServeHTTP(w, r)
		if w.Code != http.StatusAccepted {
			t.Errorf("got: %d, want: %d", w.Code, http.StatusAccepted)
		}
		if body := w.Body.String(); !strings.Contains(body, string(payload)) {
			t.Errorf("unexpected body, got: %s, want: %s", body, string(payload))
		}
	})
}

func TestHandeDeleteBeer(t *testing.T) {
	t.Run("When deleting a beer is unauthorized then it should return unauthorized", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodDelete, "/v1/beers/Guinness", nil)

		ctrl := gomock.NewController(t)
		s := mock.NewService(ctrl)
		s.EXPECT().
			DeleteBeer(r.Context(), "Guinness").
			Return(burp.ErrUnauthorized{})
		h := httpreq.NewHandler(s)
		h.ServeHTTP(w, r)
		if w.Code != http.StatusUnauthorized {
			t.Errorf("got: %d, want: %d", w.Code, http.StatusUnauthorized)
		}
	})
	t.Run("When deleting a beer is successful then it should return no content", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodDelete, "/v1/beers/Heineken", nil)

		ctrl := gomock.NewController(t)
		s := mock.NewService(ctrl)
		s.EXPECT().
			DeleteBeer(r.Context(), "Heineken").
			Return(nil)
		h := httpreq.NewHandler(s)
		h.ServeHTTP(w, r)
		if w.Code != http.StatusNoContent {
			t.Errorf("got: %d, want: %d", w.Code, http.StatusNoContent)
		}
	})
}

func TestHandleGetBeer(t *testing.T) {
	t.Run("When beer is not found then it should return not found", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/v1/beers/Heineken", nil)

		ctrl := gomock.NewController(t)
		s := mock.NewService(ctrl)
		s.EXPECT().
			GetBeerByName(r.Context(), "Heineken").
			Return(burp.Beer{}, burp.ErrNotFound{ResourceName: "heineken"})

		h := httpreq.NewHandler(s)
		h.ServeHTTP(w, r)
		if w.Code != http.StatusNotFound {
			t.Errorf("got: %d, want: %d", w.Code, http.StatusNotFound)
		}
	})

	t.Run("When beer exist then it should be returned", func(t *testing.T) {
		heineken := burp.Beer{Name: "Heineken", Rating: burp.VeryGood}
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/v1/beers/Heineken", nil)

		ctrl := gomock.NewController(t)
		s := mock.NewService(ctrl)
		s.EXPECT().
			GetBeerByName(r.Context(), "Heineken").
			Return(heineken, nil)

		h := httpreq.NewHandler(s)
		h.ServeHTTP(w, r)
		wantPayload, _ := json.Marshal(heineken)
		if body := w.Body.String(); !strings.Contains(body, string(wantPayload)) {
			t.Errorf("got body: %s, want: %s", body, string(wantPayload))
		}
	})
}

func TestGetBeers(t *testing.T) {
	t.Run("When searching for all beers and none is found, it should return an empty list", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/v1/beers?name=Gu&ratings=1,3", nil)

		ratingsFilter := []burp.Rating{burp.Bad, burp.Good}
		nameFilter := "Gu"

		ctrl := gomock.NewController(t)
		s := mock.NewService(ctrl)
		s.EXPECT().
			SearchBeers(r.Context(), burp.BeerFilter{
				NameContains: &nameFilter,
				Ratings:      &ratingsFilter,
			}).
			Return([]burp.Beer{}, nil)

		h := httpreq.NewHandler(s)
		h.ServeHTTP(w, r)

		if body := w.Body.String(); body != "[]\n" {
			t.Errorf("want empty array of beers, got: %s", body)
		}
	})
}

package burp_test

import (
	"burp"
	"burp/burptest"
	"errors"
	"testing"
	"time"
)

func TestValidateBeerWithoutID(t *testing.T) {
	beer := burptest.RandBeer()
	beer.ID = burp.ID{}
	want := burp.ErrIDEmpty

	err := beer.Validate()
	if !errors.Is(err, want) {
		t.Errorf("RandBeer %+v Validate() got error %s, want %s", beer, err, want)
	}
}

func TestValidateBeerWithoutCreationDate(t *testing.T) {
	beer := burptest.RandBeer()
	beer.CreatedAt = time.Time{}
	want := burp.ErrBeerCreateDateMissing

	err := beer.Validate()
	if !errors.Is(err, want) {
		t.Errorf("RandBeer %+v Validate() returnd error %s, want %s", beer, err, want)
	}
}

func TestValidateBeerWithoutUpdateDate(t *testing.T) {
	beer := burptest.RandBeer()
	beer.UpdatedAt = time.Time{}
	want := burp.ErrBeerUpdateDateMissing

	err := beer.Validate()
	if !errors.Is(err, want) {
		t.Errorf("RandBeer %+v Validate() returnd error %s, want %s", beer, err, want)
	}
}

func TestValidateBeerWithoutPrice(t *testing.T) {
	beer := burptest.RandBeer()
	beer.Price = burp.Price{}
	want := burp.ErrCurrencyNotSupported

	err := beer.Validate()
	if !errors.Is(err, want) {
		t.Errorf("RandBeer %+v Validate() returnd error %s, want %s", beer, err, want)
	}
}

func TestValidateBeerWithoutName(t *testing.T) {
	beer := burptest.RandBeer()
	beer.Name = ""
	want := burp.ErrBeerNameMissing

	err := beer.Validate()
	if !errors.Is(err, want) {
		t.Errorf("RandBeer %+v Validate() returnd error %s, want %s", beer, err, want)
	}
}

func TestValidateBeerWithTooLongName(t *testing.T) {
	beer := burptest.RandBeer()
	beer.Name = burptest.RandString(16)
	want := burp.ErrBeerNameTooLong

	err := beer.Validate()
	if !errors.Is(err, want) {
		t.Errorf("RandBeer %+v Validate() returnd error %s, want %s", beer, err, want)
	}
}

func TestIDValidateWithEmptyUUID(t *testing.T) {
	id := burp.ID{}
	err := id.Validate()
	want := burp.ErrIDEmpty

	if !errors.Is(err, want) {
		t.Errorf("id %q Validate() got error %s, want %s", id, err, want)
	}
}

func TestNewPriceWithInvalidCurrency(t *testing.T) {
	currency := burp.Currency("HELLO")
	amount := uint(1)
	price := burp.Price{Currency: currency, Amount: amount}
	want := burp.ErrCurrencyNotSupported

	err := price.Validate()
	if !errors.Is(err, want) {
		t.Errorf("price %v Validate() got error %s, want %s", price, err, want)
	}
}

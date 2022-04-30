package burp

import "encoding/json"

type Rating int

const (
	VeryBad Rating = iota
	Bad
	Average
	Good
	VeryGood
)

type Beer struct {
	Name   string `json:"name"`
	Rating Rating `json:"rating"`
}

type BeerFilter struct {
	NameContains *string
	Ratings      *[]Rating
}

func (b *Beer) EnsureValid() error {
	if len(b.Name) < 2 {
		return ErrBeerNameTooShort
	}
	if len(b.Name) > 15 {
		return ErrBeerNameTooLong
	}
	if b.Rating < VeryBad || b.Rating > VeryGood {
		return ErrBeerRatingInvalid
	}
	return nil
}

func (b *Beer) UnmarshalJSON(payload []byte) error {
	var beerPayload struct {
		Name   string `json:"name"`
		Rating Rating `json:"rating"`
	}
	err := json.Unmarshal(payload, &beerPayload)
	if err != nil {
		return ErrBadRequest{Hint: "could not decode payload into a beer: " + err.Error()}
	}
	*b = beerPayload
	return b.EnsureValid()
}

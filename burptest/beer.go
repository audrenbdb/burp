package burptest

import (
	"burp"
	"github.com/google/uuid"
	"math/rand"
)

func RandBeer() *burp.Beer {
	return &burp.Beer{
		ID:        burp.ID{UUID: uuid.New()},
		CreatedAt: RandTime(),
		UpdatedAt: RandTime(),

		Name: RandString(15),
		Price: burp.Price{
			Currency: burp.EUR,
			Amount:   uint(rand.Intn(1000-10) + 10),
		},
	}
}

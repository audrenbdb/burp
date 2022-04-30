package burp

import (
	"github.com/google/uuid"
	"time"
)

type Beer struct {
	ID        ID        `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Name  string `json:"name"`
	Price Price  `json:"price"`
}

type ID struct {
	uuid.UUID
}

type Currency string

var (
	EUR Currency = "Euro"
	USD Currency = "Dollar"
)

type Price struct {
	Currency Currency `json:"currency"`
	Amount   uint     `json:"amount"`
}

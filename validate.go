package burp

import "github.com/google/uuid"

func (p Price) Validate() error {
	switch p.Currency {
	case EUR, USD:
	default:
		return ErrCurrencyNotSupported
	}

	return nil
}

func (id ID) Validate() error {
	if id.UUID == uuid.Nil {
		return ErrIDEmpty
	}

	return nil
}

func (b *Beer) Validate() error {
	if err := b.ID.Validate(); err != nil {
		return Errorf("invalid id: %w", err)
	}

	if b.CreatedAt.IsZero() {
		return ErrBeerCreateDateMissing
	}

	if b.UpdatedAt.IsZero() {
		return ErrBeerUpdateDateMissing
	}

	if err := b.Price.Validate(); err != nil {
		return Errorf("invalid price: %w", err)
	}

	if b.Name == "" {
		return ErrBeerNameMissing
	}

	if len(b.Name) > 15 {
		return ErrBeerNameTooLong
	}

	return nil
}

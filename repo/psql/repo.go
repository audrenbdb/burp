package psql

import (
	"burp"
	"burp/repo"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
)

type Repo struct {
	Conn *pgx.Conn
}

func (r *Repo) SaveBeer(ctx context.Context, beer *burp.Beer) error {
	q := `INSERT INTO beer(id, created_at, updated_at, name, price_currency, price_amount) 
	VALUES($1, $2, $3, $4, $5, $6)
	ON CONFLICT (id)
	DO
	UPDATE SET created_at = $2, updated_at = $3, name = $4, price_currency = $5, price_amount = $6`

	_, err := r.Conn.Exec(ctx, q, beer.ID, beer.CreatedAt, beer.UpdatedAt, beer.Name, beer.Price.Currency, beer.Price.Amount)
	if err != nil {
		return repo.Error(err.Error())
	}

	return nil
}

func (r *Repo) RemoveBeer(ctx context.Context, id burp.ID) error {
	q := `DELETE FROM beer WHERE id = $1`

	_, err := r.Conn.Exec(ctx, q, id)
	return err
}

func (r *Repo) SelectBeer(ctx context.Context, id burp.ID) (*burp.Beer, error) {
	var beer burp.Beer
	q := `SELECT id, created_at, updated_at, name, price_currency, price_amount FROM beer WHERE id = $1`
	row := r.Conn.QueryRow(ctx, q, id)
	err := row.Scan(
		&beer.ID,
		&beer.CreatedAt,
		&beer.UpdatedAt,
		&beer.Name,
		&beer.Price.Currency,
		&beer.Price.Amount,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, repo.Errorf(
			"beer not found with id %q: %w",
			id,
			repo.ErrNotFound,
		)
	}

	return &beer, err
}

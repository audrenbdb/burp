package repo

import (
	"burp"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"os"
	"strconv"
	"strings"
)

func NewPSQL() burp.Repo {
	conn, err := pgxpool.Connect(context.Background(), os.Getenv("BURP_PSQL_URL"))
	if err != nil {
		log.Fatal(err)
	}
	return &psqlRepo{conn: conn}
}

type psqlRepo struct {
	conn *pgxpool.Pool
}

func (r *psqlRepo) GetBeerByName(ctx context.Context, name string) (beer burp.Beer, err error) {
	beer.Name = name
	query := "SELECT rating FROM beers WHERE name = $1"
	row := r.conn.QueryRow(ctx, query, name)
	return beer, row.Scan(&beer.Rating)
}

func (r *psqlRepo) GetBeers(ctx context.Context, filter burp.BeerFilter) ([]burp.Beer, error) {
	beers := make([]burp.Beer, 0)
	query, args := psqlQueryFromBeerFilter(filter)
	rows, err := r.conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var b burp.Beer
		err = rows.Scan(&b.Name, &b.Rating)
		if err != nil {
			return nil, err
		}
		beers = append(beers, b)
	}

	return beers, nil
}

func psqlQueryFromBeerFilter(filter burp.BeerFilter) (query string, args []any) {
	query = "SELECT name, rating FROM beers WHERE true"
	if filter.Ratings != nil {
		var ratings []string
		for _, r := range *filter.Ratings {
			ratings = append(ratings, strconv.Itoa(int(r)))
			query += fmt.Sprintf("$%d", len(args)+1)
		}
		query += " AND rating IN("
		query += strings.Join(ratings, ",")
		query += ")"
	}
	if filter.NameContains != nil {
		query += fmt.Sprintf(" AND name LIKE '%%' || $%d || '%%'", len(args)+1)
		args = append(args, filter.NameContains)
	}
	query += " ORDER BY created_at ASC"
	return query, args
}

func (r *psqlRepo) DeleteBeer(ctx context.Context, name string) error {
	query := "DELETE FROM beers WHERE name = $1"
	_, err := r.conn.Exec(ctx, query, name)
	return err
}

func (r *psqlRepo) SaveBeer(ctx context.Context, beer burp.Beer) error {
	insert := "INSERT INTO beers(name, rating) VALUES($1, $2)"
	update := "ON CONFLICT (name) DO UPDATE SET name = $1, rating = $2"
	upsert := fmt.Sprintf("%s %s", insert, update)
	_, err := r.conn.Exec(ctx, upsert, beer.Name, beer.Rating)
	return err
}

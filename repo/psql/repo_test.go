package psql_test

import (
	"burp"
	"burp/burptest"
	"burp/repo"
	"errors"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"strings"
	"testing"
)

func TestSaveBeer(t *testing.T) {
	beer := burptest.RandBeer()

	err := appRepo.SaveBeer(ctx, beer)
	if err != nil {
		t.Errorf("SaveBeer(ctx, %+v) should not return an error", beer)
	}

	selectQuery := "SELECT id, created_at, updated_at, name, price_currency, price_amount FROM beer WHERE id = $1"

	got, err := scanBeerRow(conn.QueryRow(ctx, selectQuery, beer.ID))
	if err != nil {
		t.Fatalf("Querying %selectQuery should not fail, got %s", selectQuery, err)
	}

	if diff := cmp.Diff(beer, got); diff != "" {
		t.Errorf("SaveBeer(ctx, %+v did not save new beer in database (-want/+got):\n%s", beer, diff)
	}

	beer.UpdatedAt = burptest.RandTime()
	err = appRepo.SaveBeer(ctx, beer)
	if err != nil {
		t.Errorf("SaveBeer(ctx, %+v) returned unexpected error: %s", beer, err)
	}

	got, err = scanBeerRow(conn.QueryRow(ctx, selectQuery, beer.ID))
	if err != nil {
		t.Fatalf("Querying %selectQuery should not fail, got %s", selectQuery, err)
	}

	if diff := cmp.Diff(beer, got); diff != "" {
		t.Errorf("SaveBeer(ctx, %+v did not update existing beer in database (-want/+got):\n%s", beer, diff)
	}
}

func TestRemoveBeer(t *testing.T) {
	beer := burptest.RandBeer()

	insertBeer(t, beer)

	err := appRepo.RemoveBeer(ctx, beer.ID)
	if err != nil {
		t.Errorf("RemoveBeer(ctx, %s) returnd error %s, want none", beer.ID, err)
	}

	var count int
	conn.QueryRow(ctx, "SELECT COUNT(*) FROM beer WHERE id = $1", beer.ID).Scan(&count)

	if count != 0 {
		t.Errorf("Selecting beer count after deletion returned %d, want 0", count)
	}
}

func TestSaveBeerWithRepoErr(t *testing.T) {
	beer := burptest.RandBeer()
	beer.Name = burptest.RandString(300)

	err := appRepo.SaveBeer(ctx, beer)
	if !errors.As(err, &repo.Err{}) {
		t.Errorf("SaveBeer(ctx, %+v) returned error %s, want error of type app.RepoErr{}", beer, err)
	}
}

func TestSelectBeer(t *testing.T) {
	beer := burptest.RandBeer()

	insertBeer(t, beer)

	got, err := appRepo.SelectBeer(ctx, beer.ID)
	if err != nil {
		t.Errorf("SelectBeer(ctx, %+v) returned error %s, want none", beer, err)
	}

	if diff := cmp.Diff(beer, got); diff != "" {
		t.Errorf("SelectBeer(ctx, %+v) returned unexpected beer, (-want/+got):\n%q", beer, diff)
	}
}

func TestSelectBeerNotFound(t *testing.T) {
	var repoErr repo.Err
	id := burp.ID{UUID: uuid.New()}
	want := fmt.Sprintf("beer not found with id %q", id)

	_, err := appRepo.SelectBeer(ctx, id)
	if !errors.Is(err, repo.ErrNotFound) || !errors.As(err, &repoErr) || !strings.Contains(err.Error(), want) {
		t.Errorf("SelectBeer(ctx, %q) got error %s, want app.RepoErr{} that contains %q", id, err, want)
	}
}

func insertBeer(t *testing.T, beer *burp.Beer) {
	t.Helper()

	insertQuery := `INSERT INTO beer(id, created_at, updated_at, name, price_currency, price_amount)
	VALUES($1, $2, $3, $4, $5, $6)`

	_, err := conn.Exec(
		ctx,
		insertQuery,
		beer.ID,
		beer.CreatedAt,
		beer.UpdatedAt,
		beer.Name,
		beer.Price.Currency,
		beer.Price.Amount,
	)
	if err != nil {
		t.Fatalf("Executing query %q returned error %v, required none", insertQuery, err)
	}
}

func scanBeerRow(row pgx.Row) (*burp.Beer, error) {
	var got burp.Beer
	err := row.Scan(
		&got.ID,
		&got.CreatedAt,
		&got.UpdatedAt,
		&got.Name,
		&got.Price.Currency,
		&got.Price.Amount,
	)
	return &got, err
}

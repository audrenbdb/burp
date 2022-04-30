package psql_test

import (
	"burp/repo/psql"
	"context"
	_ "embed"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"log"
	"os"
	"strings"
	"testing"
	"time"
)

var (
	ctx     context.Context
	conn    *pgx.Conn
	appRepo *psql.Repo

	//go:embed schema.sql
	schema string
)

func TestMain(m *testing.M) {
	tctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	ctx = tctx

	d, err := os.MkdirTemp("", "migration")
	if err != nil {
		log.Fatal("Could not create a temporary directory")
	}
	defer os.RemoveAll(d)

	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "latest",
		Env: []string{
			"POSTGRES_PASSWORD=secret",
			"POSTGRES_USER=user_name",
			"POSTGRES_DB=dbname",
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	hostAndPort := resource.GetHostPort("5432/tcp")
	databaseUrl := fmt.Sprintf("postgres://user_name:secret@%s/dbname?sslmode=disable", hostAndPort)

	log.Println("Connecting to database on url: ", databaseUrl)

	resource.Expire(120) // Tell docker to hard kill the container in 120 seconds

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	pool.MaxWait = 120 * time.Second
	if err = pool.Retry(func() error {
		conn, err = pgx.Connect(context.Background(), databaseUrl)
		if err != nil {
			return err
		}
		return conn.Ping(context.Background())
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	if err = initDatabase(); err != nil {
		log.Fatalf("Could not initialize database: %s", err)
	}

	appRepo = &psql.Repo{Conn: conn}

	//Run tests
	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func initDatabase() error {
	statements := strings.Split(schema, ";")
	for _, statement := range statements {
		_, err := conn.Exec(ctx, statement)
		if err != nil {
			return err
		}
	}
	return nil
}

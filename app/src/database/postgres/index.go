package postgres_connection

import (
	custom_errors "cloudview/app/src/api/errors"
	"cloudview/app/src/helpers"
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-jet/jet/v2/qrm"
	_ "github.com/lib/pq"
)

type PostgresStore struct {
	Postgres qrm.DB
	RawDB    *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	var (
		DB_NAME = helpers.GoDotEnvVariable("DB_NAME")
		DB_HOST = helpers.GoDotEnvVariable("DB_HOST")
		DB_PORT = 5432
		// DB_PORT     = helpers.GoDotEnvVariable("DB_PORT")
		DB_PASSWORD = helpers.GoDotEnvVariable("DB_PASSWORD")
		DB_USER     = helpers.GoDotEnvVariable("DB_USER")
	)
	if DB_HOST == "" || DB_USER == "" || DB_NAME == "" || DB_PASSWORD == "" {
		return nil, errors.New("Missing Database credentials. Make sure the env vars are valid and accessible")
	}
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, custom_errors.DBErrors(err)
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	// defer db.Close()
	return &PostgresStore{
		Postgres: db,
		RawDB:    db,
	}, nil
}

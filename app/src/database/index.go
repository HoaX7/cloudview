package database

import postgres_connection "cloudview/app/src/database/postgres"

type DB struct {
	*postgres_connection.PostgresStore
}

func NewDB() (*DB, error) {
	db, err := postgres_connection.NewPostgresStore()
	return &DB{db}, err
}

package postgres

import (
	"database/sql"
	_ "github.com/lib/pq"
	"time"
)

type DB struct {
	Conn *sql.DB
}

func New(dsn string) (*DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)
	return &DB{Conn: db}, nil
}
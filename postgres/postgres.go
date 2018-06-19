// Package postgres is the package, that contains queries for handling databases
// It handles two databases: players and tournaments
package postgres

import (
	"database/sql"

	"github.com/Tournament/errors"
)

// Postgres is the struct that contains DB and implements methods for handling databases
type Postgres struct {
	DB *sql.DB
}

// NewDB returns postgres database with configuration conf
func NewDB(conf string) (*Postgres, error) {
	db, err := sql.Open("postgres", conf)
	if err != nil {
		return nil, errors.Error{Code: errors.DatabaseOpenError, Message: "cannot open database, config: " + conf}
	}
	return &Postgres{DB: db}, nil
}

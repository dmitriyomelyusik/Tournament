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

// UpdateTourAndPlayer updates tournament participants and player balance in one transaction
func (p *Postgres) UpdateTourAndPlayer(tourID, playerID string) error {
	tx, err := p.DB.Begin()
	if err != nil {
		return errors.Error{Code: errors.UnexpectedError, Message: "update tournament and player: failed to start transaction"}
	}
	err = updateTxParticipants(tx, tourID, playerID)
	if err != nil {
		err2 := tx.Rollback()
		return errors.Join(err, err2)
	}
	dep, err := p.GetDeposit(tourID)
	if err != nil {
		err2 := tx.Rollback()
		return errors.Join(err, err2)
	}
	err = updateTxPlayer(tx, playerID, -1*dep)
	if err != nil {
		err2 := tx.Rollback()
		return errors.Join(err, err2)
	}
	return tx.Commit()
}

// ResultError checks query result and returns errors, if they occured
func ResultError(res sql.Result, possibleErr string) error {
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n != 1 {
		return errors.Error{Code: errors.NotFoundError, Message: possibleErr}
	}
	return nil
}

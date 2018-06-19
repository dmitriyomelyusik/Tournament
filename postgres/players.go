package postgres

import (
	"strconv"

	"github.com/Tournament/entity"
	"github.com/Tournament/errors"
)

// CreatePlayer creates new player with id and points
func (p *Postgres) CreatePlayer(id string, points uint) error {
	res, err := p.DB.Exec("INSERT INTO players (id, points) values ($1, $2)", id, points)
	if err != nil {
		return errors.Error{Code: errors.DuplicatedIDError, Message: "using duplicated id to create player"}
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n != 1 {
		return errors.Error{Code: errors.UnexpectedError, Message: "cannot create player with id " + id}
	}
	return nil
}

// GetPlayer returns player by its id
func (p *Postgres) GetPlayer(id string) (entity.Player, error) {
	row := p.DB.QueryRow("SELECT points FROM players WHERE id=$1", id)
	var points uint
	err := row.Scan(&points)
	if err != nil {
		return entity.Player{}, errors.Error{Code: errors.PlayerNotFoundError, Message: "cannot find player with id " + id}
	}
	return entity.Player{ID: id, Points: points}, nil
}

// UpdatePlayer updates player number of points
func (p *Postgres) UpdatePlayer(id string, dif int) error {
	res, err := p.DB.Exec("UPDATE players SET points=points+$1 WHERE id=$2", dif, id)
	if err != nil {
		return errors.Error{Code: errors.NegativePointsNumberError, Message: "cannot update points numbers with dif " + strconv.Itoa(dif)}
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n != 1 {
		return errors.Error{Code: errors.UnexpectedError, Message: "cannot update player points: id " + id + ", dif" + strconv.Itoa(dif)}
	}
	return nil
}

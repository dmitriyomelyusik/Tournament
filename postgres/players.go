package postgres

import (
	"database/sql"
	"strconv"

	"github.com/dmitriyomelyusik/Tournament/entity"
	"github.com/dmitriyomelyusik/Tournament/errors"
)

// CreatePlayer creates new player with id and points
func (p *Postgres) CreatePlayer(id string, points int) (entity.Player, error) {
	res, err := p.DB.Exec("INSERT INTO players (id, points) values ($1, $2)", id, points)
	if err != nil {
		return entity.Player{}, errors.Error{Code: errors.DuplicatedIDError, Message: "create player: using duplicated id to create player, id " + id}
	}
	err = resultError(res, "creating player: cannot create player, id "+id)
	if err != nil {
		return entity.Player{}, err
	}
	return entity.Player{ID: id, Points: points}, nil
}

// GetPlayer returns player by its id
func (p *Postgres) GetPlayer(id string) (entity.Player, error) {
	row := p.DB.QueryRow("SELECT points FROM players WHERE id=$1", id)
	var points int
	err := row.Scan(&points)
	if err != nil {
		return entity.Player{}, errors.Error{Code: errors.NotFoundError, Message: "get player: cannot find player, id " + id}
	}
	return entity.Player{ID: id, Points: points}, nil
}

// UpdatePlayer updates player points
func (p *Postgres) UpdatePlayer(id string, dif int) error {
	res, err := p.DB.Exec("UPDATE players SET points=points+$1 WHERE id=$2", dif, id)
	if err != nil {
		return errors.Error{Code: errors.NegativePointsNumberError, Message: "update player: cannot update points numbers, dif " + strconv.Itoa(dif)}
	}
	return resultError(res, "update player: cannot find player, id "+id)
}

func updateTxPlayer(tx *sql.Tx, id string, dif int) error {
	res, err := tx.Exec("UPDATE players SET points=points+$1 WHERE id=$2", dif, id)
	if err != nil {
		return errors.Error{Code: errors.NegativePointsNumberError, Message: "update player: cannot update points numbers, dif " + strconv.Itoa(dif)}
	}
	return resultError(res, "update player: cannot find player, id "+id)
}

// DeletePlayer deletes player from database
func (p *Postgres) DeletePlayer(id string) error {
	res, err := p.DB.Exec("DELETE FROM players WHERE id=$1", id)
	if err != nil {
		return errors.Error{Code: errors.UnexpectedError, Message: "delete player: " + err.Error()}
	}
	return resultError(res, "delete player: player does not exist, id "+id)
}

package postgres

import (
	"database/sql"
	"strconv"

	"github.com/Tournament/entity"
	"github.com/Tournament/errors"
)

// CreatePlayer creates new player with id and points
func (p *Postgres) CreatePlayer(id string, points int) error {
	res, err := p.DB.Exec("INSERT INTO players (id, points) values ($1, $2)", id, points)
	if err != nil {
		return errors.Error{Code: errors.DuplicatedIDError, Message: "creating player: using duplicated id to create player"}
	}
	return resultError(res, "creating player: cannot create player with id "+id)
}

// GetPlayer returns player by its id
func (p *Postgres) GetPlayer(id string) (entity.Player, error) {
	row := p.DB.QueryRow("SELECT points FROM players WHERE id=$1", id)
	var points int
	err := row.Scan(&points)
	if err != nil {
		return entity.Player{}, errors.Error{Code: errors.PlayerNotFoundError, Message: "getting player: cannot find player with id " + id}
	}
	return entity.Player{ID: id, Points: points}, nil
}

// UpdatePlayer updates player number of points
func (p *Postgres) UpdatePlayer(id string, dif int) error {
	res, err := p.DB.Exec("UPDATE players SET points=points+$1 WHERE id=$2", dif, id)
	if err != nil {
		return errors.Error{Code: errors.NegativePointsNumberError, Message: "updating player: cannot update points numbers with dif " + strconv.Itoa(dif)}
	}
	return resultError(res, "updating player: cannot found player with id "+id)
}

func updateTxPlayer(tx *sql.Tx, id string, dif int) error {
	res, err := tx.Exec("UPDATE players SET points=points+$1 WHERE id=$2", dif, id)
	if err != nil {
		return errors.Error{Code: errors.NegativePointsNumberError, Message: "updating player: cannot update points numbers with dif " + strconv.Itoa(dif)}
	}
	return resultError(res, "updating player: cannot find player: id "+id)
}

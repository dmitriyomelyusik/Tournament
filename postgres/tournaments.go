package postgres

import (
	"encoding/json"

	"github.com/Tournament/entity"
	"github.com/Tournament/errors"
)

// CreateTournament creates tournament with id and deposit
func (p *Postgres) CreateTournament(id string, deposit uint) error {
	res, err := p.DB.Exec("INSERT INTO tournaments (id, deposit, prize) values ($1, $2, '0')", id, deposit)
	if err != nil {
		return errors.Error{Code: errors.DuplicatedIDError, Message: "using duplicated id to create tournament, id: " + id}
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n != 1 {
		return errors.Error{Code: errors.UnexpectedError, Message: "cannot create tournament with id " + id}
	}
	return nil
}

// UpdateParticipants add new paricipant in the tournament
func (p *Postgres) UpdateParticipants(id string, player entity.Player) error {
	rawPlayer, err := json.Marshal(player)
	if err != nil {
		return errors.Error{Code: errors.JSONError, Message: "updating participants: cannot marshal player, id: " + player.ID}
	}
	res, err := p.DB.Exec("UPDATE tournaments SET participants=array_append(participants, $1) WHERE id=$2", rawPlayer, id)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n != 1 {
		return errors.Error{Code: errors.TournamentNotFoundError, Message: "cannot update participants from not existing tournament, id: " + id}
	}
	return nil
}

// SetTournamentWinner sets winner
func (p *Postgres) SetTournamentWinner(id string, player entity.Player) error {
	rawPlayer, err := json.Marshal(player)
	if err != nil {
		return errors.Error{Code: errors.JSONError, Message: "setting winner: cannot marshal player, id: " + player.ID}
	}
	res, err := p.DB.Exec("UPDATE tournaments SET winner=$1 WHERE id=$2", rawPlayer, id)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n != 1 {
		return errors.Error{Code: errors.TournamentNotFoundError, Message: "cannot set winner in not existing tournament, id: " + id}
	}
	return nil
}

// GetParticipants returns all participants in the tournament
func (p *Postgres) GetParticipants(id string) ([]entity.Player, error) {
	row := p.DB.QueryRow("SELECT participant FROM tournaments WHERE id=$1", id)
	var players []entity.Player
	err := row.Scan(players)
	if err != nil {
		return nil, errors.Error{Code: errors.TournamentNotFoundError, Message: "cannot get participants from not existing tournament, id: " + id}
	}
	return players, nil
}

// GetDeposit returns tournament deposit
func (p *Postgres) GetDeposit(id string) (uint, error) {
	row := p.DB.QueryRow("SELECT deposit FROM tournaments WHERE id=$1", id)
	var deposit uint
	err := row.Scan(&deposit)
	if err != nil {
		return 0, errors.Error{Code: errors.TournamentNotFoundError, Message: "cannot get deposit from not existing tournament, id: " + id}
	}
	return deposit, nil
}

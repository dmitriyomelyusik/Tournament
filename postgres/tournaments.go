package postgres

import (
	"database/sql"
	"encoding/json"

	"github.com/lib/pq"

	"github.com/Tournament/entity"
	"github.com/Tournament/errors"
)

// CloseTournament closes tournament in transaction
func (p *Postgres) CloseTournament(id string) error {
	res, err := p.DB.Exec("UPDATE tournaments SET isOpen='false' WHERE id=$1", id)
	if err != nil {
		return err
	}
	return ResultError(res, "closing tournament: cannot close not existing tournament, id: "+id)
}

// CreateTournament creates tournament with id and deposit
func (p *Postgres) CreateTournament(id string, deposit int) error {
	res, err := p.DB.Exec("INSERT INTO tournaments (id, deposit, prize, isOpen) values ($1, $2, '0', 'true')", id, deposit)
	if err != nil {
		return errors.Error{Code: errors.DuplicatedIDError, Message: "creating tournament: using duplicated id to create tournament, id: " + id, Info: err.Error()}
	}
	return ResultError(res, "creating tournament: cannot create tournament with id "+id)
}

// GetParticipants returns tournament participants
func (p *Postgres) GetParticipants(id string) ([]string, error) {
	row := p.DB.QueryRow("SELECT participants FROM tournaments WHERE id=$1", id)
	var playerIDs []string
	err := row.Scan(pq.Array(&playerIDs))
	if err != nil {
		return nil, errors.Error{Code: errors.NotFoundError, Message: "getting participants: cannot get participants from not existing tournament, id: " + id, Info: err.Error()}
	}
	return playerIDs, nil
}

// GetTournamentState returns true, if tournament opens for joining
func (p *Postgres) GetTournamentState(id string) (bool, error) {
	row := p.DB.QueryRow("SELECT isOpen FROM tournaments WHERE id=$1", id)
	var isOpen bool
	err := row.Scan(&isOpen)
	if err != nil {
		return false, errors.Error{Code: errors.NotFoundError, Message: "getting state: cannot get tournament state from not existing tournament, id: " + id, Info: err.Error()}
	}
	return isOpen, nil
}

// GetWinner returns tournament winner
func (p *Postgres) GetWinner(id string) (entity.Winners, error) {
	row := p.DB.QueryRow("SELECT winner FROM tournaments WHERE id=$1", id)
	var rawWinner []byte
	err := row.Scan(&rawWinner)
	if err != nil {
		return entity.Winners{}, errors.Error{Code: errors.NotFoundError, Message: "getting winner: cannot get winner from not existing tournament, id: " + id, Info: err.Error()}
	}
	var winner entity.Winner
	err = json.Unmarshal(rawWinner, &winner)
	if err != nil {
		return entity.Winners{}, errors.Error{Code: errors.JSONError, Message: "getting winner: cannot unmarshal winner, id: " + winner.ID, Info: err.Error()}
	}
	return entity.Winners{Winners: []entity.Winner{winner}}, nil
}

// GetDeposit returns tournament deposit
func (p *Postgres) GetDeposit(id string) (int, error) {
	row := p.DB.QueryRow("SELECT deposit FROM tournaments WHERE id=$1", id)
	var deposit int
	err := row.Scan(&deposit)
	if err != nil {
		return 0, errors.Error{Code: errors.NotFoundError, Message: "getting deposit: cannot get deposit from not existing tournament, id: " + id, Info: err.Error()}
	}
	return deposit, nil
}

// SetTournamentWinner sets winner
func (p *Postgres) SetTournamentWinner(id string, winner entity.Winner) error {
	tx, err := p.DB.Begin()
	if err != nil {
		return err
	}
	row := tx.QueryRow("SELECT prize FROM tournaments WHERE id=$1", id)
	var prize int
	err = row.Scan(&prize)
	if err != nil {
		err2 := tx.Rollback()
		return errors.Join(err, err2).SetPrefix("setting winner: tournament not exist, id: " + id + "\n").SetCode(errors.NotFoundError)
	}
	err = updateTxPlayer(tx, winner.ID, prize)
	if err != nil {
		err2 := tx.Rollback()
		return errors.Join(err, err2).SetPrefix("setting winner: ")
	}
	winner.Prize = prize
	rawWinner, err := json.Marshal(winner)
	if err != nil {
		err2 := tx.Rollback()
		return errors.Join(err, err2).SetPrefix("setting winner: cannot marshal player").SetCode(errors.JSONError)
	}
	res, err := tx.Exec("UPDATE tournaments SET winner=$1 WHERE id=$2", rawWinner, id)
	if err != nil {
		err2 := tx.Rollback()
		return errors.Join(err, err2).SetPrefix("setting winner: ")
	}
	err = ResultError(res, "setting winner: cannot close not existing tournament, id: "+id)
	if err != nil {
		err2 := tx.Rollback()
		return errors.Join(err, err2).SetPrefix("setting winner: ")
	}
	return tx.Commit()
}

func updateTxParticipants(tx *sql.Tx, tourID, playerID string) error {
	res, err := tx.Exec("UPDATE tournaments SET participants=array_append(participants, $1), prize=prize+deposit WHERE id=$2", playerID, tourID)
	if err != nil {
		return err
	}
	return ResultError(res, "updating participiants: cannot update participants from not existing tournament, id: "+tourID)
}

// DeleteTournament deletes tournament
func (p *Postgres) DeleteTournament(id string) error {
	res, err := p.DB.Exec("DELETE FROM tournaments WHERE id=$1", id)
	if err != nil {
		return err
	}
	return ResultError(res, "deleting tournament: tournament does not exist, id "+id)
}

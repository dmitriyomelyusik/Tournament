package mongo

import (
	"github.com/dmitriyomelyusik/Tournament/entity"
	"github.com/dmitriyomelyusik/Tournament/errors"
	"gopkg.in/mgo.v2/bson"
)

// CreateTournament creates tournament with id and deposit
func (m *Mongo) CreateTournament(id string, deposit int) error {
	err := m.tournaments.Insert(bson.M{"_id": id, "deposit": deposit, "isOpen": true, "participants": []string{}, "prize": 0, "winner": entity.Winners{}})
	if err != nil {
		return errors.Error{Code: errors.UnexpectedError, Message: err.Error()}.SetPrefix("create tournament: ")
	}
	return nil
}

// CloseTournament closes tournament in transaction
func (m *Mongo) CloseTournament(id string) error {
	err := m.tournaments.UpdateId(id, bson.M{"$set": bson.M{"isOpen": false}})
	if err != nil {
		return errors.Error{Code: errors.NotFoundError, Message: "close tournament: tournament is not found, id " + id}
	}
	return nil
}

// GetParticipants returns tournament participants
func (m *Mongo) GetParticipants(id string) ([]string, error) {
	var participants []string
	err := m.tournaments.FindId(id).Select(bson.M{"participants": 1}).One(participants)
	if err != nil {
		return nil, errors.Error{Code: errors.NotFoundError, Message: "get participants: tournament is not found, id " + id}
	}
	return participants, nil
}

// GetTournamentState returns true, if tournament opens for joining
func (m *Mongo) GetTournamentState(id string) (bool, error) {
	var state bool
	err := m.tournaments.FindId(id).Select(bson.M{"isOpen": 1}).One(&state)
	if err != nil {
		return false, errors.Error{Code: errors.NotFoundError, Message: "get tournament state: tournament is not found, id " + id}
	}
	return state, nil
}

// GetWinner returns tournament winner
func (m *Mongo) GetWinner(id string) (entity.Winners, error) {
	var winner entity.Winners
	err := m.tournaments.FindId(id).Select(bson.M{"winners": 1}).One(&winner)
	if err != nil {
		return entity.Winners{}, errors.Error{Code: errors.NotFoundError, Message: "get winner: tournament is not found, id " + id}
	}
	return winner, nil
}

// SetTournamentWinner sets winner in one transaction
func (m *Mongo) SetTournamentWinner(id string, winner entity.Winner) error {

	return nil
}

// DeleteTournament deletes tournament
func (m *Mongo) DeleteTournament(id string) error {
	err := m.tournaments.RemoveId(id)
	if err != nil {
		return errors.Error{Code: errors.NotFoundError, Message: "delete tournament: tournament is not found, id " + id}
	}
	return nil
}

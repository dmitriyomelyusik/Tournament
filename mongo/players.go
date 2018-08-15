package mongo

import (
	"log"
	"strconv"

	"github.com/dmitriyomelyusik/Tournament/entity"
	"github.com/dmitriyomelyusik/Tournament/errors"
	"github.com/dmitriyomelyusik/Tournament/mongo/logs"
	"gopkg.in/mgo.v2/bson"
)

// CreatePlayer creates new player with id and points
func (m *Mongo) CreatePlayer(id string, points int) (entity.Player, error) {
	player := entity.Player{ID: id, Points: points}
	err := m.players.Insert(player)
	if err != nil {
		return entity.Player{}, err
	}
	err = m.logger.Log(id, logger.Fund, points)
	if err != nil {
		return entity.Player{}, err
	}
	return player, nil
}

// GetPlayer returns player by its id
func (m *Mongo) GetPlayer(id string) (entity.Player, error) {
	var p entity.Player
	err := m.players.FindId(id).One(&p)
	return p, err
}

// UpdatePlayer updates player points
func (m *Mongo) UpdatePlayer(id string, points int) error {
	if points < 0 {
		return m.getPoints(id, points)
	}
	return m.setPoints(id, points)
}

func (m *Mongo) setPoints(id string, points int) error {
	err := m.players.UpdateId(id, bson.M{"$inc": bson.M{"points": points}})
	if err != nil {
		return errors.Error{Code: errors.UnexpectedError, Message: err.Error()}
	}
	return m.logger.Log(id, logger.Fund, points)
}

func (m *Mongo) getPoints(id string, points int) error {
	s0, err := logSum(m.logger, id)
	if err != nil {
		return errors.Error{Code: errors.UnexpectedError, Message: err.Error()}
	}
	if s0 < 0 {
		return errors.Error{Code: errors.NegativePointsNumberError, Message: "update player: negative balance, player id " + id}
	}
	if s0 < -points {
		return errors.Error{Code: errors.NegativePointsNumberError, Message: "update player: cannot update points numbers, dif " + strconv.Itoa(points)}
	}
	err = m.players.UpdateId(id, bson.M{"$inc": bson.M{"points": points}})
	if err != nil {
		return errors.Error{Code: errors.UnexpectedError, Message: err.Error()}.SetPrefix("update player: ")
	}
	err = m.logger.Log(id, logger.Take, points)
	if err != nil {
		log.Println(err)
		return m.rollback(id, -points)
	}
	s1, err := logSum(m.logger, id)
	if err != nil {
		log.Println(err)
		return m.rollback(id, -points)
	}
	if s1 < 0 {
		return m.rollback(id, -points)
	}
	return nil
}

func logSum(log *logger.Logger, id string) (int, error) {
	data, err := log.GetLogs(id)
	if err != nil {
		return 0, err
	}
	var sum int
	for _, v := range data {
		sum += v.Points
	}
	return sum, nil
}

func (m *Mongo) rollback(id string, points int) error {
	err := m.players.UpdateId(id, bson.M{"$inc": bson.M{"points": points}})
	if err != nil {
		return errors.Error{Code: errors.CriticalError, Message: "rollback: cannot rollback, next operations can be dangerous", Info: err.Error()}
	}
	return errors.Error{Code: errors.RollbackError, Message: "rollback: got negative balance or disconect, operation aborted"}
}

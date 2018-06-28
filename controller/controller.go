// Package controller has methods to controll, what data send to database queries
package controller

import (
	"math/rand"
	"time"

	"github.com/dmitriyomelyusik/Tournament/entity"
	"github.com/dmitriyomelyusik/Tournament/errors"
)

// PlayerDB is an interface for database, that used to controll player activity methods
type PlayerDB interface {
	GetPlayer(string) (entity.Player, error)
	CreatePlayer(string, int) (entity.Player, error)
	UpdatePlayer(string, int) error
	DeletePlayer(string) error
}

// TourDB is an interface for database, that used to controll tournament activity methods
type TourDB interface {
	CreateTournament(string, int) error
	GetTournamentState(string) (bool, error)
	GetWinner(string) (entity.Winners, error)
	CloseTournament(string) error
	GetParticipants(string) ([]string, error)
	SetTournamentWinner(string, entity.Winner) error
	DeleteTournament(string) error
}

// Database is an interface for database, that uses tournament and player database interfaces
// and adds method to join that two databases
type Database interface {
	PlayerDB
	TourDB
	UpdateTourAndPlayer(string, string) error
}

// Game is a struct which methods controlls activity within database interface
type Game struct {
	DB Database
}

// Fund controlls funding player
func (g Game) Fund(id string, points int) (entity.Player, error) {
	if points < 0 {
		return entity.Player{}, errors.Error{Code: errors.NegativePointsNumberError, Message: "fund: cannot fund negative number of points"}
	}
	if id == "" {
		return entity.Player{}, errors.Error{Code: errors.NotFoundError, Message: "fund: id must be not nil"}
	}
	_, err := g.DB.GetPlayer(id)
	if err != nil {
		return g.DB.CreatePlayer(id, points)
	}
	return entity.Player{}, g.DB.UpdatePlayer(id, points)
}

// Take controlls taking points
func (g Game) Take(id string, points int) error {
	if points < 0 {
		return errors.Error{Code: errors.NegativePointsNumberError, Message: "take: cannot take negative number of points"}
	}
	if id == "" {
		return errors.Error{Code: errors.NotFoundError, Message: "take: id must be not nil"}
	}
	err := g.DB.UpdatePlayer(id, -1*points)
	if err != nil {
		err := errors.Transform(err)
		if err.Code == errors.NegativePointsNumberError {
			err.Message = "take: cannot take points, player doesn't have enough points"
			return err
		}
		return err
	}
	return nil
}

// Balance controlls getting actual player balance
func (g Game) Balance(id string) (entity.Player, error) {
	if id == "" {
		return entity.Player{}, errors.Error{Code: errors.NotFoundError, Message: "balance: id must be not nil"}
	}
	return g.DB.GetPlayer(id)
}

// AnnounceTournament controlls announcing tournament
func (g Game) AnnounceTournament(id string, deposit int) error {
	if deposit <= 0 {
		return errors.Error{Code: errors.NegativeDepositError, Message: "announce: cannot create tournament with not positive deposite, id: " + id}
	}
	if id == "" {
		return errors.Error{Code: errors.NotFoundError, Message: "announce: id must be not nil"}
	}
	return g.DB.CreateTournament(id, deposit)
}

// JoinTournament controlls joining player to tournament
func (g Game) JoinTournament(tourID, playerID string) error {
	if tourID == "" {
		return errors.Error{Code: errors.NotFoundError, Message: "join tournament: tournament id must be not nil"}
	}
	if playerID == "" {
		return errors.Error{Code: errors.NotFoundError, Message: "join tournament: player id must be not nil"}
	}
	isOpen, err := g.DB.GetTournamentState(tourID)
	if err != nil {
		return err
	}
	if !isOpen {
		return errors.Error{Code: errors.ClosedTournamentError, Message: "join tournament: cannot join to closed tournament, tourID: " + tourID}
	}
	p, err := g.DB.GetParticipants(tourID)
	if err != nil {
		return err
	}
	for i := range p {
		if p[i] == playerID {
			return errors.Error{Code: errors.DuplicatedIDError, Message: "join tournament: cannot join to one tournament twice, playerID: " + playerID}
		}
	}
	return g.DB.UpdateTourAndPlayer(tourID, playerID)
}

// Results controls getting results from tournament
// If tournament is opened, it closes it
func (g Game) Results(tourID string) (entity.Winners, error) {
	if tourID == "" {
		return entity.Winners{}, errors.Error{Code: errors.NotFoundError, Message: "results: id must be not nil"}
	}
	isOpen, err := g.DB.GetTournamentState(tourID)
	if err != nil {
		return entity.Winners{}, err
	}
	if isOpen {
		err = g.DB.CloseTournament(tourID)
		if err != nil {
			return entity.Winners{}, err
		}
		winner, err := chooseWinner(g, tourID)
		if err != nil {
			return entity.Winners{}, err
		}
		err = g.DB.SetTournamentWinner(tourID, winner)
		if err != nil {
			return entity.Winners{}, err
		}
	}
	return g.DB.GetWinner(tourID)
}

func chooseWinner(g Game, tourID string) (entity.Winner, error) {
	p, err := g.DB.GetParticipants(tourID)
	if err != nil {
		return entity.Winner{}, err
	}
	if len(p) == 0 {
		return entity.Winner{}, errors.Error{Code: errors.NoneParticipantsError, Message: "cannot choose winner: tournament has no participants, id: " + tourID}
	}
	rand.Seed(time.Now().UnixNano())
	win, err := g.DB.GetPlayer(p[rand.Intn(len(p))])
	if err != nil {
		return entity.Winner{}, err
	}
	return entity.Winner{ID: win.ID, Points: win.Points}, nil
}

// DeletePlayer controls deleting player
func (g Game) DeletePlayer(playerID string) error {
	if playerID == "" {
		return errors.Error{Code: errors.NotFoundError, Message: "delete player: cannot delete player with empty id"}
	}
	return g.DB.DeletePlayer(playerID)
}

// DeleteTournament controls deleting tournament
func (g Game) DeleteTournament(tourID string) error {
	if tourID == "" {
		return errors.Error{Code: errors.NotFoundError, Message: "delete tournament: cannot delete tournament with empty id"}
	}
	return g.DB.DeleteTournament(tourID)
}

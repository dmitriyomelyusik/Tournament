// Package controller has some methods to controll, what data send to database queries
package controller

import (
	"github.com/Tournament/entity"
	"github.com/Tournament/errors"
)

// PlayerDB is an interface for database, that used to controll player activity methods
type PlayerDB interface {
	GetPlayer(string) (entity.Player, error)
	CreatePlayer(string, uint) error
	UpdatePlayer(string, int) error
}

// TourDB is an interface for database, that used to controll tournament activity methods
type TourDB interface {
	CreateTournament(id string, deposit uint) error
	GetDeposit(id string) (uint, error)
	GetParticipants(id string) ([]entity.Player, error)
	SetTournamentWinner(id string, player entity.Player) error
	UpdateParticipants(id string, player entity.Player) error
}

// Game methods controll game activity
type Game struct {
	PDB PlayerDB
	TDB TourDB
}

// Fund controlls funding player
func (g Game) Fund(id string, points int) error {
	if points < 0 {
		return errors.Error{Code: errors.NegativePointsNumberError, Message: "cannot fund negative number of points"}
	}
	_, err := g.PDB.GetPlayer(id)
	if err != nil {
		return g.PDB.CreatePlayer(id, uint(points))
	}
	return g.PDB.UpdatePlayer(id, points)
}

// Take controlls taking points
func (g Game) Take(id string, points int) error {
	if points < 0 {
		return errors.Error{Code: errors.NegativePointsNumberError, Message: "cannot take negative number of points"}
	}
	oldP, err := g.PDB.GetPlayer(id)
	if err != nil {
		return errors.Error{Code: errors.PlayerNotFoundError, Message: "cannot found player with id " + id}
	}
	err = g.PDB.UpdatePlayer(id, points*-1)
	if err != nil {
		return err
	}
	newP, err := g.PDB.GetPlayer(id)
	if err != nil {
		return err
	}
	if oldP.Points != newP.Points+uint(points) {
		err = g.PDB.UpdatePlayer(id, points)
		if err != nil {
			return errors.Error{Code: errors.TransactionError, Message: "failed to make backup of transaction, player id: " + id}
		}
		return errors.Error{Code: errors.TransactionError, Message: "failed to take player points, because it was used during operation, player id: " + id}
	}
	return nil
}

// Balance controlls getting actual player balance
func (g Game) Balance(id string) (entity.Player, error) {
	return g.PDB.GetPlayer(id)

}

// AnnounceTournament controlls announcing tournament
func (g Game) AnnounceTournament(id string, deposit int) error {
	if deposit <= 0 {
		return errors.Error{Code: errors.NegativeDepositError, Message: "cannot create tournament with non positive deposite, id: " + id}
	}
	return g.TDB.CreateTournament(id, uint(deposit))
}

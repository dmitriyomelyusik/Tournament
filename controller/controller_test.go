package controller

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dmitriyomelyusik/Tournament/entity"
	"github.com/dmitriyomelyusik/Tournament/errors"
)

var (
	g  Game
	db *MockDatabase
)

func TestMain(m *testing.M) {
	db = &MockDatabase{}
	g.DB = db
	code := m.Run()
	os.Exit(code)
}

func TestController_Fund(t *testing.T) {
	players := []entity.Player{
		{ID: "fund_ok", Points: 100},
		{ID: "fund_not_found", Points: 100},
		{ID: "fund_negative_points", Points: -100},
	}
	db.On("GetPlayer", players[0].ID).Return(players[0], nil)
	db.On("GetPlayer", players[1].ID).Return(entity.Player{}, errors.Error{Code: errors.NotFoundError})

	db.On("UpdatePlayer", players[0].ID, players[0].Points).Return(nil)

	db.On("CreatePlayer", players[1].ID, players[0].Points).Return(players[1], nil)
	tt := []struct {
		name           string
		playerID       string
		fund           int
		expectedPlayer entity.Player
		expectedError  error
	}{
		{
			name:           "fund: ok",
			playerID:       players[0].ID,
			fund:           players[0].Points,
			expectedPlayer: entity.Player{},
			expectedError:  nil,
		},
		{
			name:           "fund: create",
			playerID:       players[1].ID,
			fund:           players[1].Points,
			expectedPlayer: players[1],
			expectedError:  nil,
		},
		{
			name:           "fund: negative points number",
			playerID:       players[2].ID,
			fund:           players[2].Points,
			expectedPlayer: entity.Player{},
			expectedError:  errors.Error{Code: errors.NegativePointsNumberError, Message: "fund: cannot fund negative number of points"},
		},
		{
			name:           "fund: empty id",
			playerID:       "",
			fund:           0,
			expectedPlayer: entity.Player{},
			expectedError:  errors.Error{Code: errors.NotFoundError, Message: "fund: id must be not nil"},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			p, err := g.Fund(tc.playerID, tc.fund)
			assert.Equal(t, tc.expectedPlayer, p)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestController_Take(t *testing.T) {
	players := []entity.Player{
		{ID: "take_ok", Points: 100},
		{ID: "take_negative_points", Points: -100},
		{ID: "take_not_found", Points: 0},
	}
	db.On("UpdatePlayer", players[0].ID, -1*players[0].Points).Return(nil)
	db.On("UpdatePlayer", players[1].ID, 2*players[1].Points).Return(errors.Error{Code: errors.NegativePointsNumberError})
	db.On("UpdatePlayer", players[2].ID, -1*players[2].Points).Return(errors.Error{Code: errors.NotFoundError})
	tt := []struct {
		name          string
		playerID      string
		take          int
		expectedError error
	}{
		{
			name:          "take: ok",
			playerID:      players[0].ID,
			take:          players[0].Points,
			expectedError: nil,
		},
		{
			name:          "take: negative take",
			playerID:      players[1].ID,
			take:          players[1].Points,
			expectedError: errors.Error{Code: errors.NegativePointsNumberError, Message: "take: cannot take negative number of points"},
		},
		{
			name:          "take: empty id",
			playerID:      "",
			take:          0,
			expectedError: errors.Error{Code: errors.NotFoundError, Message: "take: id must be not nil"},
		},
		{
			name:          "take: more than can",
			playerID:      players[1].ID,
			take:          players[1].Points * -2,
			expectedError: errors.Error{Code: errors.NegativePointsNumberError, Message: "take: cannot take points, player doesn't have enough points"},
		},
		{
			name:          "take: not existing player",
			playerID:      players[2].ID,
			take:          players[2].Points,
			expectedError: errors.Error{Code: errors.NotFoundError},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := g.Take(tc.playerID, tc.take)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestController_Balance(t *testing.T) {
	players := []entity.Player{
		{ID: "balance_id", Points: 100},
	}
	db.On("GetPlayer", players[0].ID).Return(players[0], nil)
	tt := []struct {
		name           string
		playerID       string
		expectedPlayer entity.Player
		expectedError  error
	}{
		{
			name:           "balance: ok",
			playerID:       players[0].ID,
			expectedPlayer: players[0],
			expectedError:  nil,
		},
		{
			name:           "balance: empty id",
			playerID:       "",
			expectedPlayer: entity.Player{},
			expectedError:  errors.Error{Code: errors.NotFoundError, Message: "balance: id must be not nil"},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			p, err := g.Balance(tc.playerID)
			assert.Equal(t, tc.expectedPlayer, p)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestController_Announce(t *testing.T) {
	tournaments := []entity.Tournament{
		{ID: "announce_ok", Deposit: 100},
		{ID: "announce_negative_deposit", Deposit: -100},
	}
	db.On("CreateTournament", tournaments[0].ID, tournaments[0].Deposit).Return(nil)
	tt := []struct {
		name          string
		tourID        string
		deposit       int
		expectedError error
	}{
		{
			name:          "announce: ok",
			tourID:        tournaments[0].ID,
			deposit:       tournaments[0].Deposit,
			expectedError: nil,
		},
		{
			name:          "announce: negative deposit",
			tourID:        tournaments[1].ID,
			deposit:       tournaments[1].Deposit,
			expectedError: errors.Error{Code: errors.NegativeDepositError, Message: "announce: cannot create tournament with not positive deposite, id: " + tournaments[1].ID},
		},
		{
			name:          "announce: empty id",
			tourID:        "",
			deposit:       1,
			expectedError: errors.Error{Code: errors.NotFoundError, Message: "announce: id must be not nil"},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := g.AnnounceTournament(tc.tourID, tc.deposit)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestController_Join(t *testing.T) {
	players := []entity.Player{
		{ID: "join_ok", Points: 100},
		{ID: "join_duplicate", Points: 200},
	}
	tournaments := []entity.Tournament{
		{ID: "join_ok", Deposit: 50, IsOpen: true},
		{ID: "join_not_found", Deposit: 50},
		{ID: "join_closed_tournament", Deposit: 15, IsOpen: false},
		{ID: "join_getparticipants_error", Deposit: 20, IsOpen: true},
		{ID: "join_duplicate", Deposit: 33, IsOpen: true},
	}
	db.On("GetTournamentState", tournaments[0].ID).Return(tournaments[0].IsOpen, nil)
	db.On("GetTournamentState", tournaments[1].ID).Return(false, errors.Error{Code: errors.NotFoundError})
	db.On("GetTournamentState", tournaments[2].ID).Return(tournaments[2].IsOpen, nil)
	db.On("GetTournamentState", tournaments[3].ID).Return(tournaments[3].IsOpen, nil)
	db.On("GetTournamentState", tournaments[4].ID).Return(tournaments[4].IsOpen, nil)

	db.On("GetParticipants", tournaments[0].ID).Return(nil, nil)
	db.On("GetParticipants", tournaments[3].ID).Return(nil, errors.Error{Code: errors.NotFoundError})
	db.On("GetParticipants", tournaments[4].ID).Return([]string{players[1].ID}, nil)

	db.On("UpdateTourAndPlayer", tournaments[0].ID, players[0].ID).Return(nil)
	tt := []struct {
		name          string
		tourID        string
		playerID      string
		expectedError error
	}{
		{
			name:          "join: ok",
			tourID:        tournaments[0].ID,
			playerID:      players[0].ID,
			expectedError: nil,
		},
		{
			name:          "join: empty tournament id",
			tourID:        "",
			playerID:      players[0].ID,
			expectedError: errors.Error{Code: errors.NotFoundError, Message: "join tournament: tournament id must be not nil"},
		},
		{
			name:          "join: empty player id",
			tourID:        tournaments[0].ID,
			playerID:      "",
			expectedError: errors.Error{Code: errors.NotFoundError, Message: "join tournament: player id must be not nil"},
		},
		{
			name:          "join: not found tournament",
			tourID:        tournaments[1].ID,
			playerID:      players[0].ID,
			expectedError: errors.Error{Code: errors.NotFoundError},
		},
		{
			name:          "join: closed tournament",
			tourID:        tournaments[2].ID,
			playerID:      players[0].ID,
			expectedError: errors.Error{Code: errors.ClosedTournamentError, Message: "join tournament: cannot join to closed tournament, tourID: " + tournaments[2].ID},
		},
		{
			name:          "join: get participants error",
			tourID:        tournaments[3].ID,
			playerID:      players[0].ID,
			expectedError: errors.Error{Code: errors.NotFoundError},
		},
		{
			name:          "join: duplicated player",
			tourID:        tournaments[4].ID,
			playerID:      players[1].ID,
			expectedError: errors.Error{Code: errors.DuplicatedIDError, Message: "join tournament: cannot join to one tournament twice, playerID: " + players[1].ID},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := g.JoinTournament(tc.tourID, tc.playerID)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

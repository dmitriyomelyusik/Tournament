package handlers

import (
	"encoding/json"
	e "errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dmitriyomelyusik/Tournament/entity"
	"github.com/dmitriyomelyusik/Tournament/errors"
)

var (
	ts         *httptest.Server
	controller *mockCtlr
)

func TestMain(m *testing.M) {
	controller = new(mockCtlr)
	r := NewRouter(Server{Controller: controller})
	ts = httptest.NewServer(r)
	defer ts.Close()
	code := m.Run()
	os.Exit(code)
}

func TestHandlers_FundHandler(t *testing.T) {
	players := []entity.Player{
		{ID: "fundplayer_test1", Points: 200},
		{ID: "fundplayer_test2", Points: 200},
		{ID: "fundplayer_test3", Points: 200},
	}
	client := http.Client{}
	tt := []struct {
		name                   string
		player                 entity.Player
		fund                   interface{}
		expectedNegativeStatus int
		expectedError          error
		expectedPlayer         entity.Player
		expectedStatus         int
		returnedPlayerOne      entity.Player
		returnedPlayerTwo      entity.Player
		negativeError          error
	}{
		{
			name:                   "ok_createfund0",
			player:                 players[0],
			fund:                   0,
			expectedError:          nil,
			expectedNegativeStatus: http.StatusOK,
			expectedPlayer:         players[0],
			expectedStatus:         http.StatusCreated,
			returnedPlayerOne:      players[0],
			returnedPlayerTwo:      entity.Player{},
			negativeError:          nil,
		},
		{
			name:                   "ok_createfund1",
			player:                 players[1],
			fund:                   0,
			expectedError:          nil,
			expectedNegativeStatus: http.StatusOK,
			expectedPlayer:         players[1],
			expectedStatus:         http.StatusCreated,
			returnedPlayerOne:      players[1],
			returnedPlayerTwo:      entity.Player{},
			negativeError:          nil,
		},
		{
			name:                   "ok_createfund2",
			player:                 players[2],
			fund:                   0,
			expectedError:          nil,
			expectedNegativeStatus: http.StatusOK,
			expectedPlayer:         players[2],
			expectedStatus:         http.StatusCreated,
			returnedPlayerOne:      players[2],
			returnedPlayerTwo:      entity.Player{},
			negativeError:          nil,
		},
		{
			name:                   "ok_fund0",
			player:                 players[0],
			fund:                   -100,
			expectedError:          nil,
			expectedNegativeStatus: http.StatusNotFound,
			expectedPlayer:         entity.Player{ID: players[0].ID, Points: players[0].Points * 2},
			expectedStatus:         http.StatusOK,
			returnedPlayerOne:      entity.Player{},
			returnedPlayerTwo:      entity.Player{},
			negativeError:          errors.Error{Code: errors.NegativePointsNumberError},
		},
		{
			name:                   "ok_fund1",
			player:                 players[1],
			fund:                   -100,
			expectedError:          nil,
			expectedNegativeStatus: http.StatusNotFound,
			expectedPlayer:         entity.Player{ID: players[1].ID, Points: players[1].Points * 2},
			expectedStatus:         http.StatusOK,
			returnedPlayerOne:      entity.Player{},
			returnedPlayerTwo:      entity.Player{},
			negativeError:          errors.Error{Code: errors.NegativePointsNumberError},
		},
		{
			name:                   "ok_fund2",
			player:                 players[2],
			fund:                   -100,
			expectedError:          nil,
			expectedNegativeStatus: http.StatusNotFound,
			expectedPlayer:         entity.Player{ID: players[2].ID, Points: players[2].Points * 2},
			expectedStatus:         http.StatusOK,
			returnedPlayerOne:      entity.Player{},
			returnedPlayerTwo:      entity.Player{},
			negativeError:          errors.Error{Code: errors.NegativePointsNumberError},
		},
		{
			name:                   "notnumber_fund",
			player:                 players[0],
			fund:                   "abc",
			expectedError:          nil,
			expectedNegativeStatus: http.StatusNotFound,
			expectedPlayer:         entity.Player{},
			expectedStatus:         http.StatusOK,
			returnedPlayerOne:      entity.Player{},
			returnedPlayerTwo:      entity.Player{},
			negativeError:          errors.Error{Code: errors.NotFoundError},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			controller.On("Fund", tc.player.ID, tc.player.Points).
				Return(tc.returnedPlayerOne, tc.expectedError).Once()
			req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%v/fund?playerId=%v&points=%v", ts.URL, tc.player.ID, tc.player.Points), nil)
			assert.Equal(t, tc.expectedError, err)
			res, err := client.Do(req)
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedStatus, res.StatusCode)

			controller.On("Balance", tc.player.ID).
				Return(tc.expectedPlayer, tc.expectedError).Once()
			req, err = http.NewRequest(http.MethodGet, fmt.Sprintf("%v/balance?playerId=%v", ts.URL, tc.player.ID), nil)
			assert.Equal(t, tc.expectedError, err)
			res, err = client.Do(req)
			assert.Equal(t, tc.expectedError, err)
			decoder := json.NewDecoder(res.Body)
			var player entity.Player
			err = decoder.Decode(&player)
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedPlayer, player)

			controller.On("Fund", tc.player.ID, tc.fund).
				Return(tc.returnedPlayerTwo, tc.negativeError).Once()
			req, err = http.NewRequest(http.MethodPut, fmt.Sprintf("%v/fund?playerId=%v&points=%v", ts.URL, tc.player.ID, tc.fund), nil)
			assert.Equal(t, tc.expectedError, err)
			res, err = client.Do(req)
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedNegativeStatus, res.StatusCode)
		})
	}
}

func TestHandlers_TakeHandler(t *testing.T) {
	players := []entity.Player{
		{ID: "takeplayer_test1", Points: 200},
		{ID: "takeplayer_test2", Points: 200},
		{ID: "takeplayer_test3", Points: 200},
	}
	take := []int{150, 150, 150}
	client := http.Client{}
	tt := []struct {
		name           string
		player         entity.Player
		take           interface{}
		takeError      error
		balancePlayer  entity.Player
		balanceError   error
		expectedError  error
		expectedPlayer entity.Player
		expectedStatus int
	}{
		{
			name:           "ok_take0",
			player:         players[0],
			take:           take[0],
			takeError:      nil,
			balancePlayer:  entity.Player{ID: players[0].ID, Points: players[0].Points - take[0]},
			balanceError:   nil,
			expectedError:  nil,
			expectedPlayer: entity.Player{ID: players[0].ID, Points: players[0].Points - take[0]},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "ok_take1",
			player:         players[1],
			take:           take[1],
			takeError:      nil,
			balancePlayer:  entity.Player{ID: players[1].ID, Points: players[1].Points - take[1]},
			balanceError:   nil,
			expectedError:  nil,
			expectedPlayer: entity.Player{ID: players[1].ID, Points: players[1].Points - take[1]},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "ok_take2",
			player:         players[2],
			take:           take[2],
			takeError:      nil,
			balancePlayer:  entity.Player{ID: players[2].ID, Points: players[2].Points - take[2]},
			balanceError:   nil,
			expectedError:  nil,
			expectedPlayer: entity.Player{ID: players[2].ID, Points: players[2].Points - take[2]},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "err_take0",
			player:         players[0],
			take:           take[0],
			takeError:      errors.Error{Code: errors.NegativePointsNumberError},
			balancePlayer:  entity.Player{ID: players[0].ID, Points: players[0].Points - take[0]},
			balanceError:   nil,
			expectedError:  nil,
			expectedPlayer: entity.Player{ID: players[0].ID, Points: players[0].Points - take[0]},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "err_take1",
			player:         players[1],
			take:           take[1],
			takeError:      errors.Error{Code: errors.NegativePointsNumberError},
			balancePlayer:  entity.Player{ID: players[1].ID, Points: players[1].Points - take[1]},
			balanceError:   nil,
			expectedError:  nil,
			expectedPlayer: entity.Player{ID: players[1].ID, Points: players[1].Points - take[1]},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "err_take2",
			player:         players[2],
			take:           take[2],
			takeError:      errors.Error{Code: errors.NegativePointsNumberError},
			balancePlayer:  entity.Player{ID: players[2].ID, Points: players[2].Points - take[2]},
			balanceError:   nil,
			expectedError:  nil,
			expectedPlayer: entity.Player{ID: players[2].ID, Points: players[2].Points - take[2]},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "notnumber_take",
			player:         players[0],
			take:           "abc",
			takeError:      errors.Error{Code: errors.NotNumberError},
			balancePlayer:  players[0],
			balanceError:   nil,
			expectedError:  nil,
			expectedPlayer: players[0],
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			controller.On("Take", tc.player.ID, tc.take).
				Return(tc.takeError).Once()
			req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%v/take?playerId=%v&points=%v", ts.URL, tc.player.ID, tc.take), nil)
			assert.Equal(t, tc.expectedError, err)
			res, err := client.Do(req)
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedStatus, res.StatusCode)
			controller.On("Balance", tc.player.ID).
				Return(tc.balancePlayer, tc.balanceError).Once()
			req, err = http.NewRequest(http.MethodGet, fmt.Sprintf("%v/balance?playerId=%v", ts.URL, tc.player.ID), nil)
			assert.Equal(t, tc.expectedError, err)
			res, err = client.Do(req)
			assert.Equal(t, tc.expectedError, err)
			decoder := json.NewDecoder(res.Body)
			var player entity.Player
			err = decoder.Decode(&player)
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedPlayer, player)
		})
	}
}

func TestHandlers_BalanceHandler(t *testing.T) {
	players := []entity.Player{
		{ID: "balanceplayer_test1", Points: 200},
		{ID: "balanceplayer_test2", Points: 200},
		{ID: "balanceplayer_test3", Points: 200},
	}
	take := []int{150, 150, 150}
	fund := []int{100, 100, 100}
	client := http.Client{}
	tt := []struct {
		name           string
		player         entity.Player
		take           int
		fund           int
		takeError      error
		fundPlayer     entity.Player
		fundError      error
		balancePlayer  entity.Player
		balanceError   error
		expectedError  error
		expectedPlayer entity.Player
		expectedStatus int
	}{
		{
			name:           "balance_fund0",
			player:         players[0],
			take:           0,
			fund:           fund[0],
			takeError:      nil,
			fundPlayer:     entity.Player{},
			fundError:      nil,
			balancePlayer:  entity.Player{ID: players[0].ID, Points: players[0].Points + fund[0]},
			balanceError:   nil,
			expectedError:  nil,
			expectedPlayer: entity.Player{ID: players[0].ID, Points: players[0].Points + fund[0]},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "balance_fund1",
			player:         players[1],
			take:           0,
			fund:           fund[1],
			takeError:      nil,
			fundPlayer:     entity.Player{},
			fundError:      nil,
			balancePlayer:  entity.Player{ID: players[1].ID, Points: players[1].Points + fund[1]},
			balanceError:   nil,
			expectedError:  nil,
			expectedPlayer: entity.Player{ID: players[1].ID, Points: players[1].Points + fund[1]},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "balance_fund2",
			player:         players[2],
			take:           0,
			fund:           fund[2],
			takeError:      nil,
			fundPlayer:     entity.Player{},
			fundError:      nil,
			balancePlayer:  entity.Player{ID: players[2].ID, Points: players[2].Points + fund[2]},
			balanceError:   nil,
			expectedError:  nil,
			expectedPlayer: entity.Player{ID: players[2].ID, Points: players[2].Points + fund[2]},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "balance_take0",
			player:         players[0],
			take:           take[0],
			fund:           0,
			takeError:      nil,
			fundPlayer:     entity.Player{},
			fundError:      nil,
			balancePlayer:  entity.Player{ID: players[0].ID, Points: players[0].Points + fund[0] - take[0]},
			balanceError:   nil,
			expectedError:  nil,
			expectedPlayer: entity.Player{ID: players[0].ID, Points: players[0].Points + fund[0] - take[0]},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "balance_take1",
			player:         players[1],
			take:           take[1],
			fund:           fund[1],
			takeError:      nil,
			fundPlayer:     entity.Player{},
			fundError:      nil,
			balancePlayer:  entity.Player{ID: players[1].ID, Points: players[1].Points + fund[1] + fund[1] - take[1]},
			balanceError:   nil,
			expectedError:  nil,
			expectedPlayer: entity.Player{ID: players[1].ID, Points: players[1].Points + fund[1] + fund[1] - take[1]},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "balance_take2",
			player:         players[2],
			take:           take[2],
			fund:           0,
			takeError:      nil,
			fundPlayer:     entity.Player{},
			fundError:      nil,
			balancePlayer:  entity.Player{ID: players[2].ID, Points: players[2].Points + fund[2] - take[2]},
			balanceError:   nil,
			expectedError:  nil,
			expectedPlayer: entity.Player{ID: players[2].ID, Points: players[2].Points + fund[2] - take[2]},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "balance_fake",
			player:         entity.Player{},
			take:           0,
			fund:           0,
			takeError:      errors.Error{Code: errors.NotFoundError},
			fundPlayer:     entity.Player{},
			fundError:      errors.Error{Code: errors.NotFoundError},
			balancePlayer:  entity.Player{},
			balanceError:   errors.Error{Code: errors.NotFoundError},
			expectedError:  nil,
			expectedPlayer: entity.Player{},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			controller.On("Take", tc.player.ID, tc.take).
				Return(tc.takeError).Once()
			req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%v/take?playerId=%v&points=%v", ts.URL, tc.player.ID, tc.take), nil)
			assert.Equal(t, tc.expectedError, err)
			res, err := client.Do(req)
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedStatus, res.StatusCode)
			controller.On("Fund", tc.player.ID, tc.fund).
				Return(tc.fundPlayer, tc.fundError).Once()
			req, err = http.NewRequest(http.MethodPut, fmt.Sprintf("%v/fund?playerId=%v&points=%v", ts.URL, tc.player.ID, tc.fund), nil)
			assert.Equal(t, tc.expectedError, err)
			res, err = client.Do(req)
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedStatus, res.StatusCode)
			controller.On("Balance", tc.player.ID).
				Return(tc.balancePlayer, tc.balanceError).Once()
			req, err = http.NewRequest(http.MethodGet, fmt.Sprintf("%v/balance?playerId=%v", ts.URL, tc.player.ID), nil)
			assert.Equal(t, tc.expectedError, err)
			res, err = client.Do(req)
			assert.Equal(t, tc.expectedError, err)
			decoder := json.NewDecoder(res.Body)
			var player entity.Player
			err = decoder.Decode(&player)
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedPlayer, player)
		})
	}
}

func TestHandlers_AnnounceHandler(t *testing.T) {
	tournaments := []entity.Tournament{
		{ID: "announce_tour1", Deposit: 100},
		{ID: "announce_tour2", Deposit: 100},
		{ID: "announce_tour3", Deposit: 100},
	}
	client := http.Client{}
	tt := []struct {
		name           string
		tournamentID   string
		deposit        interface{}
		expectedError  error
		expectedStatus int
		announceError  error
	}{
		{
			name:           "announce_ok0",
			tournamentID:   tournaments[0].ID,
			deposit:        tournaments[0].Deposit,
			expectedError:  nil,
			expectedStatus: http.StatusOK,
			announceError:  nil,
		},
		{
			name:           "announce_ok1",
			tournamentID:   tournaments[1].ID,
			deposit:        tournaments[1].Deposit,
			expectedError:  nil,
			expectedStatus: http.StatusOK,
			announceError:  nil,
		},
		{
			name:           "announce_ok2",
			tournamentID:   tournaments[2].ID,
			deposit:        tournaments[2].Deposit,
			expectedError:  nil,
			expectedStatus: http.StatusOK,
			announceError:  nil,
		},
		{
			name:           "announce_duplicated0",
			tournamentID:   tournaments[0].ID,
			deposit:        tournaments[0].Deposit,
			expectedError:  nil,
			expectedStatus: http.StatusNotFound,
			announceError:  errors.Error{Code: errors.DuplicatedIDError},
		},
		{
			name:           "announce_duplicated1",
			tournamentID:   tournaments[1].ID,
			deposit:        tournaments[1].Deposit,
			expectedError:  nil,
			expectedStatus: http.StatusNotFound,
			announceError:  errors.Error{Code: errors.DuplicatedIDError},
		},
		{
			name:           "announce_duplicated2",
			tournamentID:   tournaments[2].ID,
			deposit:        tournaments[2].Deposit,
			expectedError:  nil,
			expectedStatus: http.StatusNotFound,
			announceError:  errors.Error{Code: errors.DuplicatedIDError},
		},
		{
			name:           "announce_notNumberDeposit",
			tournamentID:   "fakeID:TROLOLO",
			deposit:        "mayoneZ",
			expectedError:  nil,
			expectedStatus: http.StatusNotFound,
			announceError:  errors.Error{Code: errors.NotNumberError},
		},
		{
			name:           "announce_brokenDB",
			tournamentID:   "itisnotnecessary",
			deposit:        9999,
			expectedError:  nil,
			expectedStatus: http.StatusInternalServerError,
			announceError:  e.New("UnexpectedError"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			controller.On("AnnounceTournament", tc.tournamentID, tc.deposit).
				Return(tc.announceError).Once()
			req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%v/announceTournament?tournamentId=%v&deposit=%v", ts.URL, tc.tournamentID, tc.deposit), nil)
			assert.Equal(t, tc.expectedError, err)
			res, err := client.Do(req)
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedStatus, res.StatusCode)
		})
	}
}

func TestHandlers_JoinHandler(t *testing.T) {
	tournaments := []entity.Tournament{
		{ID: "join_tour1", Deposit: 100},
		{ID: "join_tour2", Deposit: 100},
		{ID: "join_tour3", Deposit: 100},
	}
	players := []entity.Player{
		{ID: "join_player1", Points: 350},
		{ID: "join_player2", Points: 250},
		{ID: "join_player3", Points: 150},
	}
	client := http.Client{}
	tt := []struct {
		name               string
		tournament         entity.Tournament
		participants       []entity.Player
		expectedPlayers    []entity.Player
		expectedJoinStatus []int
		expectedJoinError  []error
		expectedError      error
		joinErrors         []error
		balancePlayers     []entity.Player
		balanceErrors      []error
	}{
		{
			name:         "join_ok0",
			tournament:   tournaments[0],
			participants: players,
			expectedPlayers: []entity.Player{
				{ID: players[0].ID, Points: players[0].Points - tournaments[0].Deposit},
				{ID: players[1].ID, Points: players[1].Points - tournaments[0].Deposit},
				{ID: players[2].ID, Points: players[2].Points - tournaments[0].Deposit}},
			expectedJoinStatus: []int{http.StatusOK, http.StatusOK, http.StatusOK},
			expectedJoinError:  []error{nil, nil, nil},
			expectedError:      nil,
			joinErrors:         []error{nil, nil, nil},
			balancePlayers: []entity.Player{
				{ID: players[0].ID, Points: players[0].Points - tournaments[0].Deposit},
				{ID: players[1].ID, Points: players[1].Points - tournaments[0].Deposit},
				{ID: players[2].ID, Points: players[2].Points - tournaments[0].Deposit}},
			balanceErrors: []error{nil, nil, nil},
		},
		{
			name:         "join_ok1",
			tournament:   tournaments[1],
			participants: players,
			expectedPlayers: []entity.Player{
				{ID: players[0].ID, Points: players[0].Points - tournaments[1].Deposit - tournaments[0].Deposit},
				{ID: players[1].ID, Points: players[1].Points - tournaments[1].Deposit - tournaments[0].Deposit},
				{ID: players[2].ID, Points: players[2].Points - tournaments[0].Deposit}},
			expectedJoinError:  []error{nil, nil, errors.Error{Code: errors.NegativePointsNumberError, Message: "update player: cannot update points numbers, dif -100"}},
			expectedJoinStatus: []int{http.StatusOK, http.StatusOK, http.StatusNotFound},
			expectedError:      nil,
			joinErrors:         []error{nil, nil, errors.Error{Code: errors.NegativePointsNumberError, Message: "update player: cannot update points numbers, dif -100"}},
			balancePlayers: []entity.Player{
				{ID: players[0].ID, Points: players[0].Points - tournaments[1].Deposit - tournaments[0].Deposit},
				{ID: players[1].ID, Points: players[1].Points - tournaments[1].Deposit - tournaments[0].Deposit},
				{ID: players[2].ID, Points: players[2].Points - tournaments[0].Deposit}},
			balanceErrors: []error{nil, nil, nil},
		},
		{
			name:         "join_ok2",
			tournament:   tournaments[2],
			participants: players,
			expectedPlayers: []entity.Player{
				{ID: players[0].ID, Points: players[0].Points - tournaments[2].Deposit - tournaments[1].Deposit - tournaments[0].Deposit},
				{ID: players[1].ID, Points: players[1].Points - tournaments[1].Deposit - tournaments[0].Deposit},
				{ID: players[2].ID, Points: players[2].Points - tournaments[0].Deposit}},
			expectedJoinError:  []error{nil, errors.Error{Code: errors.NegativePointsNumberError, Message: "update player: cannot update points numbers, dif -100"}, errors.Error{Code: errors.NegativePointsNumberError, Message: "update player: cannot update points numbers, dif -100"}},
			expectedJoinStatus: []int{http.StatusOK, http.StatusNotFound, http.StatusNotFound},
			expectedError:      nil,
			joinErrors:         []error{nil, errors.Error{Code: errors.NegativePointsNumberError, Message: "update player: cannot update points numbers, dif -100"}, errors.Error{Code: errors.NegativePointsNumberError, Message: "update player: cannot update points numbers, dif -100"}},
			balancePlayers: []entity.Player{
				{ID: players[0].ID, Points: players[0].Points - tournaments[2].Deposit - tournaments[1].Deposit - tournaments[0].Deposit},
				{ID: players[1].ID, Points: players[1].Points - tournaments[1].Deposit - tournaments[0].Deposit},
				{ID: players[2].ID, Points: players[2].Points - tournaments[0].Deposit}},
			balanceErrors: []error{nil, nil, nil},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			for i := range tc.participants {
				controller.On("JoinTournament", tc.tournament.ID, tc.participants[i].ID).
					Return(tc.joinErrors[i]).Once()
				req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%v/joinTournament?tournamentId=%v&playerId=%v", ts.URL, tc.tournament.ID, tc.participants[i].ID), nil)
				assert.Equal(t, tc.expectedError, err)
				res, err := client.Do(req)
				assert.Equal(t, tc.expectedError, err)
				assert.Equal(t, tc.expectedJoinStatus[i], res.StatusCode)
				if tc.expectedJoinError[i] != nil {
					var joinErr errors.Error
					decoder := json.NewDecoder(res.Body)
					err = decoder.Decode(&joinErr)
					assert.Equal(t, tc.expectedError, err)
					assert.Equal(t, tc.expectedJoinError[i], joinErr)
				}

				controller.On("Balance", tc.participants[i].ID).
					Return(tc.balancePlayers[i], tc.balanceErrors[i]).Once()
				req, err = http.NewRequest(http.MethodGet, fmt.Sprintf("%v/balance?playerId=%v", ts.URL, tc.participants[i].ID), nil)
				assert.Equal(t, tc.expectedError, err)
				res, err = client.Do(req)
				assert.Equal(t, tc.expectedError, err)
				var player entity.Player
				decoder := json.NewDecoder(res.Body)
				err = decoder.Decode(&player)
				assert.Equal(t, tc.expectedError, err)
				assert.Equal(t, tc.expectedPlayers[i], player)
			}
		})
	}
}

func TestHandlers_ResultHandler(t *testing.T) {
	tournaments := []entity.Tournament{
		{ID: "result_tour1", Deposit: 100},
		{ID: "result_tour2", Deposit: 200},
		{ID: "result_tour3", Deposit: 300},
	}
	players := []entity.Player{
		{ID: "result_player1", Points: 300},
		{ID: "result_player2", Points: 200},
		{ID: "result_player3", Points: 100},
	}
	client := http.Client{}
	tt := []struct {
		name               string
		id                 string
		expectedError      error
		expectedStatus     int
		expectedPrize      int
		expectedJoinStatus int
		expectedJoinError  error
		resultWinners      entity.Winners
		resultError        error
		joinError          error
	}{
		{
			name:               "result_ok0",
			id:                 tournaments[0].ID,
			expectedError:      nil,
			expectedStatus:     http.StatusOK,
			expectedPrize:      tournaments[0].Deposit * 3,
			expectedJoinStatus: http.StatusNotFound,
			expectedJoinError:  errors.Error{Code: errors.ClosedTournamentError, Message: "join tournament: cannot join to closed tournament, tourID: " + tournaments[0].ID},
			resultWinners:      entity.Winners{Winners: []entity.Winner{entity.Winner{ID: players[0].ID, Points: players[0].Points + tournaments[0].Deposit*3, Prize: tournaments[0].Deposit * 3}}},
			resultError:        nil,
			joinError:          errors.Error{Code: errors.ClosedTournamentError, Message: "join tournament: cannot join to closed tournament, tourID: " + tournaments[0].ID},
		},
		{
			name:               "result_ok1",
			id:                 tournaments[1].ID,
			expectedError:      nil,
			expectedStatus:     http.StatusOK,
			expectedPrize:      tournaments[1].Deposit * 2,
			expectedJoinStatus: http.StatusNotFound,
			expectedJoinError:  errors.Error{Code: errors.ClosedTournamentError, Message: "join tournament: cannot join to closed tournament, tourID: " + tournaments[1].ID},
			resultWinners:      entity.Winners{Winners: []entity.Winner{entity.Winner{ID: players[2].ID, Points: players[2].Points + tournaments[1].Deposit*3, Prize: tournaments[1].Deposit * 2}}},
			resultError:        nil,
			joinError:          errors.Error{Code: errors.ClosedTournamentError, Message: "join tournament: cannot join to closed tournament, tourID: " + tournaments[1].ID},
		},
		{
			name:               "result_ok2",
			id:                 tournaments[2].ID,
			expectedError:      nil,
			expectedStatus:     http.StatusOK,
			expectedPrize:      tournaments[2].Deposit,
			expectedJoinStatus: http.StatusNotFound,
			expectedJoinError:  errors.Error{Code: errors.ClosedTournamentError, Message: "join tournament: cannot join to closed tournament, tourID: " + tournaments[2].ID},
			resultWinners:      entity.Winners{Winners: []entity.Winner{entity.Winner{ID: players[2].ID, Points: players[2].Points + tournaments[2].Deposit*3, Prize: tournaments[2].Deposit * 1}}},
			resultError:        nil,
			joinError:          errors.Error{Code: errors.ClosedTournamentError, Message: "join tournament: cannot join to closed tournament, tourID: " + tournaments[2].ID},
		},
		{
			name:               "result_fake",
			id:                 "result_fake",
			expectedError:      nil,
			expectedStatus:     http.StatusNotFound,
			expectedPrize:      0,
			expectedJoinStatus: http.StatusNotFound,
			expectedJoinError:  errors.Error{Code: errors.NotFoundError, Message: "get state: cannot get tournament state from not existing tournament, id: announce_fake"},
			resultWinners:      entity.Winners{},
			resultError:        errors.Error{Code: errors.NotFoundError, Message: "get state: cannot get tournament state from not existing tournament, id: announce_fake"},
			joinError:          errors.Error{Code: errors.NotFoundError, Message: "get state: cannot get tournament state from not existing tournament, id: announce_fake"},
		},
		{
			name:               "result_emptyTour",
			id:                 tournaments[1].ID,
			expectedError:      nil,
			expectedStatus:     http.StatusOK,
			expectedPrize:      0,
			expectedJoinStatus: http.StatusNotFound,
			expectedJoinError:  errors.Error{Code: errors.ClosedTournamentError, Message: "join tournament: cannot join to closed tournament, tourID: " + tournaments[1].ID},
			resultWinners:      entity.Winners{},
			resultError:        errors.Error{Code: errors.NoneParticipantsError},
			joinError:          errors.Error{Code: errors.ClosedTournamentError, Message: "join tournament: cannot join to closed tournament, tourID: " + tournaments[1].ID},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			controller.On("Results", tc.id).
				Return(tc.resultWinners, tc.resultError).Once()
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%v/resultTournament?tournamentId=%v", ts.URL, tc.id), nil)
			assert.Equal(t, tc.expectedError, err)
			res, err := client.Do(req)
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedStatus, res.StatusCode)
			decoder := json.NewDecoder(res.Body)
			var winner entity.Winners
			err = decoder.Decode(&winner)
			assert.Equal(t, tc.expectedError, err)
			if len(winner.Winners) > 0 {
				assert.Equal(t, tc.expectedPrize, winner.Winners[0].Prize)
			}
			controller.On("JoinTournament", tc.id, players[0].ID).
				Return(tc.joinError).Once()
			req, err = http.NewRequest(http.MethodGet, fmt.Sprintf("%v/joinTournament?tournamentId=%v&playerId=%v", ts.URL, tc.id, players[0].ID), nil)
			assert.Equal(t, tc.expectedError, err)
			res, err = client.Do(req)
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedJoinStatus, res.StatusCode)
			decoder = json.NewDecoder(res.Body)
			var joinErr errors.Error
			err = decoder.Decode(&joinErr)
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedJoinError, joinErr)
		})
	}
}

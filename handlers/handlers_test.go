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
		{ID: "fund_ok", Points: 200},
		{ID: "fund_negative_points_number", Points: -100},
	}
	controller.On("Fund", players[0].ID, players[0].Points).Return(players[0], nil)
	controller.On("Fund", players[1].ID, players[1].Points).Return(entity.Player{}, errors.Error{Code: errors.NegativePointsNumberError})
	client := http.Client{}
	tt := []struct {
		name           string
		playerID       string
		fund           interface{}
		err            error
		expectedError  error
		expectedPlayer entity.Player
		expectedStatus int
	}{
		{
			name:           "fund: ok",
			playerID:       players[0].ID,
			fund:           players[0].Points,
			expectedPlayer: players[0],
			expectedStatus: http.StatusCreated,
			expectedError:  errors.Error{},
		},
		{
			name:           "fund: incorrect points format",
			playerID:       players[0].ID,
			fund:           "incorrect_format",
			expectedPlayer: entity.Player{},
			expectedStatus: http.StatusNotFound,
			expectedError:  errors.Error{Code: errors.NotNumberError, Message: "cannot fund player, points is not number: incorrect_format", Info: "strconv.Atoi: parsing \"incorrect_format\": invalid syntax"},
		},
		{
			name:           "fund: negative points number",
			playerID:       players[1].ID,
			fund:           players[1].Points,
			expectedPlayer: entity.Player{},
			expectedStatus: http.StatusNotFound,
			expectedError:  errors.Error{Code: errors.NegativePointsNumberError},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%v/fund?playerId=%v&points=%v", ts.URL, tc.playerID, tc.fund), nil)
			assert.Equal(t, tc.err, err)
			res, err := client.Do(req)
			assert.Equal(t, tc.err, err)
			assert.Equal(t, tc.expectedStatus, res.StatusCode)
			decoder := json.NewDecoder(res.Body)
			var expErr errors.Error
			err = decoder.Decode(&expErr)
			assert.Equal(t, tc.err, err)
			assert.Equal(t, tc.expectedError, expErr)
		})
	}
}

func TestHandlers_TakeHandler(t *testing.T) {
	players := []entity.Player{
		{ID: "take_ok", Points: 200},
		{ID: "take_unexpected", Points: 200},
		{ID: "take_not_found", Points: 200},
	}
	controller.On("Take", players[0].ID, players[0].Points).Return(nil)
	controller.On("Take", players[1].ID, players[1].Points).Return(errors.Error{Code: errors.UnexpectedError})
	controller.On("Take", players[2].ID, players[2].Points).Return(errors.Error{Code: errors.NotFoundError})
	client := http.Client{}
	tt := []struct {
		name           string
		playerID       string
		take           interface{}
		err            error
		expectedError  errors.Error
		expectedStatus int
	}{
		{
			name:           "take: ok",
			playerID:       players[0].ID,
			take:           players[0].Points,
			expectedError:  errors.Error{},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "take: incorrect format",
			playerID:       players[0].ID,
			take:           "incorrect_format",
			expectedError:  errors.Error{Code: errors.NotNumberError, Message: "cannot take points, points is not number: incorrect_format", Info: "strconv.Atoi: parsing \"incorrect_format\": invalid syntax"},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "take: unexpected error",
			playerID:       players[1].ID,
			take:           players[1].Points,
			expectedError:  errors.Error{Code: errors.UnexpectedError},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "take: unexisting player",
			playerID:       players[2].ID,
			take:           players[2].Points,
			expectedError:  errors.Error{Code: errors.NotFoundError},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%v/take?playerId=%v&points=%v", ts.URL, tc.playerID, tc.take), nil)
			assert.Equal(t, tc.err, err)
			res, err := client.Do(req)
			assert.Equal(t, tc.err, err)
			assert.Equal(t, tc.expectedStatus, res.StatusCode)
			if tc.expectedStatus != http.StatusOK {
				decoder := json.NewDecoder(res.Body)
				var expErr errors.Error
				err = decoder.Decode(&expErr)
				assert.Equal(t, tc.err, err)
				assert.Equal(t, tc.expectedError, expErr)
			}
		})
	}
}

func TestHandlers_BalanceHandler(t *testing.T) {
	players := []entity.Player{
		{ID: "balance_ok", Points: 200},
		{ID: "balance_not_found", Points: 100},
	}
	controller.On("Balance", players[0].ID).Return(players[0], nil)
	controller.On("Balance", players[1].ID).Return(entity.Player{}, errors.Error{Code: errors.NotFoundError})
	client := http.Client{}
	tt := []struct {
		name           string
		playerID       string
		err            error
		expectedError  errors.Error
		expectedPlayer entity.Player
		expectedStatus int
	}{
		{
			name:           "balance: ok",
			playerID:       players[0].ID,
			expectedError:  errors.Error{},
			expectedPlayer: players[0],
			expectedStatus: http.StatusOK,
		},
		{
			name:           "balance: not found",
			playerID:       players[1].ID,
			expectedError:  errors.Error{Code: errors.NotFoundError},
			expectedPlayer: entity.Player{},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%v/balance?playerId=%v", ts.URL, tc.playerID), nil)
			assert.Equal(t, tc.err, err)
			res, err := client.Do(req)
			assert.Equal(t, tc.err, err)
			assert.Equal(t, tc.expectedStatus, res.StatusCode)
			decoder := json.NewDecoder(res.Body)
			var player entity.Player
			err = decoder.Decode(&player)
			assert.Equal(t, tc.err, err)
			assert.Equal(t, tc.expectedPlayer, player)
			var expErr errors.Error
			err = decoder.Decode(&expErr)
			assert.Equal(t, tc.err, err)
			assert.Equal(t, tc.expectedError, expErr)
		})
	}
}

func TestHandlers_AnnounceHandler(t *testing.T) {
	tournaments := []entity.Tournament{
		{ID: "announce_ok_and_duplicated", Deposit: 100},
	}
	controller.On("AnnounceTournament", tournaments[0].ID, tournaments[0].Deposit).Return(nil).Once()
	controller.On("AnnounceTournament", tournaments[0].ID, tournaments[0].Deposit).Return(errors.Error{Code: errors.DuplicatedIDError})
	client := http.Client{}
	tt := []struct {
		name           string
		tournamentID   string
		deposit        interface{}
		err            error
		expectedError  errors.Error
		expectedStatus int
	}{
		{
			name:           "announce: ok",
			tournamentID:   tournaments[0].ID,
			deposit:        tournaments[0].Deposit,
			expectedError:  errors.Error{},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "announce: duplicated id",
			tournamentID:   tournaments[0].ID,
			deposit:        tournaments[0].Deposit,
			expectedError:  errors.Error{Code: errors.DuplicatedIDError},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "announce: incorrect deposit",
			tournamentID:   tournaments[0].ID,
			deposit:        "incorrect_deposit",
			expectedError:  errors.Error{Code: errors.NotNumberError, Message: "cannot create tournament, deposit is not number: incorrect_deposit", Info: "strconv.Atoi: parsing \"incorrect_deposit\": invalid syntax"},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%v/announceTournament?tournamentId=%v&deposit=%v", ts.URL, tc.tournamentID, tc.deposit), nil)
			assert.Equal(t, tc.err, err)
			res, err := client.Do(req)
			assert.Equal(t, tc.err, err)
			assert.Equal(t, tc.expectedStatus, res.StatusCode)
			if tc.expectedStatus != http.StatusOK {
				decoder := json.NewDecoder(res.Body)
				var expErr errors.Error
				err = decoder.Decode(&expErr)
				assert.Equal(t, tc.err, err)
				assert.Equal(t, tc.expectedError, expErr)
			}
		})
	}
}

func TestHandlers_JoinHandler(t *testing.T) {
	tournaments := []entity.Tournament{
		{ID: "join_ok", Deposit: 100},
		{ID: "join_unexpected", Deposit: 100},
	}
	players := []entity.Player{
		{ID: "join_ok", Points: 350},
		{ID: "join_unexpected", Points: 200},
	}
	controller.On("JoinTournament", tournaments[0].ID, players[0].ID).Return(nil)
	controller.On("JoinTournament", tournaments[1].ID, players[1].ID).Return(e.New("unexpected"))
	client := http.Client{}
	tt := []struct {
		name           string
		tourID         string
		playerID       string
		err            error
		expectedError  errors.Error
		expectedStatus int
	}{
		{
			name:           "join: ok",
			tourID:         tournaments[0].ID,
			playerID:       players[0].ID,
			expectedError:  errors.Error{},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "join: unexpected error",
			tourID:         tournaments[1].ID,
			playerID:       players[1].ID,
			expectedError:  errors.Error{Code: "UnknownError", Message: "unexpected"},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%v/joinTournament?tournamentId=%v&playerId=%v", ts.URL, tc.tourID, tc.playerID), nil)
			assert.Equal(t, tc.err, err)
			res, err := client.Do(req)
			assert.Equal(t, tc.err, err)
			assert.Equal(t, tc.expectedStatus, res.StatusCode)
			if tc.expectedStatus != http.StatusOK {
				decoder := json.NewDecoder(res.Body)
				var expErr errors.Error
				err = decoder.Decode(&expErr)
				assert.Equal(t, tc.err, err)
				assert.Equal(t, tc.expectedError, expErr)
			}
		})
	}
}

func TestHandlers_ResultHandler(t *testing.T) {
	tournaments := []entity.Tournament{
		{ID: "result_ok", Deposit: 100},
		{ID: "result_none_participants", Deposit: 100},
	}
	winners := []entity.Winners{
		{Winners: []entity.Winner{entity.Winner{ID: "result_ok"}}},
	}
	controller.On("Results", tournaments[0].ID).Return(winners[0], nil)
	controller.On("Results", tournaments[1].ID).Return(entity.Winners{}, errors.Error{Code: errors.NoneParticipantsError})
	client := http.Client{}
	tt := []struct {
		name            string
		tourID          string
		err             error
		expectedWinners entity.Winners
		expectedError   error
		expectedStatus  int
	}{
		{
			name:            "result: ok",
			tourID:          tournaments[0].ID,
			expectedWinners: winners[0],
			expectedError:   errors.Error{},
			expectedStatus:  http.StatusOK,
		},
		{
			name:            "result: none participants",
			tourID:          tournaments[1].ID,
			expectedWinners: entity.Winners{},
			expectedError:   errors.Error{Code: errors.NoneParticipantsError},
			expectedStatus:  http.StatusOK,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%v/resultTournament?tournamentId=%v", ts.URL, tc.tourID), nil)
			assert.Equal(t, tc.err, err)
			res, err := client.Do(req)
			assert.Equal(t, tc.err, err)
			assert.Equal(t, tc.expectedStatus, res.StatusCode)
			decoder := json.NewDecoder(res.Body)
			var winner entity.Winners
			err = decoder.Decode(&winner)
			assert.Equal(t, tc.err, err)
			assert.Equal(t, tc.expectedWinners, winner)
			var expErr errors.Error
			err = decoder.Decode(&expErr)
			assert.Equal(t, tc.err, err)
			assert.Equal(t, tc.expectedError, expErr)
		})
	}
}

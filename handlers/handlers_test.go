package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Tournament/controller"
	"github.com/Tournament/entity"
	"github.com/Tournament/errors"
	"github.com/Tournament/postgres"
)

var (
	ts *httptest.Server
)

func TestMain(m *testing.M) {
	p, err := postgres.NewDB("user=postgres dbname=postgres password=password sslmode=disable")
	if err != nil {
		log.Fatalf("Cannot open database: %v", err)
	}
	r := NewRouter(Server{Controller: controller.Game{DB: p}})
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
		fund                   int
		expectedNegativeStatus int
		expectedError          error
		expectedPlayer         entity.Player
		expectedStatus         int
	}{
		{
			name:                   "ok_createfund0",
			player:                 players[0],
			fund:                   0,
			expectedError:          nil,
			expectedNegativeStatus: http.StatusOK,
			expectedPlayer:         players[0],
			expectedStatus:         http.StatusCreated,
		},
		{
			name:                   "ok_createfund1",
			player:                 players[1],
			fund:                   0,
			expectedError:          nil,
			expectedNegativeStatus: http.StatusOK,
			expectedPlayer:         players[1],
			expectedStatus:         http.StatusCreated,
		},
		{
			name:                   "ok_createfund2",
			player:                 players[2],
			fund:                   0,
			expectedError:          nil,
			expectedNegativeStatus: http.StatusOK,
			expectedPlayer:         players[2],
			expectedStatus:         http.StatusCreated,
		},
		{
			name:                   "ok_fund0",
			player:                 players[0],
			fund:                   -100,
			expectedError:          nil,
			expectedNegativeStatus: http.StatusNotFound,
			expectedPlayer:         entity.Player{ID: players[0].ID, Points: players[0].Points * 2},
			expectedStatus:         http.StatusOK,
		},
		{
			name:                   "ok_fund1",
			player:                 players[1],
			fund:                   -100,
			expectedError:          nil,
			expectedNegativeStatus: http.StatusNotFound,
			expectedPlayer:         entity.Player{ID: players[1].ID, Points: players[1].Points * 2},
			expectedStatus:         http.StatusOK,
		},
		{
			name:                   "ok_fund2",
			player:                 players[2],
			fund:                   -100,
			expectedError:          nil,
			expectedNegativeStatus: http.StatusNotFound,
			expectedPlayer:         entity.Player{ID: players[2].ID, Points: players[2].Points * 2},
			expectedStatus:         http.StatusOK,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%v/fund?playerId=%v&points=%v", ts.URL, tc.player.ID, tc.player.Points), nil)
			assert.Equal(t, tc.expectedError, err)
			res, err := client.Do(req)
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedStatus, res.StatusCode)

			req, err = http.NewRequest(http.MethodGet, fmt.Sprintf("%v/balance?playerId=%v", ts.URL, tc.player.ID), nil)
			assert.Equal(t, tc.expectedError, err)
			res, err = client.Do(req)
			assert.Equal(t, tc.expectedError, err)
			decoder := json.NewDecoder(res.Body)
			var player entity.Player
			err = decoder.Decode(&player)
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedPlayer, player)

			req, err = http.NewRequest(http.MethodPut, fmt.Sprintf("%v/fund?playerId=%v&points=%v", ts.URL, tc.player.ID, tc.fund), nil)
			assert.Equal(t, tc.expectedError, err)
			res, err = client.Do(req)
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedNegativeStatus, res.StatusCode)
		})
	}

	for i := range players {
		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%v/deletePlayer?playerId=%v", ts.URL, players[i].ID), nil)
		if err != nil {
			t.Fatal(err)
		}
		client.Do(req)
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
	for i := range players {
		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%v/fund?playerId=%v&points=%v", ts.URL, players[i].ID, players[i].Points), nil)
		if err != nil {
			t.Fatal(err)
		}
		client.Do(req)
	}
	tt := []struct {
		name           string
		player         entity.Player
		take           int
		expectedError  error
		expectedPlayer entity.Player
		expectedStatus int
	}{
		{
			name:           "ok_take0",
			player:         players[0],
			take:           take[0],
			expectedError:  nil,
			expectedPlayer: entity.Player{ID: players[0].ID, Points: players[0].Points - take[0]},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "ok_take1",
			player:         players[1],
			take:           take[1],
			expectedError:  nil,
			expectedPlayer: entity.Player{ID: players[1].ID, Points: players[1].Points - take[1]},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "ok_take2",
			player:         players[2],
			take:           take[2],
			expectedError:  nil,
			expectedPlayer: entity.Player{ID: players[2].ID, Points: players[2].Points - take[2]},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "err_take0",
			player:         players[0],
			take:           take[0],
			expectedError:  nil,
			expectedPlayer: entity.Player{ID: players[0].ID, Points: players[0].Points - take[0]},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "err_take1",
			player:         players[1],
			take:           take[1],
			expectedError:  nil,
			expectedPlayer: entity.Player{ID: players[1].ID, Points: players[1].Points - take[1]},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "err_take2",
			player:         players[2],
			take:           take[2],
			expectedError:  nil,
			expectedPlayer: entity.Player{ID: players[2].ID, Points: players[2].Points - take[2]},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%v/take?playerId=%v&points=%v", ts.URL, tc.player.ID, tc.take), nil)
			assert.Equal(t, tc.expectedError, err)
			res, err := client.Do(req)
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedStatus, res.StatusCode)
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

	for i := range players {
		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%v/deletePlayer?playerId=%v", ts.URL, players[i].ID), nil)
		if err != nil {
			t.Fatal(err)
		}
		client.Do(req)
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
	for i := range players {
		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%v/fund?playerId=%v&points=%v", ts.URL, players[i].ID, players[i].Points), nil)
		if err != nil {
			t.Fatal(err)
		}
		client.Do(req)
	}
	tt := []struct {
		name           string
		player         entity.Player
		take           int
		fund           int
		expectedError  error
		expectedPlayer entity.Player
		expectedStatus int
	}{
		{
			name:           "balance_fund0",
			player:         players[0],
			take:           0,
			fund:           fund[0],
			expectedError:  nil,
			expectedPlayer: entity.Player{ID: players[0].ID, Points: players[0].Points + fund[0]},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "balance_fund1",
			player:         players[1],
			take:           0,
			fund:           fund[1],
			expectedError:  nil,
			expectedPlayer: entity.Player{ID: players[1].ID, Points: players[1].Points + fund[1]},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "balance_fund2",
			player:         players[2],
			take:           0,
			fund:           fund[2],
			expectedError:  nil,
			expectedPlayer: entity.Player{ID: players[2].ID, Points: players[2].Points + fund[2]},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "balance_take0",
			player:         players[0],
			take:           take[0],
			fund:           0,
			expectedError:  nil,
			expectedPlayer: entity.Player{ID: players[0].ID, Points: players[0].Points + fund[0] - take[0]},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "balance_take1",
			player:         players[1],
			take:           take[1],
			fund:           fund[1],
			expectedError:  nil,
			expectedPlayer: entity.Player{ID: players[1].ID, Points: players[1].Points + fund[1] + fund[1] - take[1]},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "balance_take2",
			player:         players[2],
			take:           take[2],
			fund:           0,
			expectedError:  nil,
			expectedPlayer: entity.Player{ID: players[2].ID, Points: players[2].Points + fund[2] - take[2]},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%v/take?playerId=%v&points=%v", ts.URL, tc.player.ID, tc.take), nil)
			assert.Equal(t, tc.expectedError, err)
			res, err := client.Do(req)
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedStatus, res.StatusCode)
			req, err = http.NewRequest(http.MethodPut, fmt.Sprintf("%v/fund?playerId=%v&points=%v", ts.URL, tc.player.ID, tc.fund), nil)
			assert.Equal(t, tc.expectedError, err)
			res, err = client.Do(req)
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedStatus, res.StatusCode)
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

	for i := range players {
		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%v/deletePlayer?playerId=%v", ts.URL, players[i].ID), nil)
		if err != nil {
			t.Fatal(err)
		}
		client.Do(req)
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
		tournament     entity.Tournament
		expectedError  error
		expectedStatus int
	}{
		{
			name:           "announce_ok0",
			tournament:     tournaments[0],
			expectedError:  nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "announce_ok1",
			tournament:     tournaments[1],
			expectedError:  nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "announce_ok2",
			tournament:     tournaments[2],
			expectedError:  nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "announce_duplicated0",
			tournament:     tournaments[0],
			expectedError:  nil,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "announce_duplicated1",
			tournament:     tournaments[1],
			expectedError:  nil,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "announce_duplicated2",
			tournament:     tournaments[2],
			expectedError:  nil,
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%v/announceTournament?tournamentId=%v&deposit=%v", ts.URL, tc.tournament.ID, tc.tournament.Deposit), nil)
			assert.Equal(t, tc.expectedError, err)
			res, err := client.Do(req)
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedStatus, res.StatusCode)
		})
	}

	for i := range tournaments {
		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%v/deleteTournament?tournamentId=%v", ts.URL, tournaments[i].ID), nil)
		if err != nil {
			t.Fatal(err)
		}
		client.Do(req)
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
	for i := range tournaments {
		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%v/announceTournament?tournamentId=%v&deposit=%v", ts.URL, tournaments[i].ID, tournaments[i].Deposit), nil)
		if err != nil {
			t.Fatal(err)
		}
		_, err = client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
	}
	for i := range players {
		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%v/fund?playerId=%v&points=%v", ts.URL, players[i].ID, players[i].Points), nil)
		if err != nil {
			t.Fatal(err)
		}
		_, err = client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
	}
	tt := []struct {
		name               string
		tournament         entity.Tournament
		participants       []entity.Player
		expectedPlayers    []entity.Player
		expectedJoinStatus []int
		expectedJoinError  []error
		expectedError      error
	}{
		{
			name:         "announce_ok0",
			tournament:   tournaments[0],
			participants: players,
			expectedPlayers: []entity.Player{
				{ID: players[0].ID, Points: players[0].Points - tournaments[0].Deposit},
				{ID: players[1].ID, Points: players[1].Points - tournaments[0].Deposit},
				{ID: players[2].ID, Points: players[2].Points - tournaments[0].Deposit}},
			expectedJoinStatus: []int{http.StatusOK, http.StatusOK, http.StatusOK},
			expectedJoinError:  []error{nil, nil, nil},
			expectedError:      nil,
		},
		{
			name:         "announce_ok1",
			tournament:   tournaments[1],
			participants: players,
			expectedPlayers: []entity.Player{
				{ID: players[0].ID, Points: players[0].Points - tournaments[1].Deposit - tournaments[0].Deposit},
				{ID: players[1].ID, Points: players[1].Points - tournaments[1].Deposit - tournaments[0].Deposit},
				{ID: players[2].ID, Points: players[2].Points - tournaments[0].Deposit}},
			expectedJoinError:  []error{nil, nil, errors.Error{Code: errors.NegativePointsNumberError, Message: "update player: cannot update points numbers, dif -100"}},
			expectedJoinStatus: []int{http.StatusOK, http.StatusOK, http.StatusNotFound},
			expectedError:      nil,
		},
		{
			name:         "announce_ok2",
			tournament:   tournaments[2],
			participants: players,
			expectedPlayers: []entity.Player{
				{ID: players[0].ID, Points: players[0].Points - tournaments[2].Deposit - tournaments[1].Deposit - tournaments[0].Deposit},
				{ID: players[1].ID, Points: players[1].Points - tournaments[1].Deposit - tournaments[0].Deposit},
				{ID: players[2].ID, Points: players[2].Points - tournaments[0].Deposit}},
			expectedJoinError:  []error{nil, errors.Error{Code: errors.NegativePointsNumberError, Message: "update player: cannot update points numbers, dif -100"}, errors.Error{Code: errors.NegativePointsNumberError, Message: "update player: cannot update points numbers, dif -100"}},
			expectedJoinStatus: []int{http.StatusOK, http.StatusNotFound, http.StatusNotFound},
			expectedError:      nil,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			for i := range tc.participants {
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

	for i := range tournaments {
		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%v/deleteTournament?tournamentId=%v", ts.URL, tournaments[i].ID), nil)
		if err != nil {
			t.Fatal(err)
		}
		client.Do(req)
	}
	for i := range players {
		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%v/deletePlayer?playerId=%v", ts.URL, players[i].ID), nil)
		if err != nil {
			t.Fatal(err)
		}
		client.Do(req)
	}
}

func TestHandlers_ResultHandler(t *testing.T) {
	tournaments := []entity.Tournament{
		{ID: "result_tour1", Deposit: 100},
		{ID: "result_tour2", Deposit: 100},
		{ID: "result_tour3", Deposit: 100},
	}
	players := []entity.Player{
		{ID: "result_player1", Points: 300},
		{ID: "result_player2", Points: 200},
		{ID: "result_player3", Points: 100},
	}
	client := http.Client{}
	for i := range tournaments {
		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%v/announceTournament?tournamentId=%v&deposit=%v", ts.URL, tournaments[i].ID, tournaments[i].Deposit), nil)
		if err != nil {
			t.Fatal(err)
		}
		_, err = client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
	}
	for i := range players {
		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%v/fund?playerId=%v&points=%v", ts.URL, players[i].ID, players[i].Points), nil)
		if err != nil {
			t.Fatal(err)
		}
		_, err = client.Do(req)
		if err != nil {
			t.Fatal(err)
		}

		for j := range tournaments {
			req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%v/joinTournament?tournamentId=%v&playerId=%v", ts.URL, tournaments[j].ID, players[i].ID), nil)
			if err != nil {
				t.Fatal(err)
			}
			_, err = client.Do(req)
			if err != nil {
				t.Fatal(err)
			}
		}
	}
	tt := []struct {
		name               string
		id                 string
		expectedError      error
		expectedStatus     int
		expectedPrize      int
		expectedJoinStatus int
		expectedJoinError  error
	}{
		{
			name:               "announce_ok0",
			id:                 tournaments[0].ID,
			expectedError:      nil,
			expectedStatus:     http.StatusOK,
			expectedPrize:      tournaments[0].Deposit * 3,
			expectedJoinStatus: http.StatusNotFound,
			expectedJoinError:  errors.Error{Code: errors.ClosedTournamentError, Message: "join tournament: cannot join to closed tournament, tourID: " + tournaments[0].ID},
		},
		{
			name:               "announce_ok1",
			id:                 tournaments[1].ID,
			expectedError:      nil,
			expectedStatus:     http.StatusOK,
			expectedPrize:      tournaments[1].Deposit * 2,
			expectedJoinStatus: http.StatusNotFound,
			expectedJoinError:  errors.Error{Code: errors.ClosedTournamentError, Message: "join tournament: cannot join to closed tournament, tourID: " + tournaments[1].ID},
		},
		{
			name:               "announce_ok2",
			id:                 tournaments[2].ID,
			expectedError:      nil,
			expectedStatus:     http.StatusOK,
			expectedPrize:      tournaments[2].Deposit,
			expectedJoinStatus: http.StatusNotFound,
			expectedJoinError:  errors.Error{Code: errors.ClosedTournamentError, Message: "join tournament: cannot join to closed tournament, tourID: " + tournaments[2].ID},
		},
		{
			name:               "announce_fake",
			id:                 "announce_fake",
			expectedError:      nil,
			expectedStatus:     http.StatusNotFound,
			expectedPrize:      0,
			expectedJoinStatus: http.StatusNotFound,
			expectedJoinError:  errors.Error{Code: errors.NotFoundError, Message: "get state: cannot get tournament state from not existing tournament, id: announce_fake"},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
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

	for i := range tournaments {
		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%v/deleteTournament?tournamentId=%v", ts.URL, tournaments[i].ID), nil)
		if err != nil {
			t.Fatal(err)
		}
		client.Do(req)
	}
	for i := range players {
		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%v/deletePlayer?playerId=%v", ts.URL, players[i].ID), nil)
		if err != nil {
			t.Fatal(err)
		}
		client.Do(req)
	}
}

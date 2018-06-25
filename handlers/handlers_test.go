package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Tournament/controller"
	"github.com/Tournament/entity"
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
	r := NewRouter(Server{Controller: controller.Game{TDB: p, PDB: p}})
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
		name           string
		player         entity.Player
		expectedError  error
		expectedPlayer entity.Player
		expectedStatus int
	}{
		{
			name:           "ok_createfund0",
			player:         players[0],
			expectedError:  nil,
			expectedPlayer: players[0],
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "ok_createfund1",
			player:         players[1],
			expectedError:  nil,
			expectedPlayer: players[1],
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "ok_createfund2",
			player:         players[2],
			expectedError:  nil,
			expectedPlayer: players[2],
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "ok_fund0",
			player:         players[0],
			expectedError:  nil,
			expectedPlayer: entity.Player{ID: players[0].ID, Points: players[0].Points * 2},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "ok_fund1",
			player:         players[1],
			expectedError:  nil,
			expectedPlayer: entity.Player{ID: players[1].ID, Points: players[1].Points * 2},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "ok_fund2",
			player:         players[2],
			expectedError:  nil,
			expectedPlayer: entity.Player{ID: players[2].ID, Points: players[2].Points * 2},
			expectedStatus: http.StatusOK,
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
			rawPlayer, err := ioutil.ReadAll(res.Body)
			assert.Equal(t, tc.expectedError, err)
			var player entity.Player
			err = json.Unmarshal(rawPlayer, &player)
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
			rawPlayer, err := ioutil.ReadAll(res.Body)
			assert.Equal(t, tc.expectedError, err)
			var player entity.Player
			err = json.Unmarshal(rawPlayer, &player)
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
			rawPlayer, err := ioutil.ReadAll(res.Body)
			assert.Equal(t, tc.expectedError, err)
			var player entity.Player
			err = json.Unmarshal(rawPlayer, &player)
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

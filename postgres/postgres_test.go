package postgres

import (
	"log"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/dmitriyomelyusik/Tournament/entity"
	"github.com/dmitriyomelyusik/Tournament/errors"
	"github.com/stretchr/testify/assert"
)

var (
	p *Postgres
)

func TestMain(m *testing.M) {
	var err error
	p, err = NewDB("user=postgres dbname=postgres password=password sslmode=disable")
	if err != nil {
		log.Fatalf("Cannot open database: %v", err)
	}
	code := m.Run()
	os.Exit(code)
}

func TestPlayer_CreatePlayer(t *testing.T) {
	players := []entity.Player{
		{ID: "createplayer_test1", Points: 200},
		{ID: "createplayer_test2", Points: 200},
		{ID: "createplayer_test3", Points: 200},
	}
	tt := []struct {
		name          string
		player        entity.Player
		expectedError error
	}{
		{
			name:          "ok0",
			player:        players[0],
			expectedError: nil,
		},
		{
			name:          "ok1",
			player:        players[1],
			expectedError: nil,
		},
		{
			name:          "duplicate0",
			player:        players[0],
			expectedError: errors.Error{Code: errors.DuplicatedIDError, Message: "create player: using duplicated id to create player, id " + players[0].ID},
		},
		{
			name:          "ok2",
			player:        players[2],
			expectedError: nil,
		},
		{
			name:          "duplicate1",
			player:        players[1],
			expectedError: errors.Error{Code: errors.DuplicatedIDError, Message: "create player: using duplicated id to create player, id " + players[1].ID},
		},
		{
			name:          "duplicate2",
			player:        players[2],
			expectedError: errors.Error{Code: errors.DuplicatedIDError, Message: "create player: using duplicated id to create player, id " + players[2].ID},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			player, err := p.CreatePlayer(tc.player.ID, tc.player.Points)
			assert.Equal(t, tc.expectedError, err)
			if err == nil {
				assert.Equal(t, tc.player, player)
			}
		})
	}

	for i := range players {
		p.DeletePlayer(players[i].ID)
	}
}

func TestPlayer_GetPlayer(t *testing.T) {
	players := []entity.Player{
		{ID: "createplayer_test1", Points: 200},
		{ID: "createplayer_test2", Points: 200},
		{ID: "createplayer_test3", Points: 200},
	}
	for i := range players {
		_, err := p.CreatePlayer(players[i].ID, players[i].Points)
		require.NoError(t, err)
		defer func(i int) {
			err = p.DeletePlayer(players[i].ID)
			require.NoError(t, err)
		}(i)
	}
	tt := []struct {
		name           string
		id             string
		expectedPlayer entity.Player
		expectedError  error
	}{
		{
			name:           "ok0",
			id:             players[0].ID,
			expectedPlayer: players[0],
			expectedError:  nil,
		},
		{
			name:           "ok1",
			id:             players[1].ID,
			expectedPlayer: players[1],
			expectedError:  nil,
		},
		{
			name:           "ok2",
			id:             players[2].ID,
			expectedPlayer: players[2],
			expectedError:  nil,
		},
		{
			name:           "not existing player",
			id:             "getplayer_fake",
			expectedPlayer: entity.Player{},
			expectedError:  errors.Error{Code: errors.NotFoundError, Message: "get player: cannot find player, id getplayer_fake"},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			p, err := p.GetPlayer(tc.id)
			assert.Equal(t, tc.expectedPlayer, p)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestPlayer_UpdatePlayer(t *testing.T) {
	players := []entity.Player{
		{ID: "updateplayer_test1", Points: 200},
		{ID: "updateplayer_test2", Points: 200},
		{ID: "updateplayer_test3", Points: 200},
	}
	for i := range players {
		_, err := p.CreatePlayer(players[i].ID, players[i].Points)
		require.NoError(t, err)
		defer func(i int) {
			err = p.DeletePlayer(players[i].ID)
			require.NoError(t, err)
		}(i)
	}
	tt := []struct {
		name          string
		id            string
		dif           int
		expectedError error
	}{
		{
			name:          "ok_fund0",
			id:            players[0].ID,
			dif:           200,
			expectedError: nil,
		},
		{
			name:          "ok_fund1",
			id:            players[1].ID,
			dif:           200,
			expectedError: nil,
		},
		{
			name:          "ok_fund2",
			id:            players[2].ID,
			dif:           200,
			expectedError: nil,
		},
		{
			name:          "not existing player",
			id:            "getplayer_fake",
			dif:           200,
			expectedError: errors.Error{Code: errors.NotFoundError, Message: "update player: cannot find player, id getplayer_fake"},
		},
		{
			name:          "ok_take0",
			id:            players[0].ID,
			dif:           -100,
			expectedError: nil,
		},
		{
			name:          "ok_take1",
			id:            players[1].ID,
			dif:           -100,
			expectedError: nil,
		},
		{
			name:          "ok_take2",
			id:            players[2].ID,
			dif:           -100,
			expectedError: nil,
		},
		{
			name:          "err_take0",
			id:            players[0].ID,
			dif:           -500,
			expectedError: errors.Error{Code: errors.NegativePointsNumberError, Message: "update player: cannot update points numbers, dif -500"},
		},
		{
			name:          "err_take1",
			id:            players[1].ID,
			dif:           -600,
			expectedError: errors.Error{Code: errors.NegativePointsNumberError, Message: "update player: cannot update points numbers, dif -600"},
		},
		{
			name:          "err_take2",
			id:            players[2].ID,
			dif:           -700,
			expectedError: errors.Error{Code: errors.NegativePointsNumberError, Message: "update player: cannot update points numbers, dif -700"},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := p.UpdatePlayer(tc.id, tc.dif)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestPlayer_DeletePlayer(t *testing.T) {
	players := []entity.Player{
		{ID: "deleteplayer_test1", Points: 200},
		{ID: "deleteplayer_test2", Points: 200},
		{ID: "deleteplayer_test3", Points: 200},
	}
	for i := range players {
		_, err := p.CreatePlayer(players[i].ID, players[i].Points)
		require.NoError(t, err)
	}
	tt := []struct {
		name          string
		id            string
		expectedError error
	}{
		{
			name:          "ok0",
			id:            players[0].ID,
			expectedError: nil,
		},
		{
			name:          "already deleted",
			id:            players[0].ID,
			expectedError: errors.Error{Code: errors.NotFoundError, Message: "delete player: player does not exist, id " + players[0].ID},
		},
		{
			name:          "ok1",
			id:            players[1].ID,
			expectedError: nil,
		},
		{
			name:          "ok2",
			id:            players[2].ID,
			expectedError: nil,
		},
		{
			name:          "already deleted",
			id:            players[1].ID,
			expectedError: errors.Error{Code: errors.NotFoundError, Message: "delete player: player does not exist, id " + players[1].ID},
		},
		{
			name:          "wasn't exist",
			id:            "deleteplayer_fake",
			expectedError: errors.Error{Code: errors.NotFoundError, Message: "delete player: player does not exist, id deleteplayer_fake"},
		},
		{
			name:          "already deleted",
			id:            players[2].ID,
			expectedError: errors.Error{Code: errors.NotFoundError, Message: "delete player: player does not exist, id " + players[2].ID},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := p.DeletePlayer(tc.id)
			assert.Equal(t, tc.expectedError, err)
		})
	}

	for i := range players {
		p.DeletePlayer(players[i].ID)
	}
}

func TestTournament_CreateTournament(t *testing.T) {
	tournaments := []entity.Tournament{
		{ID: "createtournament_test1", Deposit: 100},
		{ID: "createtournament_test2", Deposit: 100},
		{ID: "createtournament_test3", Deposit: 100},
	}
	tt := []struct {
		name          string
		tournament    entity.Tournament
		expectedError error
	}{
		{
			name:          "ok0",
			tournament:    tournaments[0],
			expectedError: nil,
		},
		{
			name:          "ok1",
			tournament:    tournaments[1],
			expectedError: nil,
		},
		{
			name:          "duplicated id 0",
			tournament:    tournaments[0],
			expectedError: errors.Error{Code: errors.DuplicatedIDError, Message: "create tournament: using duplicated id to create tournament, id: " + tournaments[0].ID},
		},
		{
			name:          "ok2",
			tournament:    tournaments[2],
			expectedError: nil,
		},
		{
			name:          "duplicated id 1",
			tournament:    tournaments[1],
			expectedError: errors.Error{Code: errors.DuplicatedIDError, Message: "create tournament: using duplicated id to create tournament, id: " + tournaments[1].ID},
		},
		{
			name:          "duplicated id 2",
			tournament:    tournaments[2],
			expectedError: errors.Error{Code: errors.DuplicatedIDError, Message: "create tournament: using duplicated id to create tournament, id: " + tournaments[2].ID},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := p.CreateTournament(tc.tournament.ID, tc.tournament.Deposit)
			assert.Equal(t, tc.expectedError, err)
		})
	}

	for i := range tournaments {
		p.DeleteTournament(tournaments[i].ID)
	}
}

func TestTournament_DeleteTournament(t *testing.T) {
	tournaments := []entity.Tournament{
		{ID: "deletetournament_test1", Deposit: 100},
		{ID: "deletetournament_test2", Deposit: 100},
		{ID: "deletetournament_test3", Deposit: 100},
	}
	for i := range tournaments {
		err := p.CreateTournament(tournaments[i].ID, tournaments[i].Deposit)
		require.NoError(t, err)
	}
	tt := []struct {
		name          string
		id            string
		expectedError error
	}{
		{
			name:          "ok0",
			id:            tournaments[0].ID,
			expectedError: nil,
		},
		{
			name:          "already deleted",
			id:            tournaments[0].ID,
			expectedError: errors.Error{Code: errors.NotFoundError, Message: "delete tournament: tournament does not exist, id " + tournaments[0].ID},
		},
		{
			name:          "ok1",
			id:            tournaments[1].ID,
			expectedError: nil,
		},
		{
			name:          "fake id",
			id:            "deletetournaments_fake",
			expectedError: errors.Error{Code: errors.NotFoundError, Message: "delete tournament: tournament does not exist, id deletetournaments_fake"},
		},
		{
			name:          "ok2",
			id:            tournaments[2].ID,
			expectedError: nil,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := p.DeleteTournament(tc.id)
			assert.Equal(t, tc.expectedError, err)
		})
	}

	for i := range tournaments {
		p.DeleteTournament(tournaments[i].ID)
	}
}

func TestTournament_CloseTournament(t *testing.T) {
	tournaments := []entity.Tournament{
		{ID: "closetournament_test1", Deposit: 100},
		{ID: "closetournament_test2", Deposit: 100},
		{ID: "closetournament_test3", Deposit: 100},
	}
	for i := range tournaments {
		err := p.CreateTournament(tournaments[i].ID, tournaments[i].Deposit)
		require.NoError(t, err)
		defer func(i int) {
			err = p.DeleteTournament(tournaments[i].ID)
		}(i)
	}
	tt := []struct {
		name               string
		id                 string
		expectedState      bool
		expectedStateError error
		expectedError      error
	}{
		{
			name:               "ok0",
			id:                 tournaments[0].ID,
			expectedState:      false,
			expectedStateError: nil,
			expectedError:      nil,
		},
		{
			name:               "ok1",
			id:                 tournaments[1].ID,
			expectedState:      false,
			expectedStateError: nil,
			expectedError:      nil,
		},
		{
			name:               "ok2",
			id:                 tournaments[2].ID,
			expectedState:      false,
			expectedStateError: nil,
			expectedError:      nil,
		},
		{
			name:               "fake id",
			id:                 "closetournaments_fake",
			expectedState:      false,
			expectedStateError: errors.Error{Code: errors.NotFoundError, Message: "get state: cannot get tournament state from not existing tournament, id: closetournaments_fake"},
			expectedError:      errors.Error{Code: errors.NotFoundError, Message: "close tournament: cannot close not existing tournament, id: closetournaments_fake"},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := p.CloseTournament(tc.id)
			assert.Equal(t, tc.expectedError, err)
			state, err := p.GetTournamentState(tc.id)
			assert.Equal(t, tc.expectedState, state)
			assert.Equal(t, tc.expectedStateError, err)
		})
	}
}

func TestTournament_GetParticipants(t *testing.T) {
	tournaments := []entity.Tournament{
		{ID: "getparticipants_test1", Deposit: 100},
		{ID: "getparticipants_test2", Deposit: 100},
		{ID: "getparticipants_test3", Deposit: 100},
	}
	players := []entity.Player{
		{ID: "getpart_test1", Points: 200},
		{ID: "getpart_test2", Points: 200},
		{ID: "getpart_test3", Points: 200},
		{ID: "getpart_test4", Points: 50},
		{ID: "getpart_test5", Points: 50},
		{ID: "getpart_test6", Points: 50},
	}
	expParticipants := [][]string{nil, nil, nil}
	for i := range tournaments {
		err := p.CreateTournament(tournaments[i].ID, tournaments[i].Deposit)
		require.NoError(t, err)
		defer func(i int) {
			err = p.DeleteTournament(tournaments[i].ID)
		}(i)
	}
	for i := range players {
		_, err := p.CreatePlayer(players[i].ID, players[i].Points)
		require.NoError(t, err)
		defer func(i int) {
			err = p.DeletePlayer(players[i].ID)
		}(i)
	}
	rand.Seed(time.Now().UnixNano())
	for i := range players {
		j := rand.Intn(len(tournaments))
		err := p.UpdateTourAndPlayer(tournaments[j].ID, players[i].ID)
		if err != nil {
			continue
		}
		expParticipants[j] = append(expParticipants[j], players[i].ID)
	}
	tt := []struct {
		name                 string
		id                   string
		expectedParticipants []string
		expectedError        error
	}{
		{
			name:                 "ok0",
			id:                   tournaments[0].ID,
			expectedParticipants: expParticipants[0],
			expectedError:        nil,
		},
		{
			name:                 "ok1",
			id:                   tournaments[1].ID,
			expectedParticipants: expParticipants[1],
			expectedError:        nil,
		},
		{
			name:                 "ok2",
			id:                   tournaments[2].ID,
			expectedParticipants: expParticipants[2],
			expectedError:        nil,
		},
		{
			name:                 "fail",
			id:                   "getparts_fake",
			expectedParticipants: nil,
			expectedError:        errors.Error(errors.Error{Code: errors.NotFoundError, Message: "get participants: cannot get participants from not existing tournament, id: getparts_fake"}),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			part, err := p.GetParticipants(tc.id)
			assert.Equal(t, tc.expectedParticipants, part)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestTournament_GetDeposit(t *testing.T) {
	tournaments := []entity.Tournament{
		{ID: "getdeposit_test1", Deposit: 100},
		{ID: "getdeposit_test2", Deposit: 200},
		{ID: "getdeposit_test3", Deposit: 300},
	}
	for i := range tournaments {
		err := p.CreateTournament(tournaments[i].ID, tournaments[i].Deposit)
		require.NoError(t, err)
		defer func(i int) {
			err = p.DeleteTournament(tournaments[i].ID)
		}(i)
	}
	tt := []struct {
		name            string
		id              string
		expectedDeposit int
		expectedError   error
	}{
		{
			name:            "ok0",
			id:              tournaments[0].ID,
			expectedDeposit: tournaments[0].Deposit,
			expectedError:   nil,
		},
		{
			name:            "ok1",
			id:              tournaments[1].ID,
			expectedDeposit: tournaments[1].Deposit,
			expectedError:   nil,
		},
		{
			name:            "ok2",
			id:              tournaments[2].ID,
			expectedDeposit: tournaments[2].Deposit,
			expectedError:   nil,
		},
		{
			name:            "fake id",
			id:              "getdeposit_fake",
			expectedDeposit: 0,
			expectedError:   errors.Error{Code: errors.NotFoundError, Message: "get deposit: cannot get deposit from not existing tournament, id: getdeposit_fake"},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			dep, err := p.GetDeposit(tc.id)
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedDeposit, dep)
		})
	}
}

func TestTournament_GetState(t *testing.T) {
	tournaments := []entity.Tournament{
		{ID: "getstate_test1", Deposit: 100},
		{ID: "getstate_test2", Deposit: 200},
		{ID: "getstate_test3", Deposit: 300},
	}
	states := []bool{}
	rand.Seed(time.Now().UnixNano())
	for i := range tournaments {
		err := p.CreateTournament(tournaments[i].ID, tournaments[i].Deposit)
		require.NoError(t, err)
		defer func(i int) {
			err = p.DeleteTournament(tournaments[i].ID)
		}(i)
		if rand.Int()%2 == 0 {
			err = p.CloseTournament(tournaments[i].ID)
			require.NoError(t, err)
			states = append(states, false)
			continue
		}
		states = append(states, true)
	}
	tt := []struct {
		name          string
		id            string
		expectedState bool
		expectedError error
	}{
		{
			name:          "ok0",
			id:            tournaments[0].ID,
			expectedState: states[0],
			expectedError: nil,
		},
		{
			name:          "ok1",
			id:            tournaments[1].ID,
			expectedState: states[1],
			expectedError: nil,
		},
		{
			name:          "ok2",
			id:            tournaments[2].ID,
			expectedState: states[2],
			expectedError: nil,
		},
		{
			name:          "fake id",
			id:            "getstate_fake",
			expectedState: false,
			expectedError: errors.Error{Code: errors.NotFoundError, Message: "get state: cannot get tournament state from not existing tournament, id: getstate_fake"},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			state, err := p.GetTournamentState(tc.id)
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedState, state)
		})
	}
}

func TestTournament_GetWinner(t *testing.T) {
	tournaments := []entity.Tournament{
		{ID: "getwinner_test1", Deposit: 100},
		{ID: "getwinner_test2", Deposit: 100},
		{ID: "getwinner_test3", Deposit: 100},
	}
	winners := []entity.Winner{
		{ID: "getwinner_test1", Points: 50},
		{ID: "getwinner_test2", Points: 200},
		{ID: "getwinner_test3", Points: 200},
		{ID: "getwinner_test4", Points: 50},
		{ID: "getwinner_test5", Points: 200},
		{ID: "getwinner_test6", Points: 50},
	}
	var expWinner []entity.Winner
	for i := range tournaments {
		err := p.CreateTournament(tournaments[i].ID, tournaments[i].Deposit)
		require.NoError(t, err)
		defer func(i int) {
			err = p.DeleteTournament(tournaments[i].ID)
		}(i)
	}
	for i := range winners {
		_, err := p.CreatePlayer(winners[i].ID, winners[i].Prize)
		require.NoError(t, err)
		defer func(i int) {
			err = p.DeletePlayer(winners[i].ID)
		}(i)
	}
	rand.Seed(time.Now().UnixNano())
	for i := range tournaments {
		j := rand.Intn(len(winners))
		err := p.SetTournamentWinner(tournaments[i].ID, winners[j])
		require.NoError(t, err)
		expWinner = append(expWinner, winners[j])
	}
	tt := []struct {
		name           string
		id             string
		expectedWinner entity.Winners
		expectedError  error
	}{
		{
			name:           "ok0",
			id:             tournaments[0].ID,
			expectedWinner: entity.Winners{Winners: []entity.Winner{expWinner[0]}},
			expectedError:  nil,
		},
		{
			name:           "ok1",
			id:             tournaments[1].ID,
			expectedWinner: entity.Winners{Winners: []entity.Winner{expWinner[1]}},
			expectedError:  nil,
		},
		{
			name:           "ok2",
			id:             tournaments[2].ID,
			expectedWinner: entity.Winners{Winners: []entity.Winner{expWinner[2]}},
			expectedError:  nil,
		},
		{
			name:           "fail",
			id:             "getwinner_fake",
			expectedWinner: entity.Winners{},
			expectedError:  errors.Error(errors.Error{Code: errors.NotFoundError, Message: "get winner: cannot get winner from not existing tournament, id: getwinner_fake"}),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			winner, err := p.GetWinner(tc.id)
			assert.Equal(t, tc.expectedWinner, winner)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestTournament_SetWinner(t *testing.T) {
	tournaments := []entity.Tournament{
		{ID: "setwinner_test1", Deposit: 100},
		{ID: "setwinner_test2", Deposit: 100},
		{ID: "setwinner_test3", Deposit: 100},
	}
	winners := []entity.Winner{
		{ID: "setwinner_test1", Points: 50},
		{ID: "setwinner_test2", Points: 200},
		{ID: "setwinner_test3", Points: 200},
	}
	for i := range tournaments {
		err := p.CreateTournament(tournaments[i].ID, tournaments[i].Deposit)
		require.NoError(t, err)
		defer func(i int) {
			err = p.DeleteTournament(tournaments[i].ID)
		}(i)
	}
	for i := range winners {
		_, err := p.CreatePlayer(winners[i].ID, winners[i].Prize)
		require.NoError(t, err)
		defer func(i int) {
			err = p.DeletePlayer(winners[i].ID)
		}(i)
	}
	tt := []struct {
		name           string
		id             string
		winner         entity.Winner
		expectedWinner entity.Winners
		expectedError  error
	}{
		{
			name:           "ok0",
			id:             tournaments[0].ID,
			winner:         winners[0],
			expectedWinner: entity.Winners{Winners: []entity.Winner{winners[0]}},
			expectedError:  nil,
		},
		{
			name:           "ok1",
			id:             tournaments[1].ID,
			winner:         winners[1],
			expectedWinner: entity.Winners{Winners: []entity.Winner{winners[1]}},
			expectedError:  nil,
		},
		{
			name:           "ok2",
			id:             tournaments[2].ID,
			winner:         winners[2],
			expectedWinner: entity.Winners{Winners: []entity.Winner{winners[2]}},
			expectedError:  nil,
		},
		{
			name:           "fail",
			id:             "setwinner_fake",
			winner:         entity.Winner{},
			expectedWinner: entity.Winners{},
			expectedError:  errors.Error(errors.Error{Code: errors.NotFoundError, Message: "set winner: tournament not exist, id: setwinner_fake\nsql: no rows in result set"}),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := p.SetTournamentWinner(tc.id, tc.winner)
			assert.Equal(t, tc.expectedError, err)
			if err == nil {
				winner, err := p.GetWinner(tc.id)
				assert.Equal(t, tc.expectedError, err)
				assert.Equal(t, tc.expectedWinner, winner)
			}
		})
	}
}

func TestGama_UpdateTourAndPlayer(t *testing.T) {
	tournaments := []entity.Tournament{
		{ID: "updategame_test1", Deposit: 50},
		{ID: "updategame_test2", Deposit: 50},
		{ID: "updategame_test3", Deposit: 50},
	}
	players := []entity.Player{
		{ID: "updategame_test1", Points: 50},
		{ID: "updategame_test2", Points: 100},
		{ID: "updategame_test3", Points: 150},
	}
	for i := range tournaments {
		err := p.CreateTournament(tournaments[i].ID, tournaments[i].Deposit)
		require.NoError(t, err)
		defer func(i int) {
			err = p.DeleteTournament(tournaments[i].ID)
		}(i)
	}
	for i := range players {
		_, err := p.CreatePlayer(players[i].ID, players[i].Points)
		require.NoError(t, err)
		defer func(i int) {
			err = p.DeletePlayer(players[i].ID)
		}(i)
	}
	tt := []struct {
		name                 string
		id                   string
		participants         []entity.Player
		expectedPoints       []int
		expectedParticipants []string
		expectedUpdateErrors []error
		expectedGetPartError error
		expectedGetPlayError []error
	}{
		{
			name:                 "tour 1",
			id:                   tournaments[0].ID,
			participants:         players,
			expectedPoints:       []int{0, 50, 100},
			expectedParticipants: []string{players[0].ID, players[1].ID, players[2].ID},
			expectedUpdateErrors: []error{nil, nil, nil},
			expectedGetPartError: nil,
			expectedGetPlayError: []error{nil, nil, nil},
		},
		{
			name:                 "tour 2",
			id:                   tournaments[1].ID,
			participants:         players,
			expectedPoints:       []int{0, 0, 50},
			expectedParticipants: []string{players[1].ID, players[2].ID},
			expectedUpdateErrors: []error{errors.Error{Code: errors.NegativePointsNumberError, Message: "update player: cannot update points numbers, dif -50"}, nil, nil},
			expectedGetPartError: nil,
			expectedGetPlayError: []error{nil, nil, nil},
		},
		{
			name:                 "tour 3",
			id:                   tournaments[2].ID,
			participants:         players,
			expectedPoints:       []int{0, 0, 0},
			expectedParticipants: []string{players[2].ID},
			expectedUpdateErrors: []error{errors.Error{Code: errors.NegativePointsNumberError, Message: "update player: cannot update points numbers, dif -50"}, errors.Error{Code: errors.NegativePointsNumberError, Message: "update player: cannot update points numbers, dif -50"}, nil},
			expectedGetPartError: nil,
			expectedGetPlayError: []error{nil, nil, nil},
		},
		{
			name:                 "fake tour",
			id:                   "updatingfaketour",
			participants:         players,
			expectedPoints:       []int{0, 0, 0},
			expectedParticipants: nil,
			expectedUpdateErrors: []error{errors.Error{Code: errors.NotFoundError, Message: "update participiants: cannot update participants in not existing tournament, id: updatingfaketour"}, errors.Error{Code: errors.NotFoundError, Message: "update participiants: cannot update participants in not existing tournament, id: updatingfaketour"}, errors.Error{Code: errors.NotFoundError, Message: "update participiants: cannot update participants in not existing tournament, id: updatingfaketour"}},
			expectedGetPartError: errors.Error{Code: errors.NotFoundError, Message: "get participants: cannot get participants from not existing tournament, id: updatingfaketour"},
			expectedGetPlayError: []error{nil, nil, nil},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			for i, v := range tc.participants {
				err := p.UpdateTourAndPlayer(tc.id, v.ID)
				assert.Equal(t, tc.expectedUpdateErrors[i], err)
			}
			part, err := p.GetParticipants(tc.id)
			assert.Equal(t, tc.expectedGetPartError, err)
			assert.Equal(t, tc.expectedParticipants, part)
			for i, v := range tc.participants {
				player, err := p.GetPlayer(v.ID)
				assert.Equal(t, tc.expectedGetPlayError[i], err)
				assert.Equal(t, tc.expectedPoints[i], player.Points)
			}
		})
	}
}

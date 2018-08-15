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
		{ID: "createplayer_1", Points: 200},
		{ID: "createplayer_2", Points: 200},
		{ID: "createplayer_3", Points: 200},
	}
	tt := []struct {
		name          string
		player        entity.Player
		expectedError error
	}{
		{
			name:          "create player: ok0",
			player:        players[0],
			expectedError: nil,
		},
		{
			name:          "create player: ok1",
			player:        players[1],
			expectedError: nil,
		},
		{
			name:          "create player: duplicate0",
			player:        players[0],
			expectedError: errors.Error{Code: errors.DuplicatedIDError, Message: "create player: using duplicated id to create player, id " + players[0].ID},
		},
		{
			name:          "create player: ok2",
			player:        players[2],
			expectedError: nil,
		},
		{
			name:          "create player: duplicate1",
			player:        players[1],
			expectedError: errors.Error{Code: errors.DuplicatedIDError, Message: "create player: using duplicated id to create player, id " + players[1].ID},
		},
		{
			name:          "create player: duplicate2",
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
		{ID: "getplayer_1", Points: 200},
		{ID: "getplayer_2", Points: 200},
		{ID: "getplayer_3", Points: 200},
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
			name:           "get player: ok0",
			id:             players[0].ID,
			expectedPlayer: players[0],
			expectedError:  nil,
		},
		{
			name:           "get player: ok1",
			id:             players[1].ID,
			expectedPlayer: players[1],
			expectedError:  nil,
		},
		{
			name:           "get player: ok2",
			id:             players[2].ID,
			expectedPlayer: players[2],
			expectedError:  nil,
		},
		{
			name:           "get player: not existing player",
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
		{ID: "updateplayer_1", Points: 200},
		{ID: "updateplayer_2", Points: 200},
		{ID: "updateplayer_3", Points: 200},
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
			name:          "update player: ok fund0",
			id:            players[0].ID,
			dif:           200,
			expectedError: nil,
		},
		{
			name:          "update player: okfund1",
			id:            players[1].ID,
			dif:           200,
			expectedError: nil,
		},
		{
			name:          "update player: ok fund2",
			id:            players[2].ID,
			dif:           200,
			expectedError: nil,
		},
		{
			name:          "update player: not existing player",
			id:            "getplayer_fake",
			dif:           200,
			expectedError: errors.Error{Code: errors.NotFoundError, Message: "update player: cannot find player, id getplayer_fake"},
		},
		{
			name:          "update player: ok take0",
			id:            players[0].ID,
			dif:           -100,
			expectedError: nil,
		},
		{
			name:          "update player: ok take1",
			id:            players[1].ID,
			dif:           -100,
			expectedError: nil,
		},
		{
			name:          "update player: ok take2",
			id:            players[2].ID,
			dif:           -100,
			expectedError: nil,
		},
		{
			name:          "update player: err take0",
			id:            players[0].ID,
			dif:           -500,
			expectedError: errors.Error{Code: errors.NegativePointsNumberError, Message: "update player: cannot update points numbers, dif -500"},
		},
		{
			name:          "update player: err take1",
			id:            players[1].ID,
			dif:           -600,
			expectedError: errors.Error{Code: errors.NegativePointsNumberError, Message: "update player: cannot update points numbers, dif -600"},
		},
		{
			name:          "update player: err take2",
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
		{ID: "deleteplayer_1", Points: 200},
		{ID: "deleteplayer_2", Points: 200},
		{ID: "deleteplayer_3", Points: 200},
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
			name:          "delete player: ok0",
			id:            players[0].ID,
			expectedError: nil,
		},
		{
			name:          "delete player: already deleted",
			id:            players[0].ID,
			expectedError: errors.Error{Code: errors.NotFoundError, Message: "delete player: player does not exist, id " + players[0].ID},
		},
		{
			name:          "delete player: ok1",
			id:            players[1].ID,
			expectedError: nil,
		},
		{
			name:          "delete player: ok2",
			id:            players[2].ID,
			expectedError: nil,
		},
		{
			name:          "delete player: already deleted",
			id:            players[1].ID,
			expectedError: errors.Error{Code: errors.NotFoundError, Message: "delete player: player does not exist, id " + players[1].ID},
		},
		{
			name:          "delete player: wasn't exist",
			id:            "deleteplayer_fake",
			expectedError: errors.Error{Code: errors.NotFoundError, Message: "delete player: player does not exist, id deleteplayer_fake"},
		},
		{
			name:          "delete player: already deleted",
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
		{ID: "createtournament_1", Deposit: 100},
		{ID: "createtournament_2", Deposit: 100},
		{ID: "createtournament_3", Deposit: 100},
	}
	tt := []struct {
		name          string
		tournament    entity.Tournament
		expectedError error
	}{
		{
			name:          "create tournament: ok0",
			tournament:    tournaments[0],
			expectedError: nil,
		},
		{
			name:          "create tournament: ok1",
			tournament:    tournaments[1],
			expectedError: nil,
		},
		{
			name:          "create tournament: duplicated id 0",
			tournament:    tournaments[0],
			expectedError: errors.Error{Code: errors.DuplicatedIDError, Message: "create tournament: using duplicated id to create tournament, id: " + tournaments[0].ID},
		},
		{
			name:          "create tournament: ok2",
			tournament:    tournaments[2],
			expectedError: nil,
		},
		{
			name:          "create tournament: duplicated id 1",
			tournament:    tournaments[1],
			expectedError: errors.Error{Code: errors.DuplicatedIDError, Message: "create tournament: using duplicated id to create tournament, id: " + tournaments[1].ID},
		},
		{
			name:          "create tournament: duplicated id 2",
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
		{ID: "deletetournament_1", Deposit: 100},
		{ID: "deletetournament_2", Deposit: 100},
		{ID: "deletetournament_3", Deposit: 100},
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
			name:          "delete tournament: ok0",
			id:            tournaments[0].ID,
			expectedError: nil,
		},
		{
			name:          "delete tournament: already deleted",
			id:            tournaments[0].ID,
			expectedError: errors.Error{Code: errors.NotFoundError, Message: "delete tournament: tournament does not exist, id " + tournaments[0].ID},
		},
		{
			name:          "delete tournament: ok1",
			id:            tournaments[1].ID,
			expectedError: nil,
		},
		{
			name:          "delete tournament: fake id",
			id:            "deletetournaments_fake",
			expectedError: errors.Error{Code: errors.NotFoundError, Message: "delete tournament: tournament does not exist, id deletetournaments_fake"},
		},
		{
			name:          "delete tournament: ok2",
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
		{ID: "closetournament_1", Deposit: 100},
		{ID: "closetournament_2", Deposit: 100},
		{ID: "closetournament_3", Deposit: 100},
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
			name:               "close tournament: ok0",
			id:                 tournaments[0].ID,
			expectedState:      false,
			expectedStateError: nil,
			expectedError:      nil,
		},
		{
			name:               "close tournament: ok1",
			id:                 tournaments[1].ID,
			expectedState:      false,
			expectedStateError: nil,
			expectedError:      nil,
		},
		{
			name:               "close tournament: ok2",
			id:                 tournaments[2].ID,
			expectedState:      false,
			expectedStateError: nil,
			expectedError:      nil,
		},
		{
			name:               "close tournament: fake id",
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
		{ID: "getparticipants_1", Deposit: 100},
		{ID: "getparticipants_2", Deposit: 100},
		{ID: "getparticipants_3", Deposit: 100},
	}
	players := []entity.Player{
		{ID: "getpart_1", Points: 200},
		{ID: "getpart_2", Points: 200},
		{ID: "getpart_3", Points: 200},
		{ID: "getpart_4", Points: 50},
		{ID: "getpart_5", Points: 50},
		{ID: "getpart_6", Points: 50},
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
			name:                 "get participants: ok0",
			id:                   tournaments[0].ID,
			expectedParticipants: expParticipants[0],
			expectedError:        nil,
		},
		{
			name:                 "get participants: ok1",
			id:                   tournaments[1].ID,
			expectedParticipants: expParticipants[1],
			expectedError:        nil,
		},
		{
			name:                 "get participants: ok2",
			id:                   tournaments[2].ID,
			expectedParticipants: expParticipants[2],
			expectedError:        nil,
		},
		{
			name:                 "get participants: fail",
			id:                   "getparts_fake",
			expectedParticipants: nil,
			expectedError:        errors.Error{Code: errors.NotFoundError, Message: "get participants: cannot get participants from not existing tournament, id: getparts_fake"},
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
		{ID: "getdeposit_1", Deposit: 100},
		{ID: "getdeposit_2", Deposit: 200},
		{ID: "getdeposit_3", Deposit: 300},
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
			name:            "get deposit: ok0",
			id:              tournaments[0].ID,
			expectedDeposit: tournaments[0].Deposit,
			expectedError:   nil,
		},
		{
			name:            "get deposit: ok1",
			id:              tournaments[1].ID,
			expectedDeposit: tournaments[1].Deposit,
			expectedError:   nil,
		},
		{
			name:            "get deposit: ok2",
			id:              tournaments[2].ID,
			expectedDeposit: tournaments[2].Deposit,
			expectedError:   nil,
		},
		{
			name:            "get deposit: fake id",
			id:              "getdeposit_fake",
			expectedDeposit: 0,
			expectedError:   errors.Error{Code: errors.NotFoundError, Message: "get deposit: cannot get deposit from not existing tournament, id: getdeposit_fake"},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			dep, err := p.getDeposit(tc.id)
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedDeposit, dep)
		})
	}
}

func TestTournament_GetState(t *testing.T) {
	tournaments := []entity.Tournament{
		{ID: "getstate_1", Deposit: 100},
		{ID: "getstate_2", Deposit: 200},
		{ID: "getstate_3", Deposit: 300},
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
			name:          "get state: ok0",
			id:            tournaments[0].ID,
			expectedState: states[0],
			expectedError: nil,
		},
		{
			name:          "get state: ok1",
			id:            tournaments[1].ID,
			expectedState: states[1],
			expectedError: nil,
		},
		{
			name:          "get state: ok2",
			id:            tournaments[2].ID,
			expectedState: states[2],
			expectedError: nil,
		},
		{
			name:          "get state: fake id",
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
		{ID: "getwinner_1", Deposit: 100},
		{ID: "getwinner_2", Deposit: 100},
		{ID: "getwinner_3", Deposit: 100},
	}
	winners := []entity.Winner{
		{ID: "getwinner_1", Points: 50},
		{ID: "getwinner_2", Points: 200},
		{ID: "getwinner_3", Points: 200},
		{ID: "getwinner_4", Points: 50},
		{ID: "getwinner_5", Points: 200},
		{ID: "getwinner_6", Points: 50},
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
			name:           "get winner: ok0",
			id:             tournaments[0].ID,
			expectedWinner: entity.Winners{Winners: []entity.Winner{expWinner[0]}},
			expectedError:  nil,
		},
		{
			name:           "get winner: ok1",
			id:             tournaments[1].ID,
			expectedWinner: entity.Winners{Winners: []entity.Winner{expWinner[1]}},
			expectedError:  nil,
		},
		{
			name:           "get winner: ok2",
			id:             tournaments[2].ID,
			expectedWinner: entity.Winners{Winners: []entity.Winner{expWinner[2]}},
			expectedError:  nil,
		},
		{
			name:           "get winner: fail",
			id:             "getwinner_fake",
			expectedWinner: entity.Winners{},
			expectedError:  errors.Error{Code: errors.NotFoundError, Message: "get winner: cannot get winner from not existing tournament, id: getwinner_fake"},
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
		{ID: "setwinner_1", Deposit: 100},
		{ID: "setwinner_2", Deposit: 100},
		{ID: "setwinner_3", Deposit: 100},
	}
	winners := []entity.Winner{
		{ID: "setwinner_1", Points: 50},
		{ID: "setwinner_2", Points: 200},
		{ID: "setwinner_3", Points: 200},
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
			name:           "set winner: ok0",
			id:             tournaments[0].ID,
			winner:         winners[0],
			expectedWinner: entity.Winners{Winners: []entity.Winner{winners[0]}},
			expectedError:  nil,
		},
		{
			name:           "set winner: ok1",
			id:             tournaments[1].ID,
			winner:         winners[1],
			expectedWinner: entity.Winners{Winners: []entity.Winner{winners[1]}},
			expectedError:  nil,
		},
		{
			name:           "set winner: ok2",
			id:             tournaments[2].ID,
			winner:         winners[2],
			expectedWinner: entity.Winners{Winners: []entity.Winner{winners[2]}},
			expectedError:  nil,
		},
		{
			name:           "set winner: fail",
			id:             "setwinner_fake",
			winner:         entity.Winner{},
			expectedWinner: entity.Winners{},
			expectedError:  errors.Error{Code: errors.NotFoundError, Message: "set winner: tournament not exist, id: setwinner_fake\nsql: no rows in result set"},
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
		{ID: "updategame_1", Deposit: 50},
		{ID: "updategame_2", Deposit: 50},
		{ID: "updategame_3", Deposit: 50},
	}
	players := []entity.Player{
		{ID: "updategame_1", Points: 50},
		{ID: "updategame_2", Points: 100},
		{ID: "updategame_3", Points: 150},
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
			name:                 "update tournament and player: tour 1",
			id:                   tournaments[0].ID,
			participants:         players,
			expectedPoints:       []int{0, 50, 100},
			expectedParticipants: []string{players[0].ID, players[1].ID, players[2].ID},
			expectedUpdateErrors: []error{nil, nil, nil},
			expectedGetPartError: nil,
			expectedGetPlayError: []error{nil, nil, nil},
		},
		{
			name:                 "update tournament and player: tour 2",
			id:                   tournaments[1].ID,
			participants:         players,
			expectedPoints:       []int{0, 0, 50},
			expectedParticipants: []string{players[1].ID, players[2].ID},
			expectedUpdateErrors: []error{errors.Error{Code: errors.NegativePointsNumberError, Message: "update player: cannot update points numbers, dif -50"}, nil, nil},
			expectedGetPartError: nil,
			expectedGetPlayError: []error{nil, nil, nil},
		},
		{
			name:                 "update tournament and player: tour 3",
			id:                   tournaments[2].ID,
			participants:         players,
			expectedPoints:       []int{0, 0, 0},
			expectedParticipants: []string{players[2].ID},
			expectedUpdateErrors: []error{errors.Error{Code: errors.NegativePointsNumberError, Message: "update player: cannot update points numbers, dif -50"}, errors.Error{Code: errors.NegativePointsNumberError, Message: "update player: cannot update points numbers, dif -50"}, nil},
			expectedGetPartError: nil,
			expectedGetPlayError: []error{nil, nil, nil},
		},
		{
			name:                 "update tournament and player: fake tour",
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

// Package handlers contains handlers to handle application endpoints
package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/dmitriyomelyusik/Tournament/entity"
	"github.com/dmitriyomelyusik/Tournament/errors"
	"github.com/gorilla/mux"
)

type ctlr interface {
	Fund(id string, points int) (entity.Player, error)
	Take(id string, points int) error
	Balance(id string) (entity.Player, error)
	AnnounceTournament(id string, deposit int) error
	JoinTournament(tourID, playerID string) error
	Results(tourID string) (entity.Winners, error)
}

// Server uses controller in handling http methods
type Server struct {
	Controller ctlr
}

// HandleFund handles fund query
func (s Server) HandleFund() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		id := query.Get("playerId")
		points := query.Get("points")
		p, err := strconv.Atoi(points)
		if err != nil {
			jsonError(w, errors.Error{Code: errors.NotNumberError, Message: "cannot fund player, points is not number: " + points, Info: err.Error()})
			return
		}
		player, err := s.Controller.Fund(id, p)
		if err != nil {
			jsonError(w, err)
			return
		}
		if player != (entity.Player{}) {
			jsonResponse(w, player, http.StatusCreated)
			return
		}
	}
}

// HandleTake handles take query
func (s Server) HandleTake() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		id := query.Get("playerId")
		points := query.Get("points")
		p, err := strconv.Atoi(points)
		if err != nil {
			jsonError(w, errors.Error{Code: errors.NotNumberError, Message: "cannot take points, points is not number: " + points, Info: err.Error()})
			return
		}
		err = s.Controller.Take(id, p)
		if err != nil {
			jsonError(w, err)
			return
		}
	}
}

// HandleBalance handles balance query
func (s Server) HandleBalance() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("playerId")
		p, err := s.Controller.Balance(id)
		if err != nil {
			jsonError(w, err)
			return
		}
		jsonResponse(w, p, http.StatusOK)
	}
}

// HandleAnnounce handles announce query
func (s Server) HandleAnnounce() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		id := query.Get("tournamentId")
		dep := query.Get("deposit")
		deposit, err := strconv.Atoi(dep)
		if err != nil {
			jsonError(w, errors.Error{Code: errors.NotNumberError, Message: "cannot create tournament, deposit is not number: " + dep, Info: err.Error()})
			return
		}
		err = s.Controller.AnnounceTournament(id, deposit)
		if err != nil {
			jsonError(w, err)
			return
		}
	}
}

// HandleJoin handles join query
func (s Server) HandleJoin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		tourID := query.Get("tournamentId")
		playerID := query.Get("playerId")
		err := s.Controller.JoinTournament(tourID, playerID)
		if err != nil {
			jsonError(w, err)
			return
		}
	}
}

//HandleResults handles results query
func (s Server) HandleResults() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tourID := r.URL.Query().Get("tournamentId")
		res, err := s.Controller.Results(tourID)
		if err != nil {
			jsonError(w, err)
			return
		}
		jsonResponse(w, res, http.StatusOK)
	}
}

// NewRouter returns router with configurated and handled pathes
func NewRouter(s Server) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/fund", s.HandleFund())
	r.HandleFunc("/take", s.HandleTake())
	r.HandleFunc("/balance", s.HandleBalance())
	r.HandleFunc("/announceTournament", s.HandleAnnounce())
	r.HandleFunc("/joinTournament", s.HandleJoin())
	r.HandleFunc("/resultTournament", s.HandleResults())
	return r
}

func jsonError(w http.ResponseWriter, err error) {
	myErr, ok := err.(errors.Error)
	if !ok {
		myErr = errors.Error{
			Code:    "UnknownError",
			Message: err.Error(),
		}
	}
	var status int
	switch myErr.Code {
	case errors.NotFoundError, errors.NotNumberError, errors.NegativePointsNumberError, errors.NegativeDepositError, errors.DuplicatedIDError, errors.ClosedTournamentError:
		status = http.StatusNotFound
	case errors.NoneParticipantsError:
		status = http.StatusOK
	default:
		status = http.StatusInternalServerError
	}
	jsonResponse(w, myErr, status)
}

func jsonResponse(w http.ResponseWriter, data interface{}, status int) {
	w.WriteHeader(status)
	w.Header().Set("content-type", "application/json")
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(data); err != nil {
		log.Println(err)
	}
}

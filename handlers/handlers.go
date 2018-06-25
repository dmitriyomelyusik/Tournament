package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/Tournament/controller"
	"github.com/Tournament/entity"
	"github.com/Tournament/errors"
	"github.com/gorilla/mux"
)

// Server uses controller in handling http methods
type Server struct {
	Controller controller.Game
}

// HandleFund handles fund query
func (s Server) HandleFund() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("playerId")
		points := r.URL.Query().Get("points")
		p, err := strconv.Atoi(points)
		if err != nil {
			jsonError(w, err)
			return
		}
		player, err := s.Controller.Fund(id, p)
		if err != nil {
			jsonError(w, err)
			return
		}
		if player != (entity.Player{}) {
			w.WriteHeader(http.StatusCreated)
			jsonResponse(w, player)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

// HandleTake handles take query
func (s Server) HandleTake() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("playerId")
		points := r.URL.Query().Get("points")
		p, err := strconv.Atoi(points)
		if err != nil {
			jsonError(w, err)
			return
		}
		err = s.Controller.Take(id, p)
		if err != nil {
			jsonError(w, err)
			return
		}
		w.WriteHeader(http.StatusOK)
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
		jsonResponse(w, p)
	}
}

// HandleAnnounce handles announce query
func (s Server) HandleAnnounce() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("tournamentId")
		dep := r.URL.Query().Get("deposit")
		deposit, err := strconv.Atoi(dep)
		if err != nil {
			jsonError(w, err)
			return
		}
		err = s.Controller.AnnounceTournament(id, deposit)
		if err != nil {
			jsonError(w, err)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

// HandleJoin handles join query
func (s Server) HandleJoin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tourID := r.URL.Query().Get("tournamentId")
		playerID := r.URL.Query().Get("playerId")
		err := s.Controller.JoinTournament(tourID, playerID)
		if err != nil {
			jsonError(w, err)
			return
		}
		w.WriteHeader(http.StatusOK)
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
		jsonResponse(w, res)
	}
}

// HandleDeleteTour handles deleting tournament
func (s Server) HandleDeleteTour() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tourID := r.URL.Query().Get("tournamentId")
		err := s.Controller.DeleteTournament(tourID)
		if err != nil {
			jsonError(w, err)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

// HandleDeletePlayer handles deleting tournament
func (s Server) HandleDeletePlayer() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tourID := r.URL.Query().Get("playerId")
		err := s.Controller.DeletePlayer(tourID)
		if err != nil {
			jsonError(w, err)
			return
		}
		w.WriteHeader(http.StatusOK)
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
	r.HandleFunc("/deleteTournament", s.HandleDeleteTour())
	r.HandleFunc("/deletePlayer", s.HandleDeletePlayer())
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
	switch myErr.Code {
	case errors.NotFoundError, errors.NegativePointsNumberError, errors.NegativeDepositError, errors.DuplicatedIDError, errors.ClosedTournamentError:
		w.WriteHeader(http.StatusNotFound)
	case errors.NoneParticipantsError:
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	jsonResponse(w, myErr)
}

func jsonResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("content-type", "application/json")
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(data); err != nil {
		log.Println(err)
	}
}

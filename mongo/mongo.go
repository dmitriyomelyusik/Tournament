package mongo

import (
	"github.com/dmitriyomelyusik/Tournament/mongo/logs"
	mgo "gopkg.in/mgo.v2"
)

// Mongo is an implementation of needed mongodb
type Mongo struct {
	s           *mgo.Session
	db          *mgo.Database
	players     *mgo.Collection
	tournaments *mgo.Collection
	logger      *logger.Logger
}

// NewDB returns mongo database with configuration conf
func NewDB(conf string) (*Mongo, error) {
	s, err := mgo.Dial(conf)
	if err != nil {
		return nil, err
	}
	db := s.DB("mongo")
	players := db.C("players")
	tournaments := db.C("tournaments")
	log := &logger.Logger{Logger: db.C("logger")}
	return &Mongo{s, db, players, tournaments, log}, nil
}

// Close closes database connection
func (m *Mongo) Close() {
	m.s.Close()
	return
}

// Ping runs a trivial ping command just to get in touch with the server.
func (m *Mongo) Ping() error {
	return m.s.Ping()
}

// UpdateTourAndPlayer updates tournament participants and player balance
func (m *Mongo) UpdateTourAndPlayer(tourID string, playerID string) error {
	return nil
}

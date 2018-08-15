package logger

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Data is a data that is stored in log
type Data struct {
	ID     string `bson:"id"`
	Op     string `bson:"operation"`
	Points int    `bson:"points"`
}

// Block of available operations
const (
	Take = "take"
	Fund = "fund"
	Won  = "won"
)

// Logger is collection that logs all operations with players
type Logger struct {
	Logger *mgo.Collection
}

// Log logs operation
func (l *Logger) Log(id, op string, points int) error {
	return l.Logger.Insert(Data{ID: id, Op: op, Points: points})
}

// GetLogs returns all operations, that have been done with player
func (l *Logger) GetLogs(id string) ([]Data, error) {
	var d []Data
	err := l.Logger.Find(bson.M{"id": id}).All(d)
	return d, err
}

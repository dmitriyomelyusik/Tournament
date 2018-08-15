package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dmitriyomelyusik/Tournament/controller"
	"github.com/dmitriyomelyusik/Tournament/errors"
	"github.com/dmitriyomelyusik/Tournament/handlers"
	"github.com/dmitriyomelyusik/Tournament/mongo"
	"github.com/dmitriyomelyusik/Tournament/postgres"
	_ "github.com/lib/pq"
)

// Environment variables that needs to open database
const (
	DBNAME   = "DBNAME"
	PGUSER   = "PGUSER"
	PGPASS   = "PGPASS"
	PGHOST   = "PGHOST"
	SSLMODE  = "SSLMODE"
	DBDRIVER = "DBDRIVER"
)

func main() {
	var (
		db  controller.Database
		err error
	)
	switch os.Getenv(DBDRIVER) {
	case "mongo":
		db, err = getMongo()
		if err != nil {
			log.Fatalln(err)
		}
	case "postgres":
		db, err = getPostgres()
		if err != nil {
			log.Fatalln(err)
		}
	default:
		panic("You didn't set DBDRIVER variable.")
	}

	ctl := controller.Game{DB: db}
	server := handlers.Server{Controller: ctl}
	r := handlers.NewRouter(server)
	s := http.Server{
		Addr:         ":8080",
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
		Handler:      r,
	}
	err = s.ListenAndServe()
	if err != nil {
		log.Fatal(errors.Error{Code: errors.ConnectionError, Message: "listen and serve: error occured", Info: err.Error()})
	}
}

func getMongo() (*mongo.Mongo, error) {
	m, err := mongo.NewDB("localhost")
	if err != nil {
		return m, errors.Error{Code: errors.DatabaseCreatingError, Message: "creating database: cannot create database", Info: err.Error()}
	}
	err = m.Ping()
	if err != nil {
		return nil, errors.Error{Code: errors.DatabasePingError, Message: "ping database: fatal to ping database", Info: err.Error()}
	}
	return m, nil
}

func getPostgres() (*postgres.Postgres, error) {
	vars := getEnvVars()
	conf := fmt.Sprintf("user=%v dbname=%v password=%v sslmode=%v host=%v", vars[PGUSER], vars[DBNAME], vars[PGPASS], vars[SSLMODE], vars[PGHOST])
	p, err := postgres.NewDB(conf)
	if err != nil {
		return nil, errors.Error{Code: errors.DatabaseCreatingError, Message: "creating database: cannot create database", Info: err.Error()}
	}
	err = p.Ping()
	if err != nil {
		return nil, errors.Error{Code: errors.DatabasePingError, Message: "ping database: fatal to ping database", Info: err.Error()}
	}
	return p, nil
}

func getEnvVars() map[string]string {
	vars := make(map[string]string)
	vars[PGUSER] = os.Getenv(PGUSER)
	vars[PGPASS] = os.Getenv(PGPASS)
	vars[DBNAME] = os.Getenv(DBNAME)
	vars[PGHOST] = os.Getenv(PGHOST)
	vars[SSLMODE] = os.Getenv(SSLMODE)
	var lostVars bool
	for key, value := range vars {
		if value == "" {
			fmt.Printf("You didn't set environment variable: %v\n", key)
			lostVars = true
		}
	}
	if lostVars {
		panic("Set all environment variables!")
	}
	return vars
}

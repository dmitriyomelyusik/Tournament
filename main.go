package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Tournament/controller"
	"github.com/Tournament/handlers"
	"github.com/Tournament/postgres"
	_ "github.com/lib/pq"
)

// Environment variables that needs to open database
const (
	DBNAME  = "DBNAME"
	PGUSER  = "PGUSER"
	PGPASS  = "PGPASS"
	PGHOST  = "PGHOST"
	SSLMODE = "SSLMODE"
)

func main() {
	vars := getEnvVars()
	conf := fmt.Sprintf("user=%v dbname=%v password=%v sslmode=%v host=%v", vars[PGUSER], vars[DBNAME], vars[PGPASS], vars[SSLMODE], vars[PGHOST])
	p, err := postgres.NewDB(conf)
	if err != nil {
		log.Fatal(err)
	}
	err = p.DB.Ping()
	if err != nil {
		log.Fatal(err)
	}
	ctl := controller.Game{DB: p}
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
		log.Fatal(err)
	}
}

func getEnvVars() map[string]string {
	vars := make(map[string]string)
	vars[PGUSER] = os.Getenv(PGUSER)
	vars[PGPASS] = os.Getenv(PGPASS)
	vars[DBNAME] = os.Getenv(DBNAME)
	vars[PGHOST] = os.Getenv(PGHOST)
	vars[SSLMODE] = os.Getenv(SSLMODE)
	isOK := true
	for key, value := range vars {
		if value == "" {
			fmt.Printf("You didn't set environment variable: %v\n", key)
			isOK = false
		}
	}
	if !isOK {
		panic("Set all environment variables!")
	}
	return vars
}

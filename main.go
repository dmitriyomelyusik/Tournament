package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Tournament/controller"
	"github.com/Tournament/handlers"
	"github.com/Tournament/postgres"
	_ "github.com/lib/pq"
)

func main() {
	p, err := postgres.NewDB("user=postgres dbname=postgres password=password sslmode=disable host=127.0.0.1")
	if err != nil {
		log.Fatal(err)
	}
	err = p.DB.Ping()
	if err != nil {
		log.Fatal(err)
	}
	ctl := controller.Game{TDB: p, PDB: p}
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

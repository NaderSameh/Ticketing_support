package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/naderSameh/ticketing_support/api"
	db "github.com/naderSameh/ticketing_support/db/sqlc"
)

const (
	DBDriver      = "postgres"
	DBSource      = "postgres://nader:nader123@localhost:5432/ticketing_support?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {

	conn, err := sql.Open(DBDriver, DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	server.Start(serverAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}

}

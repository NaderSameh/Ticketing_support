package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/naderSameh/ticketing_support/api"
	db "github.com/naderSameh/ticketing_support/db/sqlc"
	"github.com/naderSameh/ticketing_support/util"
	"github.com/spf13/viper"
)

func main() {

	err := util.Loadconfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := sql.Open(viper.GetString("DB_DRIVER"), viper.GetString("DB_SOURCE"))
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	server.Start(viper.GetString("SERVER_ADDRESS"))
	if err != nil {
		log.Fatal("cannot start server:", err)
	}

}

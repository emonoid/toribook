package main

import (
	"database/sql"
	"log"

	"github.com/emonoid/toribook.git/api"
	db "github.com/emonoid/toribook.git/db/sqlc"
	"github.com/emonoid/toribook.git/utils"
	_ "github.com/lib/pq" // PostgreSQL driver
)

func main(){
	config, err := utils.LoadConfig(".")

	if err != nil {
		log.Fatal("Cannot load config", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("Cannot create server: ",err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("Cannot start server:", err)
	}
}

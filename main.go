package main

import (
	"database/sql"
	"log"

	"github.com/arywr/od-reconciliation-api/api"
	db "github.com/arywr/od-reconciliation-api/db/sqlc"
	"github.com/arywr/od-reconciliation-api/util"
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot load config: ", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Database connection failed: ", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("Cannot starting to server: ", err)
	}
}

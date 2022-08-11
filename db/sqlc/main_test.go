package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/arywr/od-reconciliation-api/util"
	_ "github.com/lib/pq"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("Cannot load config: ", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Database connection failed!")
	}

	testQueries = New(conn)

	os.Exit(m.Run())
}

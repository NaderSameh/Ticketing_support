package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var testQueries *Queries

const (
	DBDriver = "postgres"
	DBSource = "postgres://nader:nader123@localhost:5432/ticketing_support?sslmode=disable"
)

func TestMain(m *testing.M) {

	testDB, err := sql.Open(DBDriver, DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}

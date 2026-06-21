package db

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var TestQueries *Queries
var TestDB *sql.DB

func TestMain(m *testing.M) {

	conn, err := sql.Open("postgres", "postgresql://root:mypassword@localhost:5558/simple_bank?sslmode=disable")
	if err != nil {
		panic(err)
	}
	TestDB = conn
	TestQueries = New(conn)

	os.Exit(m.Run())
}

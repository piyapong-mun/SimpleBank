package db

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/piyapong-mun/simplebank/util"
)

var TestQueries *Queries
var TestDB *sql.DB

func TestMain(m *testing.M) {

	config, err := util.LoadConfig("./../..")
	if err != nil {
		panic(err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		panic(err)
	}

	TestDB = conn
	TestQueries = New(conn)

	os.Exit(m.Run())
}

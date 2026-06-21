package main

import (
	// "context"
	"database/sql"
	"log"

	// Import your generated db package (replace with your actual module name)

	// Import the postgres driver
	_ "github.com/lib/pq"

	util "github.com/piyapong-mun/simplebank/util"
)

func main() {
	// ctx := context.Background()

	// 1. Connect to your local Docker Postgres database
	connStr := "postgresql://root:mypassword@localhost:5558/simple_bank?sslmode=disable"
	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Cannot connect to database: %v", err)
	}
	defer conn.Close()

	// 2. Initialize the sqlc generated Queries struct
	// queries := db.New(conn)

	println(util.RandomString(15))
	println(util.RandomCurrency())
	println(util.RandomNumber(1, 100))

}

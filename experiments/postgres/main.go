package main

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	connectionString := "postgres://postgres:shadow@localhost:5432/db_crusher"
	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		fmt.Println("Something went wrong ", err)
		return
	}

	err = db.Ping()

	if err != nil {
		fmt.Print("Failed to establish Connection ", err)
	}

	//err = db.Close()
	//if err != nil {
	//	fmt.Print("Failed to Close DB Connection", err)
	//}
}

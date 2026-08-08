package database

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type Database struct {
	DB *sql.DB
}

func NewDB() (*Database, error) {
	connectionString := "postgres://postgres:shadow@localhost:5432/db_crusher"
	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &Database{
		DB: db,
	}, nil
}

func ViewData() {

	res, err := db.Query("SELECT id, name, age FROM users")
	if err != nil {
		fmt.Println("Query Error ", err)
		return
	}

	defer res.Close()
	for res.Next() {
		var id int
		var name string
		var age int
		err := res.Scan(&id, &name, &age)

		if err != nil {
			fmt.Println("Failed to fetch Data! ", err)
			return
		}

		fmt.Println(id, name, age)
	}
}

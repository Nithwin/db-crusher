package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type AnalyticsResult struct {
	Status        string  `json:"status"`
	Count         int64   `json:"count"`
	TotalAmount   float64 `json:"total_amount"`
	AverageAmount float64 `json:"average_amount"`
}

type Database struct {
	DB *sql.DB
}

func NewDB() (*Database, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")
	connectionString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		user,
		password,
		host,
		port,
		name,
	)
	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}
	return &Database{
		DB: db,
	}, nil
}

func (d *Database) GetUsers() ([]User, error) {

	rows, err := d.DB.Query(`SELECT id, name, age FROM users`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []User

	for rows.Next() {
		var user User

		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Age,
		)

		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (d *Database) GetUser(id int) (User, error) {
	var user User

	err := d.DB.QueryRow(`
		SELECT id, name, age
		FROM users
		WHERE id = $1
	`, id).Scan(
		&user.ID,
		&user.Name,
		&user.Age,
	)

	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (d *Database) CreateUser(name string, age int) error {
	_, err := d.DB.Exec(`
		INSERT INTO users (name, age)
		VALUES ($1, $2)
	`, name, age)

	return err
}

func (d *Database) DeleteUser(id int) error {
	_, err := d.DB.Exec(`
		DELETE FROM users
		WHERE id = $1
	`, id)

	return err
}

func (d *Database) GetAnalytics() ([]AnalyticsResult, error) {
	rows, err := d.DB.Query(`
		SELECT
			status,
			COUNT(*),
			SUM(amount),
			AVG(amount)
		FROM transactions
		GROUP BY status
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []AnalyticsResult

	for rows.Next() {
		var res AnalyticsResult

		if err := rows.Scan(
			&res.Status,
			&res.Count,
			&res.TotalAmount,
			&res.AverageAmount,
		); err != nil {
			return nil, err
		}

		results = append(results, res)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

package main

import (
	"db-crusher/internal/database"
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Failed to load .env:", err)
	}

	db, err := database.NewDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.DB.Close()

	createTable := `
	CREATE TABLE IF NOT EXISTS transactions (
		id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
		user_id INT NOT NULL,
		amount NUMERIC(10, 2) NOT NULL,
		status VARCHAR(20) NOT NULL,
		created_at TIMESTAMP NOT NULL
	);
	`

	if _, err := db.DB.Exec(createTable); err != nil {
		log.Fatal("Failed to create transactions table:", err)
	}

	fmt.Println("transactions table ready")

	insertData := `
	INSERT INTO transactions (user_id, amount, status, created_at)
	SELECT
		(1 + floor(random() * 100000))::INT,
		round((10 + random() * 9990)::numeric, 2),
		(
			ARRAY['completed', 'pending', 'failed', 'refunded']
		)[1 + floor(random() * 4)::INT],
		TIMESTAMP '2025-01-01'
			+ random() * (TIMESTAMP '2026-01-01' - TIMESTAMP '2025-01-01')
	FROM generate_series(1, 1000000);
	`

	fmt.Println("Generating 1,000,000 transactions...")

	if _, err := db.DB.Exec(insertData); err != nil {
		log.Fatal("Failed to insert transactions:", err)
	}

	fmt.Println("Successfully inserted 1,000,000 transactions")
}

package main

import (
	"db-crusher/internal/database"
	"db-crusher/internal/handlers"
	"fmt"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env:", err)
		return
	}

	db, err := database.NewDB()
	if err != nil {
		fmt.Println("Database connection failed:", err)
		return
	}
	defer db.DB.Close()

	userHandler := &handlers.UserHandler{
		DB: db,
	}
	http.HandleFunc("GET /health", userHandler.HealthHandler)
	http.HandleFunc("GET /users", userHandler.GetUsersHandler)
	http.HandleFunc("GET /users/{id}", userHandler.GetUserHandler)
	http.HandleFunc("POST /users", userHandler.CreateUserHandler)
	http.HandleFunc("DELETE /users/{id}", userHandler.DeleteUserHandler)

	http.HandleFunc("GET /analytics", userHandler.GetAnalytics)

	fmt.Println("Your server is running in http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Server error! ", err)
	}

}

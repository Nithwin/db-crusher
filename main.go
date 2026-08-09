package main

import (
	"context"
	"db-crusher/internal/cache"
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

	redisClient, err := cache.NewRedisClient(context.Background())
	if err != nil {
		fmt.Println("Redis connection failed:", err)
		return
	}
	defer redisClient.Client.Close()

	fmt.Println("Redis connected successfully")

	db, err := database.NewDB()
	if err != nil {
		fmt.Println("Database connection failed:", err)
		return
	}
	defer db.DB.Close()

	userHandler := &handlers.UserHandler{
		DB:    db,
		Cache: redisClient,
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

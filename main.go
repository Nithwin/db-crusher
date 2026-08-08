package main

import (
	"db-crusher/internal/database"
	"encoding/json"
	"fmt"
	"net/http"
)

type HealthResponse struct {
	Status string `json:"status"`
}

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	http.HandleFunc("GET /health", healthHandler)
	http.HandleFunc("GET /users", getUserHandler)
	http.HandleFunc("POST /users", createUserHandler)
	fmt.Println("Your server is running in http://localhost:8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Server error! ", err)
	}

}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	status := HealthResponse{
		Status: "OK",
	}
	err := json.NewEncoder(w).Encode(status)

	if err != nil {
		fmt.Println("Something went Wrong")
	}
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	user := User{
		ID:   1,
		Name: "Nithwin",
		Age:  21,
	}
	database.ViewData()
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(user); err != nil {
		fmt.Println("Failed to Fetch! ", err)
	}
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	var req User
	err := json.NewDecoder(r.Body).Decode(&req)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		fmt.Println("Can't parse the data ", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Invalid Request"})
		return
	}
	fmt.Println(req)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User Created"})
}

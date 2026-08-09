package handlers

import (
	"db-crusher/internal/cache"
	"db-crusher/internal/database"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type UserHandler struct {
	DB    *database.Database
	Cache *cache.RedisClient
}

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type HealthResponse struct {
	Status string `json:"status"`
}

type CreateUserRequest struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func (h *UserHandler) HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	status := HealthResponse{
		Status: "OK",
	}
	err := json.NewEncoder(w).Encode(status)

	if err != nil {
		fmt.Println("Something went Wrong")
	}
}

func (h *UserHandler) GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	users, err := h.DB.GetUsers()

	if err != nil {
		fmt.Println("Failed to Fetch! ", err)
	}

	if err := json.NewEncoder(w).Encode(users); err != nil {
		fmt.Println("Failed to Fetch! ", err)
	}
}

func (h *UserHandler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		fmt.Println("Can't parse the data ", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Invalid Request"})
		return
	}
	user, err := h.DB.GetUser(id)

	if err != nil {
		fmt.Println("Failed to Fetch! ", err)
	}

	if err := json.NewEncoder(w).Encode(user); err != nil {
		fmt.Println("Failed to Fetch! ", err)
	}
}

func (h *UserHandler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req CreateUserRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		fmt.Println("Can't parse the data ", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Invalid Request"})
		return
	}

	err = h.DB.CreateUser(req.Name, req.Age)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Failed to Create"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User Created"})
}

func (h *UserHandler) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		fmt.Println("Can't parse the data ", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Invalid Request"})
		return
	}
	err = h.DB.DeleteUser(id)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Failed to Deleted"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode(map[string]string{"message": "User Deleted"})
}

func (h *UserHandler) GetAnalytics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := r.Context()

	cached, err := h.Cache.Get(ctx, "analytics")

	if err == nil {
		fmt.Println("Cache HIT")
		var analytics []database.AnalyticsResult
		if err := json.Unmarshal([]byte(cached), &analytics); err != nil {
			fmt.Println("cache contains invalid data ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(analytics)
		return
	}

	if err != redis.Nil {
		fmt.Println("Redis error:", err)
	} else {
		fmt.Println("CACHE MISS")
	}

	analytics, err := h.DB.GetAnalytics()

	if err != nil {
		fmt.Println("Failed to Fetch! ", err)
		return
	}

	data, err := json.Marshal(analytics)

	if err != nil {
		fmt.Println("Failed to Encode ")
		return
	}

	err = h.Cache.Set(ctx,
		"analytics",
		string(data),
		10*time.Second)

	if err != nil {
		fmt.Println("Failed to set Cache")
		return
	}

	if err := json.NewEncoder(w).Encode(analytics); err != nil {
		fmt.Println("Failed to Fetch! ", err)
	}
}

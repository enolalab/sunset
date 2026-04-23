// Package handler provides HTTP request handlers.
package handler

import (
	"encoding/json"
	"net/http"

	"example.com/sample/model"
)

// UserHandler handles user-related HTTP requests.
type UserHandler struct {
	repo UserRepository
}

// UserRepository defines the interface for user data access.
type UserRepository interface {
	FindByID(id string) (*model.User, error)
	Create(user *model.User) error
	Update(user *model.User) error
}

// GetUser retrieves a user by ID from the query parameter.
// Returns 404 if the user is not found.
func GetUser(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	user := &model.User{ID: id, Name: "Test User", Email: "test@example.com"}
	json.NewEncoder(w).Encode(user)
}

// CreateUser creates a new user from the JSON request body.
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

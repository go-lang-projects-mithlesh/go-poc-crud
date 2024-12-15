package main

import (
	"encoding/json"
	"net/http"
)

// User represents a user in the system
type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// A mock database of users
var users = []User{
	{Name: "John Doe", Age: 30},
	{Name: "Jane Smith", Age: 25},
}

// GetUsers godoc
// @Summary Get Users
// @Description Retrieve a list of all users
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {array} User
// @Router /users [get]
func GetUsers(w http.ResponseWriter, r *http.Request) {
	// Return the list of users in JSON format
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// createUser godoc
// @Summary Create User
// @Description Create a new user
// @Tags users
// @Accept json
// @Produce json
// @Param user body User true "User data"
// @Success 201 {object} User
// @Router /create [post]
func createUser(w http.ResponseWriter, r *http.Request) {
	var newUser User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, "Invalid user data", http.StatusBadRequest)
		return
	}
	users = append(users, newUser)
	w.Write([]byte("User created successfully!"))
}

func main() {
	http.HandleFunc("/users", GetUsers)    // Endpoint for getting users
	http.HandleFunc("/create", createUser) // Endpoint for creating a user

	http.ListenAndServe(":8080", nil) // Start the server on port 8080
}

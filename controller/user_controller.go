package controller

import (
	"go-poc-crud/repository"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var tmpl *template.Template

func InitializeTemplate() {
	var err error
	tmpl, err = template.ParseFiles("views/users.html")
	if err != nil {
		log.Fatalf("Unable to load template 'views/users.html': %v", err)
	}
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := repository.GetAllUsers()
	if err != nil {
		log.Printf("Error fetching users from the database: %v", err)
		http.Error(w, "Unable to fetch users", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, users)
	if err != nil {
		log.Printf("Error executing template 'users.html': %v", err)
		http.Error(w, "Unable to execute template", http.StatusInternalServerError)
	}
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	firstName := r.FormValue("first_name")
	lastName := r.FormValue("last_name")
	email := r.FormValue("email")

	err := repository.CreateUser(firstName, lastName, email)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		http.Error(w, "Unable to create user", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" && r.FormValue("_method") == "DELETE" {
		vars := mux.Vars(r)
		userIDStr := vars["id"]
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			log.Printf("Invalid user ID: %v", err)
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		err = repository.DeleteUser(userID)
		if err != nil {
			log.Printf("Error deleting user: %v", err)
			http.Error(w, "Unable to delete user", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

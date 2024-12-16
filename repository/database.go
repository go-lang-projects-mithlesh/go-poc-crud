package repository

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite" // SQLite driver
)

// Declare DB at the package level
var DB *sql.DB

// User represents a user record
type User struct {
	ID        int
	FirstName string
	LastName  string
	Email     string
}

// InitDatabase initializes the SQLite database
func InitDatabase() error {
	var err error
	// Open the database connection and assign to the global DB variable
	DB, err = sql.Open("sqlite", "./users.db")
	if err != nil {
		return err
	}

	// Ping the database to verify it's working
	if err := DB.Ping(); err != nil {
		return err
	}
	log.Println("Database connection established.")

	// Create table if it doesn't exist
	createTable := `CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        first_name TEXT NOT NULL,
        last_name TEXT NOT NULL,
        email TEXT NOT NULL UNIQUE
    );`
	_, err = DB.Exec(createTable)
	if err != nil {
		return err
	}
	log.Println("Users table is ready.")

	return nil
}

// GetAllUsers retrieves all users from the database
func GetAllUsers() ([]User, error) {
	// Query to fetch all users
	rows, err := DB.Query("SELECT id, first_name, last_name, email FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	// Iterate through rows and build the users slice
	for rows.Next() {
		var user User
		err = rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	// Check for errors during iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// CreateUser handles the creation of a new user in the database
func CreateUser(firstName, lastName, email string) error {
	// Insert the new user into the database
	_, err := DB.Exec("INSERT INTO users (first_name, last_name, email) VALUES (?, ?, ?)", firstName, lastName, email)
	if err != nil {
		return err
	}
	return nil
}

// DeleteUser handles the deletion of a user from the database
func DeleteUser(userID int) error {
	// Delete the user from the database using the provided ID
	_, err := DB.Exec("DELETE FROM users WHERE id = ?", userID)
	if err != nil {
		return err
	}
	return nil
}

// CloseDatabase safely closes the DB connection
func CloseDatabase() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

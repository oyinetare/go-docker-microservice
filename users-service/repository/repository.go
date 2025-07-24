package repository

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
}

type Repository struct {
	db *sql.DB
}

func Connect(host, database, user, password string, port int) (*Repository, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", user, password, host, port, database)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &Repository{db: db}, nil
}

func (r *Repository) GetUsers() ([]User, error) {
	rows, err := r.db.Query("SELECT email, phone_number FROM directory")
	if err != nil {
		return nil, fmt.Errorf("error getting users: %w", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Email, &user.PhoneNumber); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *Repository) GetUserByEmail(email string) (*User, error) {
	var user User
	err := r.db.QueryRow("SELECT email, phone_number FROM directory WHERE email = ?", email).
		Scan(&user.Email, &user.PhoneNumber)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error getting user: %w", err)
	}

	return &user, nil
}

func (r *Repository) Disconnect() error {
	return r.db.Close()
}

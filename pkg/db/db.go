package db

import (
	"context"
	"database/sql"
	"errors"

	_ "github.com/lib/pq"
)

var ErrUserNotFound = errors.New("user not found")

type DB struct {
	*sql.DB
}

type User struct {
	Email           string
	IP              string
	PrivateMessages string
}

// GetUser fetches a user by ID (with context for timeouts).
func GetUser(userID string) (*User, error) {
	var user User
	err := DB.QueryRow(
		context.Background(),
		"SELECT email, ip, private_messages FROM users WHERE id = $1",
		userID,
	).Scan(&user.Email, &user.IP, &user.PrivateMessages)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}
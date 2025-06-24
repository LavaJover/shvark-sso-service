package domain

import "time"

type User struct {
	ID string
	Login string
	Username string
	Password string // hashed password
	TwoFaSecret string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUser(login, username, hashedPassword string) (*User, error) {
	return &User{
		Login: login,
		Username: username,
		Password: hashedPassword,
	}, nil
}

// Validation functions...
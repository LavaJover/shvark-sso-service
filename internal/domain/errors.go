package domain

import "errors"

var (
	ErrInvalidLogin = errors.New("invalid login format")
	ErrEmptyPassword = errors.New("password cannot be empty")
	ErrUserNotFound = errors.New("user not found")
	ErrInvalidToken = errors.New("invalid or expired token")
	ErrLoginNotFound = errors.New("login not found")
	ErrLoginAlreadyTaken = errors.New("login is already taken")
)


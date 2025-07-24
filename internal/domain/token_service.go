package domain

import "time"

type TokenService interface {
	GenerateAccessToken(user *User) (string, time.Time, error)
	GenerateRefreshToken(user *User) (string, time.Time, error)
	ValidateAccessToken(token string) (userID string, err error)
}
package usecase

import "github.com/LavaJover/shvark-sso-service/internal/domain"

type AuthUseCase interface {
	Register(login, username, password string) (userID string, err error)
	Login(login, password string) (accessToken string, err error)
	ValidateToken(token string) (userID string, err error)
	GetUserByToken(token string) (user *domain.User, err error)
}
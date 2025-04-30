package usecase

type AuthUseCase interface {
	Register(login, username, password string) (userID string, err error)
	Login(login, password string) (accessToken string, err error)
	ValidateToken(token string) (userID string, err error)
}
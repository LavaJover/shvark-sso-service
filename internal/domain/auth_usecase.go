package domain

type AuthUseCase interface {
	Register(login, username, password string) (userID string, err error)
	Login(login, password, twoFaCode string) (accessToken string, err error)
	ValidateToken(token string) (userID string, err error)
	GetUserByToken(token string) (user *User, err error)
	Setup2FA(userID string)	(string, error)
	Verify2FA(userID, code string) (bool, error)
}
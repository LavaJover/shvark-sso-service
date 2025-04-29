package usecase

type AuthUseCase interface {
	Register(login, password string) (accessToken string, err error)
	Login(login, password string) (accessToken string, err error)
}
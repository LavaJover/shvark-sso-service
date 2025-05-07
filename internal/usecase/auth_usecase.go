package usecase

import (
	"github.com/LavaJover/shvark-sso-service/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type authUseCase struct {
	repo 			domain.UserRepository
	tokenService 	domain.TokenService
}

func NewAuthUseCase(r domain.UserRepository, t domain.TokenService) domain.AuthUseCase {
	return &authUseCase{
		repo: r,
		tokenService: t,
	}
}

func (uc *authUseCase) Register(login, username, password string) (string, error) {
	// is login already in use?
	if exist, _ := uc.repo.FindByLogin(login); exist != nil{
		return "", domain.ErrInvalidLogin
	}

	// hashing password
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil{
		return "", err
	}

	// creating user
	user, err := domain.NewUser(login, username, string(hashed))
	if err != nil{
		return "", err
	}

	if err := uc.repo.Create(user); err != nil{
		return "", err
	}

	return user.ID, nil
}

func (uc *authUseCase) Login(login, password string) (string, error) {
	// searching by login
	user, err := uc.repo.FindByLogin(login)
	if err != nil{
		return "", domain.ErrInvalidLogin
	}

	// checking password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", err
	}

	// generating token
	token, err := uc.tokenService.GenerateAccessToken(user)
	if err != nil {
		return "", err
	}

	return token, nil

}

func (uc *authUseCase) ValidateToken(token string) (string, error) {
	userID, err := uc.tokenService.ValidateAccessToken(token)
	return userID, err
}

func (uc *authUseCase) GetUserByToken(token string) (*domain.User, error){
	// firstly validate the token
	userID, err := uc.tokenService.ValidateAccessToken(token)
	if err != nil{
		return nil ,err
	}

	// find concrete user by userID
	user, err := uc.repo.FindByID(userID)
	if err != nil{
		return nil, err
	}

	return user, nil
}
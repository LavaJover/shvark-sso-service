package usecase

import (
	"github.com/LavaJover/shvark-sso-service/internal/client"
	"github.com/LavaJover/shvark-sso-service/internal/domain"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type authUseCase struct {
	tokenService 	domain.TokenService
	logger 			*logrus.Entry
	userClient		*client.UserClient
}

func NewAuthUseCase(t domain.TokenService, userClient *client.UserClient) domain.AuthUseCase {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.DebugLevel)

	return &authUseCase{
		tokenService: t,
		logger: logger.WithField("component", "authUseCase"),
		userClient: userClient,
	}
}

func (uc *authUseCase) Register(login, username, password string) (string, error) {
	uc.logger.WithFields(logrus.Fields{
		"action": "Register",
		"login": "login",
	}).Info("attempt to register user")

	// is login already in use?
	if exist, _ := uc.userClient.CheckUserExists(login); exist{
		return "", domain.ErrLoginAlreadyTaken
	}

	// hashing password
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil{
		return "", err
	}

	// creating user
	user := &domain.User{
		ID: uuid.New().String(),
		Login: login,
		Username: username,
		Password: string(hashed),
	}

	// saving user on client side to database
	userID, err := uc.userClient.CreateUser(user)
	if err != nil{
		return "", err
	}

	return userID, nil
}

func (uc *authUseCase) Login(login, password string) (string, error) {
	// searching by login
	user, err := uc.userClient.GetUserByLogin(login)
	if err != nil{
		return "", nil
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
	return uc.userClient.GetUserByID(userID)
}
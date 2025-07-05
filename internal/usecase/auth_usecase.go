package usecase

import (
	"fmt"

	"github.com/LavaJover/shvark-sso-service/internal/client"
	"github.com/LavaJover/shvark-sso-service/internal/domain"
	"github.com/LavaJover/shvark-sso-service/internal/infrastructure/google2fa"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (uc *authUseCase) Login(login, password, twoFaCode string) (string, error) {
	// searching by login
	user, err := uc.userClient.GetUserByLogin(login)
	if err != nil{
		return "", err
	}

	fmt.Println(user.TwoFaEnabled, twoFaCode)

	// checking password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", err
	}

	// checking if user using google 2 FA
	if user.TwoFaEnabled && twoFaCode == "" {
		return "", status.Error(codes.Unauthenticated, "2FA_REQUIRED")
	}

	if user.TwoFaEnabled && !google2fa.Verify2FACode(user.TwoFaSecret, twoFaCode) {
		return "", status.Error(codes.Unauthenticated, "WRONG_CREDENTIALS")
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

func (uc *authUseCase) Setup2FA(userID string) (string, error) {
	user, err := uc.userClient.GetUserByID(userID)
	if err != nil {
		return "", err
	}
	// Проверка: если уже создан secret, то не пересоздавать!
	if user.TwoFaSecret != "" {
		return fmt.Sprintf("otpauth://totp/%s:%s?secret=%s&issuer=%s", "ShavrkPay", user.Login, user.TwoFaSecret, "ShvarkPay"), nil
	}
	login := user.Login
	secret, qrURL, err := google2fa.Generate2FASecret(login)
	if err != nil {
		return "", err
	}

	err = uc.userClient.SetTwoFaSecret(userID, secret)
	if err != nil {
		return "", err
	}

	err = uc.userClient.SetTwoFaEnabled(userID, false)
	if err != nil {
		return "", err
	}
	
	return qrURL, nil
}

func (uc *authUseCase) Verify2FA(userID, code string) (bool, error) {
	user, err := uc.userClient.GetUserByID(userID)
	if err != nil {
		return false, err
	}
	if google2fa.Verify2FACode(user.TwoFaSecret, code) {
		uc.userClient.SetTwoFaEnabled(userID, true)
		return true, nil
	}
	return false, nil
}
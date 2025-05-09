package usecase

import (
	"testing"
	"time"

	"github.com/LavaJover/shvark-sso-service/internal/domain"
	"github.com/LavaJover/shvark-sso-service/internal/infrastructure/jwt"
	"github.com/LavaJover/shvark-sso-service/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthUsecase_Register_Success(t *testing.T){
	mockUserRepo := new(mocks.UserRepository)
	mockToken := new(mocks.TokenService)

	login := "boriska"
	username := "britva1"
	password := "secure"

	// Поведение: FindByLogin ничего не нашёл (логин нового пользователя не занят)
	mockUserRepo.EXPECT().
		FindByLogin(login).
		Return(nil, nil)

	// Поведение: Create возвращает nil, успех
	mockUserRepo.On("Create", mock.AnythingOfType("*domain.User")).
		Run(func(args mock.Arguments) {
			user := args.Get(0).(*domain.User)
			user.ID = "some-user-ID"
		}).
		Return(nil)

	authUC := NewAuthUseCase(mockUserRepo, mockToken)

	userID, err := authUC.Register(login, username, password)

	assert.NoError(t, err)
	assert.NotEmpty(t, userID)

	mockUserRepo.AssertExpectations(t)
}

func TestAuthUsecase_Register_LoginExists(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)
	mockTocken := new(mocks.TokenService)

	login := "boris"
	username := "britva"
	password := "secure"

	mockUserRepo.EXPECT().
		FindByLogin(login).
		Return(&domain.User{Login: login}, nil)

	authUC := NewAuthUseCase(mockUserRepo, mockTocken)

	userID, err := authUC.Register(login, username, password)

	assert.ErrorIs(t, err, domain.ErrInvalidLogin)
	assert.Empty(t, userID)

	mockUserRepo.AssertExpectations(t)
}

func TestAuthUsecase_Login_Success(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)
	mockToken := new(mocks.TokenService)

	id := "some-id"
	login := "boris"
	username := "britva"
	password := "secure"

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	assert.NoError(t, err)

	user := &domain.User{ID: id, Login: login, Username: username, Password: string(hashedPassword)}

	// Поведение: логин пользователя найден
	mockUserRepo.On("FindByLogin", login).
		Return(user, nil)

	// Поведение: access-token пользователя успешно сгенерирован
	expectedToken := "some-token"
	mockToken.On("GenerateAccessToken", user).
		Return(expectedToken, nil)

	// Запуск auth usecase
	authUC := NewAuthUseCase(mockUserRepo, mockToken)
	actualToken, err := authUC.Login(login, password)
	
	assert.NoError(t, err)
	assert.Equal(t, expectedToken, actualToken)

	mockUserRepo.AssertExpectations(t)
	mockToken.AssertExpectations(t)
}

func TestAuthUsecase_Login_UserNotFound(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)
	mockToken := new(mocks.TokenService)

	login := "boris"
	password := "secure"

	// Поведение: логин не найден в системе
	mockUserRepo.On("FindByLogin", login).
		Return(nil, domain.ErrLoginNotFound)

	authUC := NewAuthUseCase(mockUserRepo, mockToken)
	accessToken, err := authUC.Login(login, password)

	assert.Empty(t, accessToken)
	assert.ErrorIs(t, err, domain.ErrLoginNotFound)

	mockUserRepo.AssertExpectations(t)
}

func TestAuthUsecase_GetUserByToken_Success(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)
	mockToken := new(mocks.TokenService)

	userID := "some-user-id" 
	user := &domain.User{ID: userID}
	accessToken, err := jwt.NewTokenService("secret", time.Minute*15).GenerateAccessToken(user)
	require.NoError(t, err)
	require.NotEmpty(t, accessToken)

	// Поведение: токен валиден
	mockToken.On("ValidateAccessToken", accessToken).
		Return(userID, nil)
	
	// Поведение: пользователь с таким ID найден
	mockUserRepo.On("FindByID", userID).
		Return(&domain.User{ID: userID}, nil)

	// Запуск usecase
	authUC := NewAuthUseCase(mockUserRepo, mockToken)
	actualUser, err := authUC.GetUserByToken(accessToken)
	require.NoError(t, err)
	assert.Equal(t, userID, actualUser.ID)

	mockUserRepo.AssertExpectations(t)
	mockToken.AssertExpectations(t)
}
package usecase

import (
	"testing"

	"github.com/LavaJover/shvark-sso-service/internal/domain"
	"github.com/LavaJover/shvark-sso-service/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
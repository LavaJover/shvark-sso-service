package usecase

import (
	"testing"

	"github.com/LavaJover/shvark-sso-service/internal/domain"
	"github.com/LavaJover/shvark-sso-service/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
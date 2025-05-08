package jwt

import (
	"testing"
	"time"

	"github.com/LavaJover/shvark-sso-service/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestValidateAccessToken(t *testing.T) {
	users := []*domain.User{
		&domain.User{ID: "some-id"},
		&domain.User{ID: "another-id"},
		&domain.User{ID: "last-id"},
	}

	tokenService := NewTokenService("secret", time.Minute*15)
	accessTokens := make([]string, 3)
	for i := range 3{
		accessToken, err := tokenService.GenerateAccessToken(users[i])
		assert.NoError(t, err)
		accessTokens[i] = accessToken
	}

	for i := range 3{
		userID, err := tokenService.ValidateAccessToken(accessTokens[i])
		assert.NoError(t, err)
		assert.Equal(t, users[i].ID, userID)
	}
}
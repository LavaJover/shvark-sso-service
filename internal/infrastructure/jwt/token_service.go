package jwt

import (
	"time"
	"errors"

	"github.com/LavaJover/shvark-sso-service/internal/domain"
	"github.com/golang-jwt/jwt/v5"
)

type jwtTokenService struct {
	secretKey string
	ttl time.Duration
}

func NewTokenService(secretKey string, ttl time.Duration) domain.TokenService {
	return &jwtTokenService{
		secretKey: secretKey,
		ttl: ttl,
	}
}

func (s *jwtTokenService) GenerateAccessToken(user *domain.User) (string, time.Time, error){
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(s.ttl).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(s.secretKey))
	return tokenStr, time.Now().Add(s.ttl), err
}

func (s *jwtTokenService) GenerateRefreshToken(user *domain.User) (string, time.Time, error) {
	// improve later...
	return s.GenerateAccessToken(user)
}

func (s *jwtTokenService) ValidateAccessToken(tokenStr string) (string, error){
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.secretKey), nil
	})

	if err != nil || !token.Valid {
		return "", domain.ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", domain.ErrInvalidToken
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", domain.ErrInvalidToken
	}

	return userID, nil
}
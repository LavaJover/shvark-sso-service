package grpcapi

import (
	"context"

	"github.com/LavaJover/shvark-sso-service/internal/usecase"
	ssopb "github.com/LavaJover/shvark-sso-service/proto/gen"
)

type AuthHandler struct {
	ssopb.UnimplementedSSOServiceServer
	usecase.AuthUseCase
}

func (h *AuthHandler) Register(ctx context.Context, req *ssopb.RegisterRequest) (*ssopb.RegisterResponse, error) {
	userID, err := h.AuthUseCase.Register(req.Login, req.Username, req.Password)
	if err != nil{
		return nil, err
	}

	return &ssopb.RegisterResponse{
		UserId: userID,
		Message: "Successfully created user",
	}, nil
}

func (h *AuthHandler) Login(ctx context.Context, req *ssopb.LoginRequest) (*ssopb.LoginResponse, error) {
	accessToken, err := h.AuthUseCase.Login(req.Login, req.Password)
	if err != nil{
		return nil, err
	}

	return &ssopb.LoginResponse{
		AccessToken: accessToken,
		RefreshToken: accessToken,
	}, nil
}

func (h *AuthHandler) ValidateToken(ctx context.Context, req *ssopb.ValidateTokenRequest) (*ssopb.ValidateTokenResponse, error) {
	userID, err := h.AuthUseCase.ValidateToken(req.AccessToken)
	if err != nil{
		return nil, err
	}

	return &ssopb.ValidateTokenResponse{
		Valid: true,
		UserId: userID,
	}, nil
}

func (h *AuthHandler) GetUserByToken(ctx context.Context, req *ssopb.GetUserByTokenRequest) (*ssopb.GetUserByTokenResponse, error) {
	user, err := h.AuthUseCase.GetUserByToken(req.AccessToken)
	if err != nil{
		return nil, err
	}

	return &ssopb.GetUserByTokenResponse{
		UserId: user.ID,
		Login: user.Login,
		Username: user.Username,
	}, nil
}
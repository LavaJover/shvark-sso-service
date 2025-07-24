package grpcapi

import (
	"context"

	"github.com/LavaJover/shvark-sso-service/internal/domain"
	ssopb "github.com/LavaJover/shvark-sso-service/proto/gen"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AuthHandler struct {
	ssopb.UnimplementedSSOServiceServer
	domain.AuthUseCase
}

func (h *AuthHandler) Register(ctx context.Context, req *ssopb.RegisterRequest) (*ssopb.RegisterResponse, error) {
	userID, err := h.AuthUseCase.Register(req.Login, req.Username, req.Password, req.Role)
	if err != nil{
		return nil, err
	}

	return &ssopb.RegisterResponse{
		UserId: userID,
		Message: "Successfully created user",
	}, nil
}

func (h *AuthHandler) Login(ctx context.Context, req *ssopb.LoginRequest) (*ssopb.LoginResponse, error) {
	accessToken, timeExp, err := h.AuthUseCase.Login(req.Login, req.Password, req.TwoFaCode)
	if err != nil{
		return nil, err
	}

	return &ssopb.LoginResponse{
		AccessToken: accessToken,
		RefreshToken: accessToken,
		TimeExp: timestamppb.New(timeExp),
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

func (h *AuthHandler) Setup2FA(ctx context.Context, r *ssopb.Setup2FARequest) (*ssopb.Setup2FAResponse, error) {
	userID := r.UserId
	qrURL, err := h.AuthUseCase.Setup2FA(userID)
	if err != nil {
		return nil, err
	}

	return &ssopb.Setup2FAResponse{
		QrUrl: qrURL,
	}, nil
}

func (h *AuthHandler) Verify2FA(ctx context.Context, r *ssopb.Verify2FARequest) (*ssopb.Verify2FAResponse, error) {
	userID, code := r.UserId, r.Code
	verif, err := h.AuthUseCase.Verify2FA(userID, code)
	if err != nil {
		return nil, err
	}

	return &ssopb.Verify2FAResponse{
		Verif: verif,
	}, nil
}
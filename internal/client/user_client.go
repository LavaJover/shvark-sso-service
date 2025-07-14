package client

import (
	"context"
	"time"

	"github.com/LavaJover/shvark-sso-service/internal/domain"
	userpb "github.com/LavaJover/shvark-user-service/proto/gen"

	"google.golang.org/grpc"
)

type UserClient struct {
	client userpb.UserServiceClient
}

func NewUserClient(conn *grpc.ClientConn) *UserClient {
	return &UserClient{
		client: userpb.NewUserServiceClient(conn),
	}
}

func (userClient *UserClient) CheckUserExists(login string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	_, err := userClient.client.GetUserByLogin(ctx, &userpb.GetUserByLoginRequest{Login: login})
	if err != nil{
		return false, err
	}

	return true, nil
}

func (userClient *UserClient) CreateUser(user *domain.User) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	response, err := userClient.client.CreateUser(
		ctx, 
		&userpb.CreateUserRequest{
			Login: user.Login,
			Username: user.Username,
			Password: user.Password,
			Role: user.Role,
		},
	)

	if err != nil{
		return "", err
	}

	return response.UserId, nil
}

func (userClient *UserClient) GetUserByLogin(login string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	response, err := userClient.client.GetUserByLogin(
		ctx,
		&userpb.GetUserByLoginRequest{
			Login: login,
		},
	)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		ID: response.UserId,
		Login: response.Login,
		Username: response.Username,
		Password: response.Password,
		TwoFaSecret: response.TwoFaSecret,
		TwoFaEnabled: response.TwoFaEnabled,
	}

	return user, nil
}

func (userClient *UserClient) GetUserByID(userID string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	response, err := userClient.client.GetUserByID(
		ctx,
		&userpb.GetUserByIDRequest{
			UserId: userID,
		},
	)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		ID: response.UserId,
		Username: response.Username,
		Login: response.Login,
		Password: response.Password,
		TwoFaSecret: response.TwoFaSecret,
		TwoFaEnabled: response.TwoFaEnabled,
	}

	return user, nil
}

func (c *UserClient) SetTwoFaSecret(userID, twoFaSecret string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := c.client.SetTwoFaSecret(
		ctx,
		&userpb.SetTwoFaSecretRequest{
			UserId: userID,
			TwoFaSecret: twoFaSecret,
		},
	)

	return err
}

func (c *UserClient) SetTwoFaEnabled(userID string, enabled bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := c.client.SetTwoFaEnabled(
		ctx,
		&userpb.SetTwoFaEnabledRequest{
			UserId: userID,
			Enabled: enabled,
		},
	)

	return err
}
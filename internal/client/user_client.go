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
		},
	)

	if err != nil{
		return "", err
	}

	return response.UserId, nil
}
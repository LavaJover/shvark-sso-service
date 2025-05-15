package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/LavaJover/shvark-sso-service/internal/client"
	"github.com/LavaJover/shvark-sso-service/internal/config"
	"github.com/LavaJover/shvark-sso-service/internal/delivery/grpcapi"
	"github.com/LavaJover/shvark-sso-service/internal/infrastructure/jwt"
	"github.com/LavaJover/shvark-sso-service/internal/usecase"
	ssopb "github.com/LavaJover/shvark-sso-service/proto/gen"
	"google.golang.org/grpc"
)

func main(){
	// processing app config
	cfg := config.MustLoad()

	// init user-service client
	conn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect to user-service: %v\n", err)
	}
	defer conn.Close()

	userClient := client.NewUserClient(conn)

	// creating token service
	tokenService := jwt.NewTokenService("my-secret-word", 15*time.Minute)

	// creating Auth usecase
	authUseCase := usecase.NewAuthUseCase(tokenService, userClient)

	// creating gRPC server
	grpcServer := grpc.NewServer()
	authHandler := grpcapi.AuthHandler{AuthUseCase: authUseCase}

	ssopb.RegisterSSOServiceServer(grpcServer, &authHandler)

	// start
	lis, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil{
		log.Fatalf("failed to listen: %v", err)
	}

	fmt.Printf("gRPC server started on %s:%s\n", cfg.Host, cfg.Port)
	if err := grpcServer.Serve(lis); err != nil{
		log.Fatalf("failed to serve: %v\n", err)
	}

}
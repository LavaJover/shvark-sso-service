package main

import (
	"fmt"
	"net"
	"time"
	"log"

	"github.com/LavaJover/shvark-sso-service/internal/config"
	"github.com/LavaJover/shvark-sso-service/internal/delivery/grpcapi"
	"github.com/LavaJover/shvark-sso-service/internal/infrastructure/jwt"
	"github.com/LavaJover/shvark-sso-service/internal/infrastructure/postgres"
	"github.com/LavaJover/shvark-sso-service/internal/logger"
	"github.com/LavaJover/shvark-sso-service/internal/usecase"
	ssopb "github.com/LavaJover/shvark-sso-service/proto/gen"
	"google.golang.org/grpc"
)

func main(){
	// processing app config
	cfg := config.MustLoad()

	// creating logger
	myLog := logger.InitLogger(cfg)
	
	// init db
	myLog.Debug("init database")

	dsn := cfg.Dsn
	db := postgres.InitDB(dsn)

	// creating user repo
	userRepo := postgres.NewUserRepository(db)

	// creating token service
	tokenService := jwt.NewTokenService("my-secret-word", 15*time.Minute)

	// creating Auth usecase
	authUseCase := usecase.NewAuthUseCase(userRepo, tokenService)

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
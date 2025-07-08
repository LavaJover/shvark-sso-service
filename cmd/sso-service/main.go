package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/LavaJover/shvark-sso-service/internal/client"
	"github.com/LavaJover/shvark-sso-service/internal/config"
	"github.com/LavaJover/shvark-sso-service/internal/delivery/grpcapi"
	"github.com/LavaJover/shvark-sso-service/internal/infrastructure/jwt"
	"github.com/LavaJover/shvark-sso-service/internal/interceptor"
	"github.com/LavaJover/shvark-sso-service/internal/usecase"
	ssopb "github.com/LavaJover/shvark-sso-service/proto/gen"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/joho/godotenv"
)

func main(){
	if err := godotenv.Load(); err != nil {
		log.Println("failed to load .env")
	}

	// processing app config
	cfg := config.MustLoad()

	fmt.Println(cfg.UserService.Host)

	// Init logger
	logger, _ := zap.NewProduction()

	// Init retry system
	// Retry all calls with retryable errors (UNAVAILABLE, DEADLINE_EXCEEDED)
	opts := []grpc_retry.CallOption{
		grpc_retry.WithMax(uint(cfg.MaxAttempts)),
		grpc_retry.WithBackoff(grpc_retry.BackoffExponential(cfg.InitialBackoff)),
		grpc_retry.WithCodes(grpc_retry.DefaultRetriableCodes...),
	}

	// init user-service client
	conn, err := grpc.Dial(
		fmt.Sprintf("%s:%s", cfg.UserService.Host, cfg.UserService.Port),
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(opts...)))
	if err != nil {
		log.Fatalf("failed to connect to user-service: %v\n", err)
	}
	defer conn.Close()

	userClient := client.NewUserClient(conn)

	// creating token service
	tokenService := jwt.NewTokenService(os.Getenv("JWT_SECRET"), 12*time.Hour)

	// creating Auth usecase
	authUseCase := usecase.NewAuthUseCase(tokenService, userClient)

	// creating gRPC server
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			interceptor.NewLoggerInterceptor(logger),
		)),
	)
	authHandler := grpcapi.AuthHandler{AuthUseCase: authUseCase}

	ssopb.RegisterSSOServiceServer(grpcServer, &authHandler)

	// start
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.GRPCServer.Host, cfg.GRPCServer.Port))
	if err != nil{
		log.Fatalf("failed to listen: %v", err)
	}

	fmt.Printf("gRPC server started on %s:%s\n", cfg.GRPCServer.Host, cfg.GRPCServer.Port)
	if err := grpcServer.Serve(lis); err != nil{
		log.Fatalf("failed to serve: %v\n", err)
	}

}
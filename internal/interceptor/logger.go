package interceptor

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func NewLoggerInterceptor(logger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		// before calling business logic
		start := time.Now()
		logger.Info("Incoming gRPC request",
			zap.String("method", info.FullMethod),
			zap.Any("request", req),
		)
		// calling business logic
		resp, err = handler(ctx, req)
		// afetr calling business logic
		logger.Info("Completed gRPC request",
			zap.String("method", info.FullMethod),
			zap.Duration("duration", time.Since(start)),
			zap.Error(err),
		)

		return
	}
}
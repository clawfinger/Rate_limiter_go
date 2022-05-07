package grpcserver

import (
	"context"
	"time"

	"github.com/clawfinger/ratelimiter/internal/logger"
	"google.golang.org/grpc"
)

func LoggerInterceptor(logger logger.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		h, err := handler(ctx, req)
		logger.Info("Called", info.FullMethod, "with execution time", time.Since(start))
		return h, err
	}
}

package ratelimiter

import (
	"context"

	"github.com/clawfinger/ratelimiter/config"
	"github.com/clawfinger/ratelimiter/internal/logger"
	grpcserver "github.com/clawfinger/ratelimiter/internal/server/grpc"
)

type LimiterService struct {
	cfg    *config.Config
	logger logger.Logger
	server *grpcserver.GrpcServer
}

func New(cfg *config.Config, logger logger.Logger, server *grpcserver.GrpcServer) *LimiterService {
	return &LimiterService{
		cfg:    cfg,
		logger: logger,
		server: server,
	}
}

func (s *LimiterService) Start(ctx context.Context) {
	go func() {
		err := s.server.Start()
		if err != nil {
			s.logger.Info("Failed to start grpc server")
		}
	}()
	<-ctx.Done()
}

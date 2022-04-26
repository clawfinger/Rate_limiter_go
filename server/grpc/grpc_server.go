package grpcserver

import (
	"context"
	"log"
	"net"

	pb "github.com/clawfinger/ratelimiter/api/generated"
	storage "github.com/clawfinger/ratelimiter/redis"
	servers "github.com/clawfinger/ratelimiter/server"
	"google.golang.org/grpc"
)

type GrpcServer struct {
	context *servers.ServerCommonContext
	server  *grpc.Server
	pb.UnimplementedLimiterServer
}

func NewGrpcServer(context *servers.ServerCommonContext) *GrpcServer {
	return &GrpcServer{
		context: context,
	}
}

func (s *GrpcServer) Start() error {
	lsn, err := net.Listen("tcp", s.context.Cfg.Data.Grpc.Addr)
	if err != nil {
		log.Fatal(err)
	}

	s.server = grpc.NewServer(grpc.ChainUnaryInterceptor(LoggerInterceptor(s.context.Logger)))
	pb.RegisterLimiterServer(s.server, s)
	return s.server.Serve(lsn)
}

func (s *GrpcServer) Stop() error {
	s.server.Stop()
	return nil
}

func (s *GrpcServer) Validate(ctx context.Context, attempt *pb.LoginAttempt) (*pb.AttemptResult, error) {
	answer := &pb.AttemptResult{}

	ipResult := s.context.Storage.CheckIP(ctx, attempt.IP)

	if ipResult.Status == storage.Whitelisted {
		answer.Result = pb.AttemptResult_OK
		return answer, nil
	}
	if ipResult.Status == storage.Blacklisted {
		answer.Result = pb.AttemptResult_DENIED
		return answer, nil
	}

	result := s.context.RateManager.Manage(attempt.IP, attempt.Login, attempt.Password)

	s.context.Logger.Info("Validation result for request", "ip", attempt.IP, "login", attempt.Login,
		"password", attempt.Password, "status", result.Ok, "reason", result.Reason)
	if result.Ok {
		answer.Result = pb.AttemptResult_OK
	} else {
		answer.Result = pb.AttemptResult_DENIED
	}
	return answer, nil
}

func (s *GrpcServer) DropStats(ctx context.Context, stats *pb.Stats) (*pb.OperationResult, error) {
	res := &pb.OperationResult{}

	s.context.RateManager.DropStats(stats.Login, stats.IP)

	res.Status = pb.OperationResult_OK
	res.Reason = "DropStats ok"
	return res, nil
}

func (s *GrpcServer) AddBlacklist(ctx context.Context, ip *pb.Subnet) (*pb.OperationResult, error) {
	res := &pb.OperationResult{}
	res.Status = pb.OperationResult_OK
	res.Reason = "AddBlacklist ok"

	s.context.Storage.SetIP(ctx, ip.IPWithMask, storage.Blacklisted)

	return res, nil
}

func (s *GrpcServer) RemoveBlacklist(ctx context.Context, subnet *pb.Subnet) (*pb.OperationResult, error) {
	res := &pb.OperationResult{}
	res.Status = pb.OperationResult_OK
	s.context.Storage.RemoveIP(ctx, subnet.IPWithMask, storage.Blacklisted)
	return res, nil
}

func (s *GrpcServer) AddWhitelist(ctx context.Context, subnet *pb.Subnet) (*pb.OperationResult, error) {
	res := &pb.OperationResult{}
	res.Status = pb.OperationResult_OK
	res.Reason = "AddWhitelist ok"
	s.context.Storage.SetIP(ctx, subnet.IPWithMask, storage.Blacklisted)

	return res, nil
}

func (s *GrpcServer) RemoveWhitelist(ctx context.Context, subnet *pb.Subnet) (*pb.OperationResult, error) {
	res := &pb.OperationResult{}
	res.Status = pb.OperationResult_OK
	res.Reason = "RemoveWhitelist ok"
	s.context.Storage.RemoveIP(ctx, subnet.IPWithMask, storage.Blacklisted)
	return res, nil
}

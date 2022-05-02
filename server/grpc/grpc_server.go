package grpcserver

import (
	"context"
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
		s.context.Logger.Error("Failed to start grpc server", err.Error())
	}

	s.server = grpc.NewServer(grpc.ChainUnaryInterceptor(LoggerInterceptor(s.context.Logger)))
	pb.RegisterLimiterServer(s.server, s)
	return s.server.Serve(lsn)
}

func (s *GrpcServer) Stop() {
	s.server.Stop()
}

//nolint
func (s *GrpcServer) Validate(ctx context.Context, attempt *pb.LoginAttempt) (*pb.AttemptResult, error) {
	answer := &pb.AttemptResult{}

	ipResult := s.context.Storage.CheckIP(ctx, attempt.IP)

	if ipResult.Err != nil {
		answer.Result = pb.AttemptResult_DENIED
		return answer, nil
	}
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

func (s *GrpcServer) DropIPStats(ctx context.Context, stats *pb.Stats) (*pb.OperationResult, error) {
	res := &pb.OperationResult{}

	s.context.RateManager.DropIPStats(stats.Data)

	res.Status = pb.OperationResult_OK
	res.Reason = "Drop ip stats ok"
	return res, nil
}

func (s *GrpcServer) DropLoginStats(ctx context.Context, stats *pb.Stats) (*pb.OperationResult, error) {
	res := &pb.OperationResult{}

	s.context.RateManager.DropLiginStats(stats.Data)

	res.Status = pb.OperationResult_OK
	res.Reason = "Drop login stats ok"
	return res, nil
}

func (s *GrpcServer) DropPasswordStats(ctx context.Context, stats *pb.Stats) (*pb.OperationResult, error) {
	res := &pb.OperationResult{}

	s.context.RateManager.DropPasswordStats(stats.Data)

	res.Status = pb.OperationResult_OK
	res.Reason = "Drop password stats ok"
	return res, nil
}

func (s *GrpcServer) AddBlacklist(ctx context.Context, ip *pb.Subnet) (*pb.OperationResult, error) {
	res := &pb.OperationResult{}
	res.Status = pb.OperationResult_OK
	res.Reason = "AddBlacklist ok"

	err := s.context.Storage.SetIP(ctx, ip.IPWithMask, storage.Blacklisted)
	if err != nil {
		res.Status = pb.OperationResult_FAIL
		res.Reason = err.Error()
	}

	return res, nil
}

func (s *GrpcServer) RemoveBlacklist(ctx context.Context, subnet *pb.Subnet) (*pb.OperationResult, error) {
	res := &pb.OperationResult{}
	res.Status = pb.OperationResult_OK
	err := s.context.Storage.RemoveIP(ctx, subnet.IPWithMask, storage.Blacklisted)
	if err != nil {
		res.Status = pb.OperationResult_FAIL
		res.Reason = err.Error()
	}
	return res, nil
}

func (s *GrpcServer) AddWhitelist(ctx context.Context, subnet *pb.Subnet) (*pb.OperationResult, error) {
	res := &pb.OperationResult{}
	res.Status = pb.OperationResult_OK
	res.Reason = "AddWhitelist ok"
	err := s.context.Storage.SetIP(ctx, subnet.IPWithMask, storage.Blacklisted)
	if err != nil {
		res.Status = pb.OperationResult_FAIL
		res.Reason = err.Error()
	}
	return res, nil
}

func (s *GrpcServer) RemoveWhitelist(ctx context.Context, subnet *pb.Subnet) (*pb.OperationResult, error) {
	res := &pb.OperationResult{}
	res.Status = pb.OperationResult_OK
	res.Reason = "RemoveWhitelist ok"
	err := s.context.Storage.RemoveIP(ctx, subnet.IPWithMask, storage.Blacklisted)
	if err != nil {
		res.Status = pb.OperationResult_FAIL
		res.Reason = err.Error()
	}
	return res, nil
}

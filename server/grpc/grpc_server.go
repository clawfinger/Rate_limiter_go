package grpcserver

import (
	"context"
	"log"
	"net"

	pb "github.com/clawfinger/ratelimiter/api/generated"
	servers "github.com/clawfinger/ratelimiter/server"
	"google.golang.org/grpc"
)

type GrpcServer struct {
	context *servers.ServerContext
	server  *grpc.Server
	pb.UnimplementedLimiterServer
}

func NewGrpcServer(context *servers.ServerContext) *GrpcServer {
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

func (s *GrpcServer) Validate(context.Context, *pb.LoginAttempt) (*pb.AttemptResult, error) {
	res := &pb.AttemptResult{}
	res.Result = pb.AttemptResult_DENIED
	return res, nil
}

func (s *GrpcServer) DropStats(context.Context, *pb.Stats) (*pb.OperationResult, error) {
	res := &pb.OperationResult{}
	res.Status = pb.OperationResult_OK
	res.Reason = "DropStats ok"
	return res, nil
}

func (s *GrpcServer) AddBlacklist(context.Context, *pb.IP) (*pb.OperationResult, error) {
	res := &pb.OperationResult{}
	res.Status = pb.OperationResult_OK
	res.Reason = "AddBlacklist ok"
	return res, nil
}

func (s *GrpcServer) RemoveBlacklist(context.Context, *pb.IP) (*pb.OperationResult, error) {
	res := &pb.OperationResult{}
	res.Status = pb.OperationResult_OK
	res.Reason = "RemoveBlacklist ok"
	return res, nil
}

func (s *GrpcServer) AddWhitelist(context.Context, *pb.IP) (*pb.OperationResult, error) {
	res := &pb.OperationResult{}
	res.Status = pb.OperationResult_OK
	res.Reason = "AddWhitelist ok"
	return res, nil
}

func (s *GrpcServer) RemoveWhitelist(context.Context, *pb.IP) (*pb.OperationResult, error) {
	res := &pb.OperationResult{}
	res.Status = pb.OperationResult_OK
	res.Reason = "RemoveWhitelist ok"
	return res, nil
}

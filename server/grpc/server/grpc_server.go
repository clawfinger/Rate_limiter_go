package grpcserver

import (
	"log"
	"net"

	servers "github.com/clawfinger/ratelimiter/server"
	"google.golang.org/grpc"
)

type GrpcServer struct {
	context *servers.ServerContext
	server  *grpc.Server
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

	return s.server.Serve(lsn)
}

func (s *GrpcServer) Stop() error {
	s.server.Stop()
	return nil
}

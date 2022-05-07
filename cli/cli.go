package cli

import (
	"context"
	"time"

	pb "github.com/clawfinger/ratelimiter/api/generated"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type CtlAgent struct {
	client pb.LimiterClient
	conn   *grpc.ClientConn
}

func NewClient(addr string) (*CtlAgent, error) {
	agent := &CtlAgent{}
	err := agent.Connect(addr)
	return agent, err
}

func (c *CtlAgent) Connect(addr string) error {
	var err error
	c.conn, err = grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	c.client = pb.NewLimiterClient(c.conn)
	return err
}

func (c *CtlAgent) DropIPStats(ctx context.Context, in *pb.Stats) (*pb.OperationResult, error) {
	timedCtx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	return c.client.DropIPStats(timedCtx, in)
}

func (c *CtlAgent) DropLoginStats(ctx context.Context, in *pb.Stats) (*pb.OperationResult, error) {
	timedCtx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	return c.client.DropLoginStats(timedCtx, in)
}

func (c *CtlAgent) DropPasswordStats(ctx context.Context, in *pb.Stats) (*pb.OperationResult, error) {
	timedCtx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	return c.client.DropPasswordStats(timedCtx, in)
}

func (c *CtlAgent) AddBlacklist(ctx context.Context, in *pb.Subnet) (*pb.OperationResult, error) {
	timedCtx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	return c.client.AddBlacklist(timedCtx, in)
}

func (c *CtlAgent) RemoveBlacklist(ctx context.Context, in *pb.Subnet) (*pb.OperationResult, error) {
	timedCtx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	return c.client.RemoveBlacklist(timedCtx, in)
}

func (c *CtlAgent) AddWhitelist(ctx context.Context, in *pb.Subnet) (*pb.OperationResult, error) {
	timedCtx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	return c.client.AddWhitelist(timedCtx, in)
}

func (c *CtlAgent) RemoveWhitelist(ctx context.Context, in *pb.Subnet) (*pb.OperationResult, error) {
	timedCtx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	return c.client.RemoveWhitelist(timedCtx, in)
}

func (c *CtlAgent) Close() {
	c.conn.Close()
}

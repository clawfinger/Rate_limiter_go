package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"syscall"

	"github.com/clawfinger/ratelimiter/config"
	"github.com/clawfinger/ratelimiter/logger"
	ratelimiter "github.com/clawfinger/ratelimiter/root"
	serverctx "github.com/clawfinger/ratelimiter/server"
	grpcserver "github.com/clawfinger/ratelimiter/server/grpc"

	"github.com/spf13/cobra"
)

var configFile string

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()
	rootCmd := &cobra.Command{
		Use: "limiter",
		Run: func(cmd *cobra.Command, args []string) {
			config := config.NewConfig()
			err := config.Init(configFile)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error on config init, Reason: %s", err.Error())
				return
			}
			logger, err := logger.New(config.Data.Logger.Level, config.Data.Logger.Filename)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error on logger init, Reason: %s", err.Error())
				return
			}

			serverContext := serverctx.NewServerContext(config, logger)
			grpcServer := grpcserver.NewGrpcServer(serverContext)
			defer grpcServer.Stop()
			app := ratelimiter.New(config, logger, grpcServer)

			app.Start(ctx)
		},
	}
	executablePath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	defaultCfgPath := path.Join(filepath.Dir(executablePath), "config.json")
	flags := rootCmd.Flags()
	flags.StringVarP(&configFile, "config", "c", defaultCfgPath, "Config file path")
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Command line error")
	}
}

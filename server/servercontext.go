package servers

import (
	"github.com/clawfinger/ratelimiter/config"
	"github.com/clawfinger/ratelimiter/logger"
)

type ServerContext struct {
	Cfg    *config.Config
	Logger logger.Logger
}

func NewServerContext(cfg *config.Config, logger logger.Logger) *ServerContext {
	return &ServerContext{
		Cfg:    cfg,
		Logger: logger,
	}
}

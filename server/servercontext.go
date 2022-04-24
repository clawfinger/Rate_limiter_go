package serverctx

import (
	"github.com/clawfinger/ratelimiter/config"
	"github.com/clawfinger/ratelimiter/logger"
	manager "github.com/clawfinger/ratelimiter/ratemanager"
)

type ServerContext struct {
	Cfg         *config.Config
	Logger      logger.Logger
	RateManager manager.AbstractRateManager
}

func NewServerContext(cfg *config.Config, logger logger.Logger, RateManager manager.AbstractRateManager) *ServerContext {
	return &ServerContext{
		Cfg:         cfg,
		Logger:      logger,
		RateManager: RateManager,
	}
}

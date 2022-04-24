package serverctx

import (
	"github.com/clawfinger/ratelimiter/config"
	"github.com/clawfinger/ratelimiter/logger"
	manager "github.com/clawfinger/ratelimiter/ratemanager"
	storage "github.com/clawfinger/ratelimiter/redis"
)

type ServerCommonContext struct {
	Cfg         *config.Config
	Logger      logger.Logger
	RateManager manager.AbstractRateManager
	Storage     storage.AbstractStorage
}

func NewServerContext(cfg *config.Config, logger logger.Logger,
	RateManager manager.AbstractRateManager, Storage storage.AbstractStorage) *ServerCommonContext {
	return &ServerCommonContext{
		Cfg:         cfg,
		Logger:      logger,
		RateManager: RateManager,
		Storage:     Storage,
	}
}

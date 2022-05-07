package serverctx

import (
	"github.com/clawfinger/ratelimiter/config"
	internalapi "github.com/clawfinger/ratelimiter/internalapi"
	"github.com/clawfinger/ratelimiter/logger"
)

type ServerCommonContext struct {
	Cfg         *config.Config
	Logger      logger.Logger
	RateManager internalapi.AbstractRateManager
	Storage     internalapi.AbstractStorage
}

func NewServerContext(cfg *config.Config, logger logger.Logger,
	ratemanager internalapi.AbstractRateManager, storage internalapi.AbstractStorage) *ServerCommonContext {
	return &ServerCommonContext{
		Cfg:         cfg,
		Logger:      logger,
		RateManager: ratemanager,
		Storage:     storage,
	}
}

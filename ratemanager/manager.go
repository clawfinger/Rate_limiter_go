package manager

import (
	"time"

	"github.com/clawfinger/ratelimiter/config"
	"github.com/clawfinger/ratelimiter/logger"
	"github.com/juju/ratelimit"
)

type BucketData struct {
	Bucket     *ratelimit.Bucket
	LastActive time.Time
}

type RateManager struct {
	ipBuckets       map[string]*ratelimit.Bucket
	loginBuckets    map[string]*ratelimit.Bucket
	passwordBuckets map[string]*ratelimit.Bucket
	cfg             *config.Config
	logger          logger.Logger
}

func New(cfg *config.Config, logger logger.Logger) *RateManager {
	return &RateManager{
		cfg:             cfg,
		logger:          logger,
		ipBuckets:       make(map[string]*ratelimit.Bucket),
		loginBuckets:    make(map[string]*ratelimit.Bucket),
		passwordBuckets: make(map[string]*ratelimit.Bucket),
	}
}

//nolint
func (m *RateManager) Manage(ip string, login string, password string) bool {
	ipBucket, ok := m.ipBuckets[ip]
	if !ok {
		ipBucket = ratelimit.NewBucket(time.Second, int64(m.cfg.Data.Buckets.IpCapacity))
		m.ipBuckets[ip] = ipBucket
	}
	ipTockensUsed := ipBucket.TakeAvailable(1)
	if ipTockensUsed == 0 {
		return false
	}

	loginBucket, ok := m.loginBuckets[ip]
	if !ok {
		loginBucket = ratelimit.NewBucket(time.Second, int64(m.cfg.Data.Buckets.IpCapacity))
		m.loginBuckets[ip] = loginBucket
	}
	loginTockensUsed := loginBucket.TakeAvailable(1)
	if loginTockensUsed == 0 {
		return false
	}

	passwordBucket, ok := m.passwordBuckets[ip]
	if !ok {
		passwordBucket = ratelimit.NewBucket(time.Second, int64(m.cfg.Data.Buckets.IpCapacity))
		m.passwordBuckets[ip] = passwordBucket
	}
	passwordTockensUsed := passwordBucket.TakeAvailable(1)
	if passwordTockensUsed == 0 {
		return false
	}

	return true
}

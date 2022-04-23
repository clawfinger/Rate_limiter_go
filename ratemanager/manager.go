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
	ipBuckets       map[string]BucketData
	loginBuckets    map[string]BucketData
	passwordBuckets map[string]BucketData
	cfg             *config.Config
	logger          logger.Logger
}

func New(cfg *config.Config, logger logger.Logger) *RateManager {
	return &RateManager{
		cfg:             cfg,
		logger:          logger,
		ipBuckets:       make(map[string]BucketData),
		loginBuckets:    make(map[string]BucketData),
		passwordBuckets: make(map[string]BucketData),
	}
}

type Result struct {
	Ok     bool
	Reason string
}

func (m *RateManager) Manage(ip string, login string, password string) *Result {
	ipOk := m.ManageIP(ip)
	if !ipOk {
		return &Result{Ok: false, Reason: "IP denied"}
	}
	loginOk := m.ManageLogin(login)
	if !loginOk {
		return &Result{Ok: false, Reason: "Login denied"}
	}
	passOk := m.ManagePassword(password)
	if !passOk {
		return &Result{Ok: false, Reason: "Password denied"}
	}
	return &Result{
		Ok:     true,
		Reason: "success",
	}
}

func (m *RateManager) ManageIP(ip string) bool {
	ipBucketData, ok := m.ipBuckets[ip]
	if !ok {
		ipBucket := ratelimit.NewBucket(time.Second, int64(m.cfg.Data.Buckets.IpCapacity))
		ipBucketData = BucketData{
			Bucket:     ipBucket,
			LastActive: time.Now(),
		}
		m.ipBuckets[ip] = ipBucketData
	}
	ipTockensUsed := ipBucketData.Bucket.TakeAvailable(1)
	return ipTockensUsed != 0
}

func (m *RateManager) ManageLogin(login string) bool {
	loginBucketData, ok := m.loginBuckets[login]
	if !ok {
		ipBucket := ratelimit.NewBucket(time.Second, int64(m.cfg.Data.Buckets.LoginCapacity))
		loginBucketData = BucketData{
			Bucket:     ipBucket,
			LastActive: time.Now(),
		}
		m.loginBuckets[login] = loginBucketData
	}
	loginTockensUsed := loginBucketData.Bucket.TakeAvailable(1)
	return loginTockensUsed != 0
}

func (m *RateManager) ManagePassword(pass string) bool {
	passwordBucketData, ok := m.passwordBuckets[pass]
	if !ok {
		ipBucket := ratelimit.NewBucket(time.Second, int64(m.cfg.Data.Buckets.PasswordCapacity))
		passwordBucketData = BucketData{
			Bucket:     ipBucket,
			LastActive: time.Now(),
		}
		m.passwordBuckets[pass] = passwordBucketData
	}
	passwordTockensUsed := passwordBucketData.Bucket.TakeAvailable(1)
	return passwordTockensUsed != 0
}

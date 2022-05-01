package manager

import (
	"context"
	"sync"
	"time"

	"github.com/clawfinger/ratelimiter/config"
	"github.com/clawfinger/ratelimiter/logger"
	"github.com/juju/ratelimit"
)

type BucketData struct {
	Bucket     *ratelimit.Bucket
	LastActive time.Time
}

type AbstractRateManager interface {
	Manage(ip string, login string, password string) *Result
	DropStats(login string, ip string)
}

type RateManager struct {
	ipBuckets       map[string]BucketData
	loginBuckets    map[string]BucketData
	passwordBuckets map[string]BucketData
	cfg             *config.Config
	logger          logger.Logger
	mutex           sync.Mutex
	ctx             context.Context
	cancFunc        context.CancelFunc
}

func New(cfg *config.Config, logger logger.Logger) *RateManager {
	ctx, cancFunc := context.WithCancel(context.Background())
	return &RateManager{
		cfg:             cfg,
		logger:          logger,
		ipBuckets:       make(map[string]BucketData),
		loginBuckets:    make(map[string]BucketData),
		passwordBuckets: make(map[string]BucketData),
		ctx:             ctx,
		cancFunc:        cancFunc,
	}
}

func (m *RateManager) Start() {
	go func() {
		done := m.ctx.Done()
		ticker := time.NewTicker(time.Minute)
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				m.Cleanup()
			}
		}
	}()
}

func (m *RateManager) Stop() {
	m.cancFunc()
}

func (m *RateManager) Cleanup() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	deadline := m.cfg.Data.Buckets.BucketDecayTime
	now := time.Now()

	for ip, bucketData := range m.ipBuckets {
		diff := now.Sub(bucketData.LastActive)
		if diff >= deadline {
			delete(m.ipBuckets, ip)
			m.logger.Debug("Cleaning ip", ip)
		}
	}

	for login, bucketData := range m.loginBuckets {
		diff := now.Sub(bucketData.LastActive)
		if diff >= deadline {
			delete(m.loginBuckets, login)
			m.logger.Debug("Cleaning login", login)
		}
	}

	for password, bucketData := range m.passwordBuckets {
		diff := now.Sub(bucketData.LastActive)
		if diff >= deadline {
			delete(m.passwordBuckets, password)
			m.logger.Debug("Cleaning password", password)
		}
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
	m.mutex.Lock()
	defer m.mutex.Unlock()
	ipBucketData, ok := m.ipBuckets[ip]
	if !ok {
		ipBucket := ratelimit.NewBucket(time.Second, int64(m.cfg.Data.Buckets.IPCapacity))
		ipBucketData = BucketData{
			Bucket: ipBucket,
		}
		m.ipBuckets[ip] = ipBucketData
	}
	ipBucketData.LastActive = time.Now()
	ipTockensUsed := ipBucketData.Bucket.TakeAvailable(1)
	return ipTockensUsed != 0
}

func (m *RateManager) ManageLogin(login string) bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	loginBucketData, ok := m.loginBuckets[login]
	if !ok {
		loginBucket := ratelimit.NewBucket(time.Second, int64(m.cfg.Data.Buckets.LoginCapacity))
		loginBucketData = BucketData{
			Bucket: loginBucket,
		}
		m.loginBuckets[login] = loginBucketData
	}
	loginBucketData.LastActive = time.Now()
	loginTockensUsed := loginBucketData.Bucket.TakeAvailable(1)
	return loginTockensUsed != 0
}

func (m *RateManager) ManagePassword(pass string) bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	passwordBucketData, ok := m.passwordBuckets[pass]
	if !ok {
		passwordBucket := ratelimit.NewBucket(time.Second, int64(m.cfg.Data.Buckets.PasswordCapacity))
		passwordBucketData = BucketData{
			Bucket: passwordBucket,
		}
		m.passwordBuckets[pass] = passwordBucketData
	}
	passwordBucketData.LastActive = time.Now()
	passwordTockensUsed := passwordBucketData.Bucket.TakeAvailable(1)
	return passwordTockensUsed != 0
}

func (m *RateManager) DropStats(login string, ip string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	loginBucket := ratelimit.NewBucket(time.Second, int64(m.cfg.Data.Buckets.LoginCapacity))
	loginBucketData := BucketData{
		Bucket:     loginBucket,
		LastActive: time.Now(),
	}
	m.loginBuckets[login] = loginBucketData

	ipBucket := ratelimit.NewBucket(time.Second, int64(m.cfg.Data.Buckets.IPCapacity))
	ipBucketData := BucketData{
		Bucket:     ipBucket,
		LastActive: time.Now(),
	}
	m.ipBuckets[ip] = ipBucketData
}

package storage

import (
	"context"

	"github.com/clawfinger/ratelimiter/config"
	"github.com/clawfinger/ratelimiter/logger"
	"github.com/go-redis/redis/v8"
)

const (
	whitelistedVal = "whitelist"
	blacklistedVal = "blacklist"
)

type Restriction uint8

const (
	Not_set     Restriction = 0
	Whitelisted Restriction = 1
	Blacklisted Restriction = 2
)

type Result struct {
	Status Restriction
	Reason string
}

type AbstractStorage interface {
	CheckIP(ip string) *Result
	SetIP(ip string, restriction Restriction)
}

type Storage struct {
	client *redis.Client
	cfg    *config.Config
	logger logger.Logger
}

func New(cfg *config.Config, logger logger.Logger) *Storage {
	return &Storage{
		cfg:    cfg,
		logger: logger,
	}
}

func (s *Storage) Start() {
	s.client = redis.NewClient(&redis.Options{
		Addr:     s.cfg.Data.RedisData.Addr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

func (s *Storage) Stop() {
	s.client.Close()
}

func NewResult(status Restriction, reason string) *Result {
	return &Result{
		Status: status,
		Reason: reason,
	}
}

func (s *Storage) CheckIP(ip string) *Result {
	ctx, cancFunc := context.WithTimeout(context.Background(), s.cfg.Data.RedisData.OpTimeout)
	defer cancFunc()
	val, err := s.client.Get(ctx, ip).Result()
	if err != nil {
		if err == redis.Nil {
			return NewResult(Not_set, "not set")
		} else {
			s.logger.Error("Got error on redis read", err.Error())
			return NewResult(Blacklisted, "redis error")
		}
	}
	switch val {
	case whitelistedVal:
		return NewResult(Whitelisted, "whitelisted")
	case blacklistedVal:
		return NewResult(Blacklisted, "blacklisted")
	default:
		return NewResult(Blacklisted, "value is wrong")
	}
}

func (s *Storage) SetIP(ip string, restriction Restriction) {
	value := ""
	switch restriction {
	case Whitelisted:
		value = whitelistedVal
	case Blacklisted:
		value = blacklistedVal
	}

	ctx, cancFunc := context.WithTimeout(context.Background(), s.cfg.Data.RedisData.OpTimeout)
	defer cancFunc()
	err := s.client.Set(ctx, ip, value, 0).Err()
	if err != nil {
		s.logger.Error("Got error on redis write", err.Error())
	}
}

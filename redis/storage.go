package storage

import (
	"context"
	"net"

	"github.com/clawfinger/ratelimiter/config"
	"github.com/clawfinger/ratelimiter/logger"
	"github.com/go-redis/redis/v8"
)

const (
	whiteLabel = "whitelist"
	blackLabel = "blacklist"
	whiteKey   = "whiteKey"
	blackKey   = "blackKey"
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
	CheckIP(ctx context.Context, ip string) *Result
	SetIP(ctx context.Context, ip string, restriction Restriction)
	RemoveIP(ctx context.Context, ip string, restriction Restriction)
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

func (s *Storage) CheckBlacklist(ctx context.Context, ip net.IP) (bool, *Result) {
	ctx, cancFunc := context.WithTimeout(ctx, s.cfg.Data.RedisData.OpTimeout)
	result := s.client.HGetAll(ctx, blackKey)
	defer cancFunc()
	if result.Err() != nil {
		s.logger.Info("Error on gettting blacklist", result.Err().Error())
		return false, NewResult(Not_set, "db error")
	}
	for subnetString := range result.Val() {
		_, subnet, err := net.ParseCIDR(subnetString)
		if err != nil {
			s.logger.Info("Error parsing ip subnet string from db", err.Error())
			continue
		}
		if subnet.Contains(ip) {
			return true, NewResult(Blacklisted, "Blacklisted")
		}
	}
	return false, NewResult(Not_set, "not in a blacklist")
}

func (s *Storage) CheckWhitelistt(ctx context.Context, ip net.IP) (bool, *Result) {
	ctx, cancFunc := context.WithTimeout(ctx, s.cfg.Data.RedisData.OpTimeout)
	result := s.client.HGetAll(ctx, whiteKey)
	defer cancFunc()
	if result.Err() != nil {
		s.logger.Info("Error on gettting whitelist", result.Err().Error())
		return false, NewResult(Not_set, "db error")
	}
	for subnetString := range result.Val() {
		_, subnet, err := net.ParseCIDR(subnetString)
		if err != nil {
			s.logger.Info("Error parsing ip subnet string from db", err.Error())
			continue
		}
		if subnet.Contains(ip) {
			return true, NewResult(Whitelisted, "Whitelisted")
		}
	}

	return false, NewResult(Not_set, "not set in lists")
}

func (s *Storage) CheckIP(ctx context.Context, ipString string) *Result {
	ip := net.ParseIP(ipString)
	if ip == nil {
		s.logger.Info("Error parsing ip subnet string from db", ipString)
		return NewResult(Not_set, "error parsing ip")
	}

	filtered, result := s.CheckBlacklist(ctx, ip)
	if filtered {
		return result
	}
	filtered, result = s.CheckWhitelistt(ctx, ip)
	if filtered {
		return result
	}
	return NewResult(Not_set, "not set in lists")
}

func (s *Storage) SetIP(ctx context.Context, ip string, restriction Restriction) {
	ctx, cancFunc := context.WithTimeout(ctx, s.cfg.Data.RedisData.OpTimeout)
	defer cancFunc()
	switch restriction {
	case Whitelisted:
		result := s.client.HSet(ctx, whiteKey, ip, whiteLabel)
		if result.Err() != nil {
			s.logger.Info("Error on adding ip to whitelist", result.Err().Error())
			return
		}
		s.logger.Info("Added", result.Val(), "entries to whitelists")

	case Blacklisted:
		result := s.client.HSet(ctx, blackKey, ip, blackLabel)
		if result.Err() != nil {
			s.logger.Info("Error on adding ip to blacklist", result.Err().Error())
			return
		}
		s.logger.Info("Added", result.Val(), "entries to blacklist")
	}
}

func (s *Storage) RemoveIP(ctx context.Context, ip string, restriction Restriction) {
	ctx, cancFunc := context.WithTimeout(ctx, s.cfg.Data.RedisData.OpTimeout)
	defer cancFunc()
	switch restriction {
	case Whitelisted:
		result := s.client.HDel(ctx, whiteKey, ip)
		if result.Err() != nil {
			s.logger.Info("Error on removing ip from whitelist", result.Err().Error())
			return
		}
		s.logger.Info("Deleted", result.Val(), "entries from whitelists")

	case Blacklisted:
		result := s.client.HDel(ctx, blackKey, ip)
		if result.Err() != nil {
			s.logger.Info("Error on removing ip from blacklist", result.Err().Error())
			return
		}
		s.logger.Info("Deleted", result.Val(), "entries from blacklist")
	}
}

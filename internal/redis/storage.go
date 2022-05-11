package storage

import (
	"context"
	"errors"
	"net"

	"github.com/clawfinger/ratelimiter/config"
	internalapi "github.com/clawfinger/ratelimiter/internal/internalapi"
	"github.com/clawfinger/ratelimiter/internal/logger"
	"github.com/go-redis/redis/v8"
)

const (
	whiteLabel = "whitelist"
	blackLabel = "blacklist"
	whiteKey   = "whiteKey"
	blackKey   = "blackKey"
)

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

func NewResult(status internalapi.Restriction, reason string, err error) *internalapi.StorageResult {
	return &internalapi.StorageResult{
		Status: status,
		Reason: reason,
		Err:    err,
	}
}

func (s *Storage) CheckBlacklist(ctx context.Context, ip net.IP) (bool, *internalapi.StorageResult) {
	ctx, cancFunc := context.WithTimeout(ctx, s.cfg.Data.RedisData.OpTimeout)
	result := s.client.HGetAll(ctx, blackKey)
	defer cancFunc()
	if result.Err() != nil {
		s.logger.Info("Error on gettting blacklist", result.Err().Error())
		return false, NewResult(internalapi.NotSet, "db error", result.Err())
	}
	for subnetString := range result.Val() {
		_, subnet, err := net.ParseCIDR(subnetString)
		if err != nil {
			s.logger.Info("Error parsing ip subnet string from db", err.Error())
			continue
		}
		if subnet.Contains(ip) {
			return true, NewResult(internalapi.Blacklisted, "Blacklisted", nil)
		}
	}
	return false, NewResult(internalapi.NotSet, "not in a blacklist", nil)
}

func (s *Storage) CheckWhitelistt(ctx context.Context, ip net.IP) (bool, *internalapi.StorageResult) {
	ctx, cancFunc := context.WithTimeout(ctx, s.cfg.Data.RedisData.OpTimeout)
	result := s.client.HGetAll(ctx, whiteKey)
	defer cancFunc()
	if result.Err() != nil {
		s.logger.Info("Error on gettting whitelist", result.Err().Error())
		return false, NewResult(internalapi.NotSet, "db error", result.Err())
	}
	for subnetString := range result.Val() {
		_, subnet, err := net.ParseCIDR(subnetString)
		if err != nil {
			s.logger.Info("Error parsing ip subnet string from db", err.Error())
			continue
		}
		if subnet.Contains(ip) {
			return true, NewResult(internalapi.Whitelisted, "Whitelisted", nil)
		}
	}

	return false, NewResult(internalapi.NotSet, "not set in lists", nil)
}

func (s *Storage) CheckIP(ctx context.Context, ipString string) *internalapi.StorageResult {
	ip := net.ParseIP(ipString)
	if ip == nil {
		s.logger.Info("Error parsing ip subnet string from db", ipString)
		return NewResult(internalapi.NotSet, "error parsing ip", errors.New("error parsing ip"))
	}

	filtered, result := s.CheckBlacklist(ctx, ip)
	if result.Err != nil || filtered {
		return result
	}

	filtered, result = s.CheckWhitelistt(ctx, ip)
	if result.Err != nil || filtered {
		return result
	}
	return NewResult(internalapi.NotSet, "not set in lists", nil)
}

func (s *Storage) SetIP(ctx context.Context, ip string, restriction internalapi.Restriction) error {
	ctx, cancFunc := context.WithTimeout(ctx, s.cfg.Data.RedisData.OpTimeout)
	defer cancFunc()
	switch restriction {
	case internalapi.Whitelisted:
		result := s.client.HSet(ctx, whiteKey, ip, whiteLabel)
		if result.Err() != nil {
			s.logger.Info("Error on adding ip to whitelist", result.Err().Error())
			return result.Err()
		}
		s.logger.Info("Added", result.Val(), "entries to whitelists")

	case internalapi.Blacklisted:
		result := s.client.HSet(ctx, blackKey, ip, blackLabel)
		if result.Err() != nil {
			s.logger.Info("Error on adding ip to blacklist", result.Err().Error())
			return result.Err()
		}
		s.logger.Info("Added", result.Val(), "entries to blacklist")
	case internalapi.NotSet:
		return nil
	}
	return nil
}

func (s *Storage) RemoveIP(ctx context.Context, ip string, restriction internalapi.Restriction) error {
	ctx, cancFunc := context.WithTimeout(ctx, s.cfg.Data.RedisData.OpTimeout)
	defer cancFunc()
	switch restriction {
	case internalapi.Whitelisted:
		result := s.client.HDel(ctx, whiteKey, ip)
		if result.Err() != nil {
			s.logger.Info("Error on removing ip from whitelist", result.Err().Error())
			return result.Err()
		}
		s.logger.Info("Deleted", result.Val(), "entries from whitelists")
	case internalapi.Blacklisted:
		result := s.client.HDel(ctx, blackKey, ip)
		if result.Err() != nil {
			s.logger.Info("Error on removing ip from blacklist", result.Err().Error())
			return result.Err()
		}
		s.logger.Info("Deleted", result.Val(), "entries from blacklist")
	case internalapi.NotSet:
		return nil
	}
	return nil
}

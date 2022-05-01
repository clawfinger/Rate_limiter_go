package test

import (
	"testing"
	"time"

	"github.com/clawfinger/ratelimiter/config"
	manager "github.com/clawfinger/ratelimiter/ratemanager"
	"github.com/stretchr/testify/require"
)

const (
	ip    = "127.0.0.1"
	login = "Vasya"
	pass  = "12345"
)

func TestIPLimit(t *testing.T) {
	config := config.NewConfig()
	config.Data.Buckets.IPCapacity = 10
	config.Data.Buckets.LoginCapacity = 100
	config.Data.Buckets.PasswordCapacity = 100

	loggerMock := &LoggerMock{}

	manager := manager.New(config, loggerMock)
	for i := 0; i < 10; i++ {
		res := manager.Manage(ip, login, pass)
		require.Truef(t, res.Ok, res.Reason)
	}
	res := manager.Manage(ip, login, pass)
	require.Falsef(t, res.Ok, res.Reason)
}

func TestLoginLimit(t *testing.T) {
	config := config.NewConfig()
	config.Data.Buckets.IPCapacity = 1000
	config.Data.Buckets.LoginCapacity = 10
	config.Data.Buckets.PasswordCapacity = 50

	loggerMock := &LoggerMock{}

	manager := manager.New(config, loggerMock)
	for i := 0; i < 10; i++ {
		res := manager.Manage(ip, login, pass)
		require.Truef(t, res.Ok, res.Reason)
	}
	res := manager.Manage(ip, login, pass)
	require.Falsef(t, res.Ok, res.Reason)
}

func TestPasswordLimit(t *testing.T) {
	config := config.NewConfig()
	config.Data.Buckets.IPCapacity = 1000
	config.Data.Buckets.LoginCapacity = 100
	config.Data.Buckets.PasswordCapacity = 10

	loggerMock := &LoggerMock{}

	manager := manager.New(config, loggerMock)
	for i := 0; i < 10; i++ {
		res := manager.Manage(ip, login, pass)
		require.Truef(t, res.Ok, res.Reason)
	}
	res := manager.Manage(ip, login, pass)
	require.Falsef(t, res.Ok, res.Reason)
}

func TestAcceptAfterDeny(t *testing.T) {
	config := config.NewConfig()
	config.Data.Buckets.IPCapacity = 1000
	config.Data.Buckets.LoginCapacity = 10
	config.Data.Buckets.PasswordCapacity = 50

	loggerMock := &LoggerMock{}

	manager := manager.New(config, loggerMock)
	for i := 0; i < 10; i++ {
		res := manager.Manage(ip, login, pass)
		require.Truef(t, res.Ok, res.Reason)
	}
	res := manager.Manage(ip, login, pass)
	require.Falsef(t, res.Ok, res.Reason)

	time.Sleep(time.Second)

	res = manager.Manage(ip, login, pass)
	require.Truef(t, res.Ok, res.Reason)
}

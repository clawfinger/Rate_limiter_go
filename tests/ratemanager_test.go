package test

import (
	"testing"

	"github.com/clawfinger/ratelimiter/config"
	manager "github.com/clawfinger/ratelimiter/ratemanager"
	"github.com/stretchr/testify/require"
)

func TestIPLimit(t *testing.T) {
	config := config.NewConfig()
	config.Data.Buckets.IpCapacity = 10
	config.Data.Buckets.LoginCapacity = 100
	config.Data.Buckets.PasswordCapacity = 100

	loggerMock := &LoggerMock{}

	ip := "127.0.0.1"
	login := "Vasya"
	pass := "12345"
	manager := manager.New(config, loggerMock)
	for i := 0; i < 10; i++ {
		require.Truef(t, manager.Manage(ip, login, pass), "Denied illegaly")
	}
	require.Falsef(t, manager.Manage(ip, login, pass), "Passed illegaly")
}
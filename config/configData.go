package config

import (
	"time"

	"github.com/spf13/viper"
)

type Data struct {
	DBData    DatabaseData
	Logger    LoggerConf
	Grpc      Grpc
	Buckets   Buckets
	RedisData RedisData
}

func newConfigData() *Data {
	return &Data{}
}

func (d *Data) SetDefault(v *viper.Viper) {
	d.Logger.SetDefault(v)
	d.Grpc.SetDefault(v)
	d.Buckets.SetDefault(v)
	d.RedisData.SetDefault(v)
}

type LoggerConf struct {
	Level    string
	Filename string
}

func (d *LoggerConf) SetDefault(v *viper.Viper) {
	v.SetDefault("Logger", map[string]interface{}{
		"Level":    "debug",
		"Filename": "limiter.log",
	})
}

type DatabaseData struct {
	Username string
	Password string
}

type RedisData struct {
	Addr      string
	OpTimeout time.Duration
}

func (d *RedisData) SetDefault(v *viper.Viper) {
	v.SetDefault("RedisData", map[string]interface{}{
		"Addr":      "localhost:6379",
		"OpTimeout": "500ms",
	})
}

type Grpc struct {
	Addr string
}

func (d *Grpc) SetDefault(v *viper.Viper) {
	v.SetDefault("Grpc", map[string]interface{}{
		"Addr": "127.0.0.1:50051",
	})
}

type Buckets struct {
	IPCapacity       uint
	LoginCapacity    uint
	PasswordCapacity uint
	BucketDecayTime  time.Duration
}

func (b *Buckets) SetDefault(v *viper.Viper) {
	v.SetDefault("Buckets", map[string]interface{}{
		"IPCapacity":       "1000",
		"LoginCapacity":    "10",
		"PasswordCapacity": "100",
		"BucketDecayTime":  "1m",
	})
}

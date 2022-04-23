package config

import (
	"github.com/spf13/viper"
)

type Data struct {
	DBData  DatabaseData
	Logger  LoggerConf
	Storage Storage
	Grpc    Grpc
}

func newConfigData() *Data {
	return &Data{}
}

func (d *Data) SetDefault(v *viper.Viper) {
	d.Logger.SetDefault(v)
	d.Storage.SetDefault(v)
	d.Grpc.SetDefault(v)
}

type LoggerConf struct {
	Level    string
	Filename string
}

func (d *LoggerConf) SetDefault(v *viper.Viper) {
	v.SetDefault("Logger", map[string]interface{}{
		"Level":    "info",
		"Filename": "limiter.log",
	})
}

type DatabaseData struct {
	Username string
	Password string
}

type Storage struct {
	Type string
}

func (d *Storage) SetDefault(v *viper.Viper) {
	v.SetDefault("storage", map[string]interface{}{
		"Type": "inmemory",
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

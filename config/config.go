package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Data *Data
}

func NewConfig() *Config {
	data := newConfigData()

	return &Config{
		Data: data,
	}
}

func (c *Config) Init(cfgFilePath string) error {
	c.Data.SetDefault(viper.GetViper())
	viper.SetConfigType("json")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	viper.SetConfigFile(cfgFilePath)

	if err := viper.ReadInConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "Error on reading config, Reason: %s", err.Error())
	}

	err := viper.Unmarshal(c.Data)
	if err != nil {
		return err
	}

	err = viper.WriteConfig()
	c.Data.DBData.Username = viper.GetString("dbdata.username")
	c.Data.DBData.Password = viper.GetString("dbdata.password")
	return err
}

// TODO

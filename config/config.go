package config

import (
	"github.com/spf13/viper"
)

func LoadConfig(configFile string) error {

	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("/etc/tinykv/")
		viper.AddConfigPath("$HOME/tinykv/")
		viper.AddConfigPath(".")
	}

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	c.addr = viper.GetString("addr")
	c.peers = viper.GetStringSlice("peers")

	return nil
}

type config struct {
	addr  string
	peers []string
}

var c config

func Addr() string {
	return c.addr
}

func Peers() []string {
	return c.peers
}

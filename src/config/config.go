package config

import (
	"fmt"
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	Postgres PostgresConfig `json:"postgres"`
	Token    TokenConfig    `json:"token"`
}

type PostgresConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
	SSLMode  string `json:"sslMode"`
}

type TokenConfig struct {
	TokenKey      string        `json:"tokenKey"`
	TokenDuration time.Duration `json:"tokenDuration"`
	Address       string        `json:"address"`
}

func LoadConfig() *Config {
	viper.AddConfigPath("conf")
	viper.SetConfigName("config-dev.json")
	viper.SetConfigType("json")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	config := &Config{}
	err = viper.Unmarshal(config)
	if err != nil {
		panic(fmt.Errorf("fatal error unmarshal config: %s", err))
	}
	return config
}

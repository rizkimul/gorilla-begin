package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DatabaseUser    string        `mapstructure:"DATABASE_USER"`
	DatabasePass    string        `mapstructure:"DATABASE_PASSWORD"`
	DatabaseName    string        `mapstructure:"DATABASE_NAME"`
	SslMode         string        `mapstructure:"SSL_MODE"`
	CloudSecretKey  string        `mapstructure:"CLOUD_SECRET_KEY"`
	AccessTokenExp  time.Duration `mapstructure:"ACCESS_TOKEN_EXPIRED_IN"`
	RefreshTokenExp time.Duration `mapstructure:"REFRESH_TOKEN_EXPIRED_IN"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName(".env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)

	return
}

package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	RedisAddress      string        `mapstructure:"REDIS_ADDRESS"`
	RedisPassword     string        `mapstructure:"REDIS_PASSWORD"`
	RedisDB           int           `mapstructure:"REDIS_DB"`
	TokenDuration     time.Duration `mapstructure:"TOKEN_DUDATION"`
	SymmetricKey      string        `mapstructure:"SYMMETRIC_KEY"`
	GRPCAddress       string        `mapstructure:"GRPC_SERVER_ADDRESS"`
	UserServiceAddr   string        `mapstructure:"USER_SERVICE_ADDRESS"`
	Environment       string        `mapstructure:"ENVIRONMENT""`
	ServiceName       string        `mapstructure:"SERVICE_NAME"`
	OtelCollectorAddr string        `mapstructure:"OTEL_COLLECTOR_STRUCTURE"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return config, err
		}
	}

	viper.AutomaticEnv()
	viper.BindEnv("REDIS_ADDRESS")
	viper.BindEnv("REDIS_PASSWORD")
	viper.BindEnv("REDIS_DB")
	viper.BindEnv("TOKEN_DUDATION")
	viper.BindEnv("GRPC_SERVER_ADDRESS")
	viper.BindEnv("USER_SERVICE_ADDRESS")
	viper.BindEnv("SYMMETRIC_KEY")
	viper.BindEnv("ENVIRONMENT")
	viper.BindEnv("SERVICE_NAME")
	viper.BindEnv("OTEL_COLLECTOR_STRUCTURE")

	err = viper.Unmarshal(&config)
	return config, err
}

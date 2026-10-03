package config

import (
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config struct maps your environment variables
type Config struct {
	ServiceName string `mapstructure:"SERVICE_NAME"`
	ServerPort  string `mapstructure:"SERVICE_PORT"`
	MongoURI    string `mapstructure:"MONGO_URI"`
	RedisAddr   string `mapstructure:"REDIS_ADDR"`
}

// LoadConfig reads configuration from file and/or environment variables
func LoadConfig(path string) (Config, error) {
	v := viper.New()
	configFile := filepath.Join(path, ".env")
	v.SetConfigFile(configFile)
	v.SetConfigType("env")
	v.SetDefault("SERVICE_PORT", "8089")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		var notFound viper.ConfigFileNotFoundError
		if !errors.As(err, &notFound) && !errors.Is(err, fs.ErrNotExist) {
			return Config{}, fmt.Errorf("read config file %s: %w", configFile, err)
		}
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return Config{}, fmt.Errorf("decode config: %w", err)
	}

	return config, nil
}

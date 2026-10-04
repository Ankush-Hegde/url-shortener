package config

import (
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

// Config struct maps your environment variables
type Config struct {
	ServiceName   string `mapstructure:"SERVICE_NAME"`
	ServerPort    string `mapstructure:"SERVICE_PORT"`
	PublicBaseURL string `mapstructure:"PUBLIC_BASE_URL"`
	MongoURI      string `mapstructure:"MONGODB_CONNECTION_STRING"`
	MongoUsername string `mapstructure:"MONGODB_USERNAME"`
	MongoPassword string `mapstructure:"MONGODB_PASSWORD"`
	RedisAddr     string `mapstructure:"REDIS_HOST"`
	RedisUsername string `mapstructure:"REDIS_USERNAME"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`
	RedisDB       int    `mapstructure:"-"`
	RedisTLS      bool   `mapstructure:"REDIS_TLS"`
}

// LoadConfig reads configuration from file and/or environment variables
func LoadConfig(path string) (Config, error) {
	v := viper.New()
	configFile := filepath.Join(path, ".env")
	v.SetConfigFile(configFile)
	v.SetConfigType("env")
	v.SetDefault("SERVICE_PORT", "8089")
	v.SetDefault("PUBLIC_BASE_URL", "")
	v.SetDefault("REDIS_DB", 0)
	v.SetDefault("REDIS_TLS", false)
	v.AutomaticEnv()

	for _, key := range []string{
		"SERVICE_NAME",
		"SERVICE_PORT",
		"PUBLIC_BASE_URL",
		"MONGODB_CONNECTION_STRING",
		"MONGODB_USERNAME",
		"MONGODB_PASSWORD",
		"REDIS_HOST",
		"REDIS_USERNAME",
		"REDIS_PASSWORD",
		"REDIS_DB",
		"REDIS_TLS",
	} {
		if err := v.BindEnv(key); err != nil {
			return Config{}, fmt.Errorf("bind environment variable %s: %w", key, err)
		}
	}

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

	redisDB := strings.TrimSpace(v.GetString("REDIS_DB"))
	if redisDB != "" {
		parsedRedisDB, err := strconv.Atoi(redisDB)
		if err != nil {
			return Config{}, fmt.Errorf("REDIS_DB must be a non-negative integer")
		}
		config.RedisDB = parsedRedisDB
	}
	if config.RedisDB < 0 {
		return Config{}, fmt.Errorf("REDIS_DB must be a non-negative integer")
	}

	if strings.TrimSpace(config.PublicBaseURL) == "" {
		config.PublicBaseURL = "http://localhost:" + config.ServerPort
	}

	return config, nil
}

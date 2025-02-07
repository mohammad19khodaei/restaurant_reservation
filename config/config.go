package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App struct {
		TestingMode     bool          `mapstructure:"testing_mode"`
		Address         string        `mapstructure:"address"`
		ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
		SecretKey       string        `mapstructure:"secret_key"`
		TokenDuration   time.Duration `mapstructure:"token_duration"`
	} `mapstructure:"app"`
	Database struct {
		Host     string `mapstructure:"host"`
		Port     string `mapstructure:"port"`
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
		Name     string `mapstructure:"name"`
	} `mapstructure:"db"`
}

func LoadConfig(path string, filename string) (*Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(filename)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("No .env file found, using environment variables instead: %v", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

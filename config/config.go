package config

import (
	"fmt"
	"github.com/spf13/viper"
)

// AppConfig holds the application configuration.
type AppConfig struct {
	App      App      `mapstructure:"app"`
	Database Database `mapstructure:"database"`
	Redis    Redis    `mapstructure:"redis"`
	JWT      JWT      `mapstructure:"jwt"`
}

// App settings
type App struct {
	Name        string `mapstructure:"name"`
	Port        string `mapstructure:"port"`
	Environment string `mapstructure:"environment"`
}

// Database settings
type Database struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

// Redis settings
type Redis struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
	DB   int    `mapstructure:"db"`
}

// JWT settings
type JWT struct {
	SecretKey               string `mapstructure:"secret_key"`
	TokenExpiryMinutes      int    `mapstructure:"token_expiry_minutes"`
	RefreshTokenExpiryHours int    `mapstructure:"refresh_token_expiry_hours"`
}

// LoadConfig loads the application configuration from a YAML file.
func LoadConfig(path string) (*AppConfig, error) {
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml") // or "yml"

	viper.AutomaticEnv() // Automatically read environment variables

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, fmt.Errorf("config file not found: %w", err)
		} else {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	var config AppConfig
	err := viper.Unmarshal(&config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

package config

import (
	"github.com/spf13/viper"
	"log"
)

func Init() {
	viper.SetConfigName(".env") // Load .env file
	viper.SetConfigType("env")
	viper.AddConfigPath(".") // Look for .env in the root directory
	viper.AutomaticEnv()     // Override with environment variables

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}
}

func Get(key string) string {
	return viper.GetString(key)
}

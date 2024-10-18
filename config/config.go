package config

import (
	"log"

	"github.com/spf13/viper"
)

// InitConfig initializes the application configuration
func InitConfig() {
	viper.SetConfigName("config") // config.yaml
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")
	viper.AutomaticEnv() // Override config values with environment variables

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("No config file found: %v", err)
	}

	log.Println("Configuration loaded successfully.")
}

// GetEnv fetches environment variable with a fallback
func GetEnv(key string, fallback string) string {
	if value := viper.GetString(key); value != "" {
		return value
	}
	return fallback
}

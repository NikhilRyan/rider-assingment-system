package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	DB    DBConfig
	Redis RedisConfig
}

type DBConfig struct {
	User     string
	Password string
	DBName   string
	SSLMode  string
	Host     string
	Port     string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

var Cfg *Config

func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	err = viper.Unmarshal(&Cfg)
	if err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}
}

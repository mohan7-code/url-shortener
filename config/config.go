package config

import (
	"errors"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort   string
	DatabaseUrl  string
	MaxDBConn    int
	BaseShortURL string

	RedisURL        string
	RedisDB         string
	RedisMaxRetries string
	RedisPoolSize   string
}

var AppConfig *Config

func LoadConfig(envFile string) (*Config, error) {

	if err := godotenv.Load(envFile); err != nil {
		log.Println("error reading the env file")
	}

	cfg := &Config{}

	// Server port
	cfg.ServerPort = os.Getenv("SERVER_PORT")
	if cfg.ServerPort == "" {
		cfg.ServerPort = "8080"
	}

	// Database URL
	cfg.DatabaseUrl = os.Getenv("DATABASE_URL")
	if cfg.DatabaseUrl == "" {
		return nil, errors.New("missing environment variable")
	}

	// Max DB connections
	maxDBConn := os.Getenv("MAX_DB_CONN")
	if maxDBConn == "" {
		cfg.MaxDBConn = 10
	} else {
		dbMaxConn, err := strconv.Atoi(maxDBConn)
		if err != nil {
			log.Println("Invalid DB connections")
			cfg.MaxDBConn = 10
		} else {
			cfg.MaxDBConn = dbMaxConn
		}
	}

	cfg.BaseShortURL = os.Getenv("BASE_SHORT_URL")
	if cfg.BaseShortURL == "" {
		cfg.BaseShortURL = "https://sho.rt"
	}

	cfg.RedisURL = os.Getenv("REDIS_URL")
	cfg.RedisDB = os.Getenv("REDIS_DB")
	cfg.RedisMaxRetries = os.Getenv("REDIS_MAX_RETRIES")
	cfg.RedisPoolSize = os.Getenv("REDIS_POOL_SIZE")

	AppConfig = cfg
	return cfg, nil
}

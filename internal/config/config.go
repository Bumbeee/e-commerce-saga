package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerAddress string
	ReadTimeout   time.Duration
	WriteTimeout  time.Duration

	DatabaseDSN       string
	DBMaxConns        int
	DBMinConns        int
	DBConnMaxIdleTime time.Duration

	LogLevel    string
	ServiceName string

	KafkaBrokers []string
}

// Load читает переменные окружения, парсит их и валидирует.
func Load() (*Config, error) {
	_ = godotenv.Load()
	cfg := &Config{
		ServerAddress: getEnv("SERVER_ADDRESS", ":8080"),

		DatabaseDSN: os.Getenv("DATABASE_DSN"),

		LogLevel:    getEnv("LOG_LEVEL", "info"),
		ServiceName: getEnv("SERVICE_NAME", "order-service"),
	}

	var err error
	cfg.ReadTimeout, err = parseDuration("READ_TIMEOUT", "10s")
	if err != nil {
		return nil, err
	}
	cfg.WriteTimeout, err = parseDuration("WRITE_TIMEOUT", "10s")
	if err != nil {
		return nil, err
	}
	cfg.DBMaxConns, err = getIntEnv("DB_MAX_CONNS", "10")
	if err != nil {
		return nil, err
	}
	cfg.DBMaxConns, err = getIntEnv("DB_MIN_CONNS", "2")
	if err != nil {
		return nil, err
	}
	cfg.DBConnMaxIdleTime, err = parseDuration("DB_CONN_MAX_IDLE_TIME", "30m")
	if err != nil {
		return nil, err
	}
	brokersEnv := os.Getenv("KAFKA_BROKERS")
	if brokersEnv != "" {
		cfg.KafkaBrokers = strings.Split(brokersEnv, ",")
	}

	return cfg, cfg.validate()
}

func (c *Config) validate() error {
	if c.DatabaseDSN == "" {
		return fmt.Errorf("DATABASE_DSN is required")
	}
	return nil
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

func getIntEnv(key, defaultVal string) (int, error) {
	val := getEnv(key, defaultVal)
	d, err := strconv.Atoi(val)
	if err != nil {
		return 0, fmt.Errorf("invalid %s value %q: %w", key, val, err)
	}
	return d, nil
}

func parseDuration(key, defaultVal string) (time.Duration, error) {
	val := getEnv(key, defaultVal)
	d, err := time.ParseDuration(val)
	if err != nil {
		return 0, fmt.Errorf("invalid %s duration %q: %w", key, val, err)
	}
	return d, nil
}

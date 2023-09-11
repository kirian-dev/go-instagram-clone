package config

import (
	"errors"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppVersion        string
	Port              string
	Mode              string
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	CxtDefaultTimeout string
	Debug             string
	JwtSecretKey      string
	PostgresDbname    string
	PostgresUser      string
	PostgresPassword  string
	PostgresPort      string
	PostgresHost      string
	PostgresSslMode   string
	PgDriver          string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load("./config/.env"); err != nil {
		return nil, err
	}

	cfg := &Config{
		AppVersion:        os.Getenv("AppVersion"),
		Port:              os.Getenv("Port"),
		Mode:              os.Getenv("Mode"),
		CxtDefaultTimeout: os.Getenv("CxtDefaultTimeout"),
		Debug:             os.Getenv("Debug"),
		JwtSecretKey:      os.Getenv("JwtSecretKey"),
		PostgresDbname:    os.Getenv("PostgresDbname"),
		PostgresUser:      os.Getenv("PostgresUser"),
		PostgresPassword:  os.Getenv("PostgresPassword"),
		PostgresPort:      os.Getenv("PostgresPort"),
		PostgresHost:      os.Getenv("PostgresHost"),
		PostgresSslMode:   os.Getenv("PostgresSslMode"),
		PgDriver:          os.Getenv("PgDriver"),
	}

	readTimeoutStr := os.Getenv("ReadTimeout")
	writeTimeoutStr := os.Getenv("WriteTimeout")

	var err error
	cfg.ReadTimeout, err = time.ParseDuration(readTimeoutStr)
	if err != nil {
		return nil, err
	}

	cfg.WriteTimeout, err = time.ParseDuration(writeTimeoutStr)
	if err != nil {
		return nil, err
	}

	if err := validateConfig(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func validateConfig(cfg *Config) error {
	fields := []interface{}{
		cfg.AppVersion, cfg.Port, cfg.Mode, cfg.CxtDefaultTimeout, cfg.Debug,
		cfg.JwtSecretKey, cfg.PostgresDbname, cfg.PostgresUser, cfg.PostgresPassword,
		cfg.PostgresPort, cfg.PostgresHost, cfg.PostgresSslMode, cfg.PgDriver,
	}

	for _, field := range fields {
		if field == "" {
			return errors.New("configuration variables must not be missing or empty")
		}
	}

	return nil
}

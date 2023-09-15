package config

import (
	"fmt"
	"os"
	"strconv"
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
	MaxOpenConns      int
	MaxIdleConns      int
	ConnMaxLifetime   time.Duration
	ConnMaxIdleTime   time.Duration
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(".env"); err != nil {
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

	cfg.MaxOpenConns, err = strconv.Atoi(os.Getenv("MaxOpenConns"))
	if err != nil {
		return nil, err
	}

	cfg.MaxIdleConns, err = strconv.Atoi(os.Getenv("MaxIdleConns"))
	if err != nil {
		return nil, err
	}

	cfg.ConnMaxLifetime, err = time.ParseDuration(os.Getenv("ConnMaxLifetime"))
	if err != nil {
		return nil, err
	}

	cfg.ConnMaxIdleTime, err = time.ParseDuration(os.Getenv("ConnMaxIdleTime"))
	if err != nil {
		return nil, err
	}

	if err := validateConfig(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func validateConfig(cfg *Config) error {
	requiredFields := []struct {
		value interface{}
		name  string
	}{
		{cfg.AppVersion, "AppVersion"},
		{cfg.Port, "Port"},
		{cfg.Mode, "Mode"},
		{cfg.CxtDefaultTimeout, "CxtDefaultTimeout"},
		{cfg.Debug, "Debug"},
		{cfg.JwtSecretKey, "JwtSecretKey"},
		{cfg.PostgresDbname, "PostgresDbname"},
		{cfg.PostgresUser, "PostgresUser"},
		{cfg.PostgresPassword, "PostgresPassword"},
		{cfg.PostgresPort, "PostgresPort"},
		{cfg.PostgresHost, "PostgresHost"},
		{cfg.PostgresSslMode, "PostgresSslMode"},
		{cfg.PgDriver, "PgDriver"},
		{cfg.MaxOpenConns, "MaxOpenConns"},
		{cfg.MaxIdleConns, "MaxIdleConns"},
		{cfg.ConnMaxLifetime, "ConnMaxLifetime"},
		{cfg.ConnMaxIdleTime, "ConnMaxIdleTime"},
	}

	for _, field := range requiredFields {
		if field.value == "" {
			return fmt.Errorf("configuration variable '%s' must not be missing or empty", field.name)
		}
	}

	return nil
}

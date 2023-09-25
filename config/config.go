package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppVersion        string
	ChatPort          string
	AnalyticsPort     string
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
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(".env"); err != nil {
		return nil, err
	}

	cfg := &Config{
		AppVersion:        os.Getenv("AppVersion"),
		ChatPort:          os.Getenv("ChatPort"),
		AnalyticsPort:     os.Getenv("AnalyticsPort"),
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
	}

	var err error
	cfg.ReadTimeout, err = time.ParseDuration(os.Getenv("ReadTimeout"))
	if err != nil {
		return nil, err
	}

	cfg.WriteTimeout, err = time.ParseDuration(os.Getenv("WriteTimeout"))
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
		{cfg.ChatPort, "ChatPort"},
		{cfg.AnalyticsPort, "AnalyticsPort"},
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
	}

	for _, field := range requiredFields {
		if field.value == "" {
			return fmt.Errorf("configuration variable '%s' must not be missing or empty", field.name)
		}
	}

	return nil
}

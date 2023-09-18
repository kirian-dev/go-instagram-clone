package config

import (
	"fmt"
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
	Driver            string
	URI               string
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
		URI:               os.Getenv("URI"),
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
		{cfg.Port, "Port"},
		{cfg.Mode, "Mode"},
		{cfg.CxtDefaultTimeout, "CxtDefaultTimeout"},
		{cfg.Debug, "Debug"},
		{cfg.JwtSecretKey, "JwtSecretKey"},
		{cfg.URI, "URI"},
	}

	for _, field := range requiredFields {
		if field.value == "" {
			return fmt.Errorf("configuration variable '%s' must not be missing or empty", field.name)
		}
	}

	return nil
}

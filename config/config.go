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
	MySQLUser         string
	MySQLPassword     string
	MySQLHost         string
	MySQLPort         string
	MySQLDBName       string
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
		MySQLUser:         os.Getenv("MySQLUser"),
		MySQLPassword:     os.Getenv("MySQLPassword"),
		MySQLHost:         os.Getenv("MySQLHost"),
		MySQLPort:         os.Getenv("MySQLPort"),
		MySQLDBName:       os.Getenv("MySQLDBName"),
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
		{cfg.MySQLDBName, "MySQLDBName"},
		{cfg.MySQLUser, "MySQLUser"},
		{cfg.MySQLPassword, "MySQLPassword"},
		{cfg.MySQLPort, "MySQLPort"},
		{cfg.MySQLHost, "MySQLHost"},
	}

	for _, field := range requiredFields {
		if field.value == "" {
			return fmt.Errorf("configuration variable '%s' must not be missing or empty", field.name)
		}
	}

	return nil
}

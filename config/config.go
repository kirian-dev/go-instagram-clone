package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppVersion                string
	ChatPort                  string
	AnalyticsPort             string
	Mode                      string
	ReadTimeout               time.Duration
	WriteTimeout              time.Duration
	CxtDefaultTimeout         string
	Debug                     string
	JwtSecretKey              string
	PostgresDbname            string
	PostgresUser              string
	PostgresPassword          string
	PostgresPort              string
	PostgresHost              string
	PostgresSslMode           string
	AnalyticsPostgresUser     string
	AnalyticsPostgresPassword string
	AnalyticsPostgresHost     string
	AnalyticsPostgresPort     string
	AnalyticsPostgresDBName   string
	AnalyticsPostgresSslMode  string
	EmailFrom                 string
	SMTPHost                  string
	SMTPPassword              string
	SMTPPort                  string
	SMTPUser                  string
	ClientOrigin              string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(".env"); err != nil {
		return nil, err
	}

	cfg := &Config{
		AppVersion:                os.Getenv("AppVersion"),
		ChatPort:                  os.Getenv("ChatPort"),
		AnalyticsPort:             os.Getenv("AnalyticsPort"),
		Mode:                      os.Getenv("Mode"),
		CxtDefaultTimeout:         os.Getenv("CxtDefaultTimeout"),
		Debug:                     os.Getenv("Debug"),
		JwtSecretKey:              os.Getenv("JwtSecretKey"),
		PostgresDbname:            os.Getenv("PostgresDbname"),
		PostgresUser:              os.Getenv("PostgresUser"),
		PostgresPassword:          os.Getenv("PostgresPassword"),
		PostgresPort:              os.Getenv("PostgresPort"),
		PostgresHost:              os.Getenv("PostgresHost"),
		PostgresSslMode:           os.Getenv("PostgresSslMode"),
		AnalyticsPostgresUser:     os.Getenv("AnalyticsPostgresUser"),
		AnalyticsPostgresPassword: os.Getenv("AnalyticsPostgresPassword"),
		AnalyticsPostgresHost:     os.Getenv("AnalyticsPostgresHost"),
		AnalyticsPostgresPort:     os.Getenv("AnalyticsPostgresPort"),
		AnalyticsPostgresDBName:   os.Getenv("AnalyticsPostgresDBName"),
		AnalyticsPostgresSslMode:  os.Getenv("AnalyticsPostgresSslMode"),
		EmailFrom:                 os.Getenv("EmailFrom"),
		SMTPHost:                  os.Getenv("SMTPHost"),
		SMTPPassword:              os.Getenv("SMTPPassword"),
		SMTPPort:                  os.Getenv("SMTPPort"),
		SMTPUser:                  os.Getenv("SMTPUser"),
		ClientOrigin:              os.Getenv("ClientOrigin"),
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
		{cfg.AnalyticsPostgresDBName, "AnalyticsPostgresDBName"},
		{cfg.AnalyticsPostgresUser, "AnalyticsPostgresUser"},
		{cfg.AnalyticsPostgresPassword, "AnalyticsPostgresPassword"},
		{cfg.AnalyticsPostgresPort, "AnalyticsPostgresPort"},
		{cfg.AnalyticsPostgresHost, "AnalyticsPostgresHost"},
		{cfg.AnalyticsPostgresSslMode, "AnalyticsPostgresSslMode"},
		{cfg.EmailFrom, "EmailFrom"},
		{cfg.SMTPHost, "SMTPHost"},
		{cfg.SMTPPassword, "SMTPPassword"},
		{cfg.SMTPPort, "SMTPPort"},
		{cfg.SMTPUser, "SMTPUser"},
		{cfg.ClientOrigin, "ClientOrigin"},
	}

	for _, field := range requiredFields {
		if field.value == "" {
			return fmt.Errorf("configuration variable '%s' must not be missing or empty", field.name)
		}
	}

	return nil
}

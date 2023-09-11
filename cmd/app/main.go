package main

import (
	"go-instagram-clone/config"
	"go-instagram-clone/internal/server"
	"go-instagram-clone/pkg/db"
	"go-instagram-clone/pkg/logger"
	"go-instagram-clone/pkg/utils"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	appLogger := logger.InitLogger(cfg)

	appLogger.Info("Starting api server")
	appLogger.Infof("App version: %s, Mode: %s", cfg.AppVersion, cfg.Mode)

	settings := db.Settings{
		Host:     cfg.PostgresHost,
		Port:     utils.ParsePort(cfg.PostgresPort),
		Database: cfg.PostgresDbname,
		User:     cfg.PostgresUser,
		Password: cfg.PostgresPassword,
		SSLMode:  cfg.PostgresSslMode,
	}
	// Create the database connection
	psDB, err := db.NewDatabaseConnection(cfg.PgDriver, settings)
	if err != nil {
		appLogger.Fatalf("Postgresql init: %s", err)
	} else {
		appLogger.Infof("Postgres connected, Status: %#v", psDB.Stats())
	}
	defer psDB.Close()

	s := server.New(cfg, *appLogger, psDB)
	if err = s.Run(); err != nil {
		appLogger.Fatal(err)
	}
}

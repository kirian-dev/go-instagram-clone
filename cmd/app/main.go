package main

import (
	"go-instagram-clone/config"
	"go-instagram-clone/internal/server"
	"go-instagram-clone/pkg/db/postgres"
	"go-instagram-clone/pkg/logger"
	"go-instagram-clone/pkg/utils"

	"log"
	"os"
)

func main() {
	log.Println("Starting api server")

	configPath := utils.GetConfigPath(os.Getenv("config"))

	cfgFile, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("LoadConfig: %v", err)
	}

	cfg, err := config.ParseConfig(cfgFile)
	if err != nil {
		log.Fatalf("ParseConfig: %v", err)
	}

	appLogger := logger.InitLogger(cfg)
	appLogger.Infof("App version: %s, Mode: %s", cfg.Server.AppVersion, cfg.Server.Mode)

	psDB, err := postgres.NewPostgresDB(cfg)
	if err != nil {
		appLogger.Fatalf("Postgres init: %s", err)
	} else {
		appLogger.Infof("Postgres connected, Status: %#v", psDB.Stats())
	}
	defer psDB.Close()

	s := server.New(cfg, *appLogger, psDB)
	if err = s.Run(); err != nil {
		log.Fatal(err)
	}
}

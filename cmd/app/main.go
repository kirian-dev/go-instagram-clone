package main

import (
	"go-instagram-clone/config"
	"go-instagram-clone/internal/server"
	"go-instagram-clone/pkg/db"
	"go-instagram-clone/pkg/logger"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	log := logger.InitLogger(cfg)

	log.Info("Starting api server")
	log.Infof("App version: %s, Mode: %s", cfg.AppVersion, cfg.Mode)

	psDB, err := db.NewDatabaseConnection(cfg)
	if err != nil {
		log.Fatalf("Postgresql init: %s", err)
	} else {
		log.Infof("Postgres connected, Status: %#v", psDB.Stats())
	}
	defer psDB.Close()

	s := server.New(cfg, log, psDB)
	if err = s.Run(); err != nil {
		log.Fatal(err)
	}
}

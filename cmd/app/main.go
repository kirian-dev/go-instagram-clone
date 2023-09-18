package main

import (
	"go-instagram-clone/config"
	"go-instagram-clone/internal/domain/models"
	"go-instagram-clone/internal/server"
	"go-instagram-clone/pkg/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "go-instagram-clone/docs/go-instagram-clone"
)

// @title GO-INSTAGRAM-CLONE
// @version 1.0
// @description This REST API for instagram clone.
// @contact.name Kirill Polozenko
// @contact.url https://github.com/kirian-dev
// @contact.email polozenko.kirill.job@gmail.com
// @BasePath /api/v1
// @host localhost:8080

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	log := logger.InitLogger(cfg)

	log.Info("Starting api server")
	log.Infof("App version: %s, Mode: %s", cfg.AppVersion, cfg.Mode)

	db, err := gorm.Open(postgres.Open(cfg.URI), &gorm.Config{})
	if err != nil {
		log.Error("error creating database, err: %v", err)
	}
	log.Info("connected to postgres database")

	db.AutoMigrate(&models.User{}, &models.Message{}, &models.Chat{})

	s := server.New(cfg, log, db)
	if err = s.Run(); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"fmt"
	"go-instagram-clone/config"
	"go-instagram-clone/pkg/logger"
	pb "go-instagram-clone/services/analytics/cmd/proto"
	"go-instagram-clone/services/chat/internal/domain/models"
	"go-instagram-clone/services/chat/internal/helpers"
	"go-instagram-clone/services/chat/internal/server"

	_ "go-instagram-clone/services/chat/docs/go-instagram-clone"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

	URI := fmt.Sprintf("host=%s port=%s dbname=%s sslmode=%s user=%s password=%s",
		cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresDbname, cfg.PostgresSslMode, cfg.PostgresUser, cfg.PostgresPassword)
	db, err := gorm.Open(postgres.Open(URI), &gorm.Config{})
	if err != nil {
		log.Error("error creating database, err: %v", err)
		return
	}
	log.Info("connected to postgres database")

	db.AutoMigrate(&models.User{}, &models.Message{}, &models.Chat{}, &models.ChatParticipant{}, &models.FileImport{})

	analyticsConn, err := grpc.Dial(cfg.AnalyticsPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("Connected to analytics service: %v", err)
		return
	}
	defer analyticsConn.Close()

	analyticsClient := pb.NewAnalyticsServiceClient(analyticsConn)
	log.Info("Connected to analytics service")

	// Create test file for upload accounts
	fileName := "test_accounts.csv"
	numAccounts := 20000

	if err := helpers.GenerateCSV(fileName, numAccounts, "public"); err != nil {
		fmt.Println("Error in create csv file:", err)
		return
	}
	log.Info("Created csv file for upload accounts ", fileName)

	s := server.New(cfg, log, db, analyticsClient)
	if err = s.Run(); err != nil {
		log.Fatal(err)
	}
}

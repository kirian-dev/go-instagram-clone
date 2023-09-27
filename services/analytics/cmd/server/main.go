package main

import (
	"fmt"
	"net"

	"go-instagram-clone/config"
	"go-instagram-clone/pkg/logger"
	deliveryGrpc "go-instagram-clone/services/analytics/internal/analytics/delivery/grpc"
	analyticsRepo "go-instagram-clone/services/analytics/internal/analytics/repository/storage/postgres"
	analyticsUC "go-instagram-clone/services/analytics/internal/analytics/useCase"
	"go-instagram-clone/services/analytics/internal/models"

	"google.golang.org/grpc"
	"gorm.io/driver/postgres"

	"gorm.io/gorm"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	log := logger.InitLogger(cfg)
	URI := fmt.Sprintf("host=%s port=%s dbname=%s sslmode=%s user=%s password=%s",
		cfg.AnalyticsPostgresHost, cfg.AnalyticsPostgresPort, cfg.AnalyticsPostgresDBName, cfg.AnalyticsPostgresSslMode, cfg.AnalyticsPostgresUser, cfg.AnalyticsPostgresPassword)
	db, err := gorm.Open(postgres.Open(URI), &gorm.Config{})
	log.Infof("Connecting to PostgreSQL using URI: %s", URI)
	if err != nil {
		log.Fatalf("Error creating postgres database connection: %v", err)
	}
	db.AutoMigrate(&models.Analytics{})

	lis, err := net.Listen("tcp", cfg.AnalyticsPort)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
		return
	}
	analyticsRepo := analyticsRepo.NewAnalyticsRepository(db)

	analyticsUC := analyticsUC.NewAnalyticsUC(cfg, log, analyticsRepo)

	grpcServer := grpc.NewServer()
	deliveryGrpc.NewAnalyticsServerGrpc(cfg, log, analyticsUC, grpcServer)

	log.Info("Server is running on..", cfg.AnalyticsPort)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
		return
	}
}

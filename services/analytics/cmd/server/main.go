package main

import (
	"fmt"
	"net"

	"go-instagram-clone/config"
	"go-instagram-clone/pkg/logger"
	deliveryGrpc "go-instagram-clone/services/analytics/internal/analytics/delivery/grpc"
	analyticsRepo "go-instagram-clone/services/analytics/internal/analytics/repository/storage/mysql"
	analyticsUC "go-instagram-clone/services/analytics/internal/analytics/useCase"
	"go-instagram-clone/services/analytics/internal/models"

	"google.golang.org/grpc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	log := logger.InitLogger(cfg)

	mysqlURI := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.MySQLUser, cfg.MySQLPassword, cfg.MySQLHost, cfg.MySQLPort, cfg.MySQLDBName)
	mysqlDB, err := gorm.Open(mysql.Open(mysqlURI), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error creating MySQL database connection: %v", err)
	}
	mysqlDB.AutoMigrate(&models.Analytics{})

	lis, err := net.Listen("tcp", cfg.AnalyticsPort)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
		return
	}
	analyticsRepo := analyticsRepo.NewAnalyticsRepository(mysqlDB)

	analyticsUC := analyticsUC.NewAnalyticsUC(cfg, log, analyticsRepo)

	grpcServer := grpc.NewServer()
	deliveryGrpc.NewAnalyticsServerGrpc(cfg, log, analyticsUC, grpcServer)

	log.Info("Server is running on %s...", cfg.AnalyticsPort)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
		return
	}
}

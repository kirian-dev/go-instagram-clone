package main

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"go-instagram-clone/config"
	"go-instagram-clone/pkg/logger"
	analytics "go-instagram-clone/services/analytics/cmd/proto"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	log := logger.InitLogger(cfg)
	conn, err := grpc.Dial(cfg.AnalyticsPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := analytics.NewAnalyticsServiceClient(conn)

	req := &analytics.LoginRequest{
		SuccessfulLogins: 1,
	}

	_, err = client.RecordLogin(context.Background(), req)
	if err != nil {
		log.Fatalf("Failed to record login: %v", err)
	}
}

package main

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"go-instagram-clone/config"
	"go-instagram-clone/pkg/logger"
	analytics "go-instagram-clone/services/analytics/cmd/proto"

	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log := logger.InitLogger(cfg)
	conn, err := grpc.Dial(cfg.AnalyticsPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := analytics.NewAnalyticsServiceClient(conn)
	ctx := context.Background()

	recordLogin(log, client, ctx)
	recordNewUser(log, client, ctx)
	getQuantityLogins(log, client, ctx)
	getQuantityRegistrations(log, client, ctx)
}

func handleError(log *logger.ZapLogger, message string, err error) {
	if err != nil {
		log.Fatalf("%s: %v", message, err)
	}
}

func recordLogin(log *logger.ZapLogger, client analytics.AnalyticsServiceClient, ctx context.Context) {
	loginReq := &analytics.LoginRequest{
		Email: "",
		Phone: "",
	}

	_, err := client.RecordLogin(ctx, loginReq)
	handleError(log, "Failed to record login", err)
	log.Info("Login recorded successfully.")
}

func recordNewUser(log *logger.ZapLogger, client analytics.AnalyticsServiceClient, ctx context.Context) {
	newUserReq := &analytics.NewUserRequest{
		Email: "",
		Phone: "",
	}

	_, err := client.RecordNewUser(ctx, newUserReq)
	handleError(log, "Failed to record new user", err)
	log.Info("New user registration recorded successfully.")
}

func getQuantityLogins(log *logger.ZapLogger, client analytics.AnalyticsServiceClient, ctx context.Context) {
	quantityLoginsReq := &emptypb.Empty{}

	quantityLoginsResp, err := client.GetQuantityLogins(ctx, quantityLoginsReq)
	handleError(log, "Failed to get quantity of logins", err)
	log.Info("Quantity of logins: %d\n", quantityLoginsResp.Quantity)
}

func getQuantityRegistrations(log *logger.ZapLogger, client analytics.AnalyticsServiceClient, ctx context.Context) {
	quantityRegisterReq := &emptypb.Empty{}

	quantityRegisterResp, err := client.GetQuantityRegister(ctx, quantityRegisterReq)
	handleError(log, "Failed to get quantity of registrations", err)
	log.Info("Quantity of registrations: %d\n", quantityRegisterResp.Quantity)
}

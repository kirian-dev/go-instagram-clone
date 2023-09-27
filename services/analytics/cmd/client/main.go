package main

import (
	"context"

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
		panic(err)
	}

	log := logger.InitLogger(cfg)
	conn, err := grpc.Dial(cfg.AnalyticsPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := analytics.NewAnalyticsServiceClient(conn)

	// Пример вызова метода RecordLogin
	loginReq := &analytics.LoginRequest{
		SuccessfulLogins: 1,
	}

	_, err = client.RecordLogin(context.Background(), loginReq)
	if err != nil {
		log.Fatalf("Failed to record login: %v", err)
	} else {
		log.Info("Login recorded successfully.")
	}

	// Пример вызова метода RecordNewUser
	newUserReq := &analytics.NewUserRequest{
		SuccessfulRegister: 1,
	}

	_, err = client.RecordNewUser(context.Background(), newUserReq)
	if err != nil {
		log.Fatalf("Failed to record new user: %v", err)
	} else {
		log.Info("New user registration recorded successfully.")
	}

	quantityLoginsReq := &emptypb.Empty{}

	quantityLoginsResp, err := client.GetQuantityLogins(context.Background(), quantityLoginsReq)
	if err != nil {
		log.Fatalf("Failed to get quantity of logins: %v", err)
	} else {
		log.Info("Quantity of logins: %d\n", quantityLoginsResp.Quantity)
	}

	quantityRegisterReq := &emptypb.Empty{}

	quantityRegisterResp, err := client.GetQuantityRegister(context.Background(), quantityRegisterReq)
	if err != nil {
		log.Fatalf("Failed to get quantity of registrations: %v", err)
	} else {
		log.Info("Quantity of registrations: %d\n", quantityRegisterResp.Quantity)
	}
}

package main

import (
	"context"
	"net"

	"go-instagram-clone/config"
	"go-instagram-clone/pkg/logger"
	analytics "go-instagram-clone/services/analytics/cmd/proto"

	"github.com/labstack/gommon/log"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type analyticsService struct {
	analytics.AnalyticsServiceServer
}

func (s *analyticsService) RecordLogin(ctx context.Context, req *analytics.LoginRequest) (*emptypb.Empty, error) {
	log.Info("Received RecordLogin request with SuccessfulLogins:", req.SuccessfulLogins)
	return &emptypb.Empty{}, nil
}
func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	log := logger.InitLogger(cfg)
	lis, err := net.Listen("tcp", cfg.AnalyticsPort)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	analytics.RegisterAnalyticsServiceServer(grpcServer, &analyticsService{})
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

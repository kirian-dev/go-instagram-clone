package delivery

import (
	"context"

	"google.golang.org/grpc/reflection"

	"go-instagram-clone/config"
	"go-instagram-clone/pkg/logger"
	analytics "go-instagram-clone/services/analytics/cmd/proto"
	"go-instagram-clone/services/analytics/internal/analytics/useCase"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type server struct {
	cfg         *config.Config
	log         *logger.ZapLogger
	analyticsUC useCase.AnalyticsUseCase
	analytics.UnimplementedAnalyticsServiceServer
}

func NewAnalyticsServerGrpc(cfg *config.Config, log *logger.ZapLogger, analyticsUC useCase.AnalyticsUseCase, gserver *grpc.Server) {
	analyticsServer := &server{
		analyticsUC: analyticsUC,
		cfg:         cfg,
		log:         log,
	}
	analytics.RegisterAnalyticsServiceServer(gserver, analyticsServer)
	reflection.Register(gserver)
}

func (h *server) RecordLogin(ctx context.Context, req *analytics.LoginRequest) (*emptypb.Empty, error) {
	err := h.analyticsUC.RecordLogin(req.Email, req.Phone)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (h *server) RecordNewUser(ctx context.Context, req *analytics.NewUserRequest) (*emptypb.Empty, error) {
	err := h.analyticsUC.RecordNewUser(req.Email, req.Phone)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (h *server) GetQuantityLogins(ctx context.Context, req *emptypb.Empty) (*analytics.QuantityResponse, error) {
	quantityLogins, err := h.analyticsUC.GetQuantityLogins()
	if err != nil {
		return nil, err
	}

	quantityResponse := &analytics.QuantityResponse{
		Quantity: quantityLogins,
	}

	return quantityResponse, nil
}

func (h *server) GetQuantityRegister(ctx context.Context, req *emptypb.Empty) (*analytics.QuantityResponse, error) {
	quantityRegister, err := h.analyticsUC.GetQuantityRegister()
	if err != nil {
		return nil, err
	}

	quantityResponse := &analytics.QuantityResponse{
		Quantity: quantityRegister,
	}

	return quantityResponse, nil
}

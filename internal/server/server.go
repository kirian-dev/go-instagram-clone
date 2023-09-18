package server

import (
	"context"
	"go-instagram-clone/config"

	"go-instagram-clone/pkg/logger"

	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

const (
	maxHeaderBytes = 1 << 20
	ctxTimeout     = 5
)

type Server struct {
	echo *echo.Echo
	cfg  *config.Config
	log  *logger.ZapLogger
	db   *gorm.DB
}

func New(cfg *config.Config, log *logger.ZapLogger, db *gorm.DB) *Server {
	return &Server{echo: echo.New(), cfg: cfg, log: log, db: db}
}

func (s *Server) Run() error {
	server := &http.Server{
		Addr:           s.cfg.Port,
		ReadTimeout:    time.Second * s.cfg.ReadTimeout,
		WriteTimeout:   time.Second * s.cfg.WriteTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}

	go func() {
		s.log.Infof("Server is listening on PORT: %s", s.cfg.Port)
		if err := s.echo.StartServer(server); err != nil {
			s.log.Fatalf("Error starting Server: ", err)
		}
	}()

	if err := s.Handlers(s.echo); err != nil {
		return err
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), ctxTimeout*time.Second)
	defer shutdown()

	s.log.Info("Server Exited Properly")
	return s.echo.Server.Shutdown(ctx)
}

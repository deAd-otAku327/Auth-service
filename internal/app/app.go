package app

import (
	"auth-service/internal/config"
	"auth-service/internal/controller"
	"auth-service/internal/repository"
	"auth-service/internal/service"
	"auth-service/internal/tokenizer"
	"auth-service/pkg/cryptor"
	"auth-service/pkg/logger"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

const AppName = "Auth-Service"

type App struct {
	Server *http.Server
}

func New(cfg *config.Config) (*App, error) {
	logger, err := logger.NewTextLogger(os.Stdout, cfg.LogLevel)
	if err != nil {
		return nil, err
	}

	repo, err := repository.NewPostgresRepo(cfg.DBConn, logger)
	if err != nil {
		return nil, err
	}

	cryptor := cryptor.New(cfg.AsyncHashingLimit)
	tokenizer := tokenizer.New(AppName, cfg.AccessTokenSecretKey, cfg.AccessTokenExpire, cfg.RefreshTokenExpire)

	service := service.New(repo, cryptor, tokenizer, logger)

	controller := controller.New(service, logger)

	return &App{
		Server: &http.Server{
			Addr:    fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
			Handler: initRoutes(controller),
		},
	}, nil
}

func (s *App) Run() error {
	slog.Info("app starting on", slog.String("address", s.Server.Addr))
	return s.Server.ListenAndServe()
}

func (s *App) Shutdown() error {
	slog.Info("app shutting down...")
	return s.Server.Shutdown(context.Background())
}

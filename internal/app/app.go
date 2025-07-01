package app

import (
	"auth-service/internal/config"
	"context"
	"fmt"
	"log/slog"
	"net/http"
)

const AppName = "Auth-Service"

type App struct {
	Server *http.Server
}

func New(cfg *config.Config) (*App, error) {
	return &App{
		Server: &http.Server{
			Addr:    fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
			Handler: nil,
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

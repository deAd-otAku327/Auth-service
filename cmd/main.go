package main

import (
	"auth-service/internal/app"
	"auth-service/internal/config"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

func main() {
	env := os.Getenv("ENV")
	if env == "" {
		log.Fatalln("error: missing app environment")
	}

	configDir := os.Getenv("CONFIG_DIR")
	configPath := filepath.Join(configDir, fmt.Sprintf("%s.yml", env))

	cfg, err := config.New(configPath)
	reportOnError(err)

	app, err := app.New(cfg)
	reportOnError(err)

	go func() {
		err = app.Run()
		if err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				log.Fatalf("app error: %v", err)
			}
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	<-ch
	if err := app.Shutdown(); err != nil {
		log.Fatalf("app shutdown failed: %v", err)
	}
	slog.Info("app shutdown completed")
}

func reportOnError(err error) {
	if err != nil {
		log.Fatalln(err.Error())
	}
}

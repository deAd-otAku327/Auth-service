package controller

import (
	"auth-service/internal/service"
	"log/slog"
	"net/http"
)

type Controller interface {
	HandleLogin() http.HandlerFunc
	HandleGetCurrentUser() http.HandlerFunc
	HandleRefresh() http.HandlerFunc
	HandleLogout() http.HandlerFunc
}

type authController struct {
	service service.Service
	logger  *slog.Logger
}

func New(service service.Service, logger *slog.Logger) Controller {
	return &authController{
		service: service,
		logger:  logger,
	}
}

func (c *authController) HandleLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func (c *authController) HandleGetCurrentUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func (c *authController) HandleRefresh() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func (c *authController) HandleLogout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

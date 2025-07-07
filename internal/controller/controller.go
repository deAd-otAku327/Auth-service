package controller

import (
	"auth-service/internal/apperrors"
	"auth-service/internal/controller/responser"
	"auth-service/internal/mappers/dtomap"
	"auth-service/internal/mappers/modelmap"
	"auth-service/internal/service"
	"auth-service/internal/types/dto"
	"log/slog"
	"net/http"

	"github.com/gorilla/schema"
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
		err := r.ParseForm()
		if err != nil {
			responser.MakeErrorResponseJSON(w, dtomap.MapToErrorResponse(apperrors.ErrInvalidRequestParams, http.StatusBadRequest))
			return
		}

		request := dto.LoginRequest{}
		err = schema.NewDecoder().Decode(&request, r.Form)
		if err != nil {
			responser.MakeErrorResponseJSON(w, dtomap.MapToErrorResponse(apperrors.ErrInvalidRequestParams, http.StatusBadRequest))
			return
		}

		clientIP := r.RemoteAddr
		originIPWithProxy := r.Header.Get("X-Forwarded-For")
		if originIPWithProxy != "" {
			clientIP = originIPWithProxy
		}

		response, refreshCookie, servErr := c.service.Login(r.Context(), modelmap.MapToLoginModel(&request, r.UserAgent(), clientIP))
		if servErr != nil {
			responser.MakeErrorResponseJSON(w, servErr)
			return
		}

		http.SetCookie(w, refreshCookie)
		responser.MakeResponseJSON(w, http.StatusOK, &response)
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

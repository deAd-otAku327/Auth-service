package controller

import (
	"auth-service/internal/controller/apierrors"
	"auth-service/internal/controller/responser"
	"auth-service/internal/mappers/dtomap"
	"auth-service/internal/mappers/modelmap"
	"auth-service/internal/middleware"
	"auth-service/internal/service"
	"auth-service/internal/tokenizer"
	"auth-service/internal/types/dto"
	"log/slog"
	"net/http"
	"strings"
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

const GUIDQueryParam = "guid"

func (c *authController) HandleLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			responser.MakeErrorResponseJSON(w, dtomap.MapToErrorResponse(apierrors.ErrInvalidRequestFormat, http.StatusBadRequest))
			return
		}

		request := dto.LoginRequest{
			UserGUID: r.URL.Query().Get(GUIDQueryParam),
		}

		clientIP := getClientIPFromRequest(r)

		response, refreshCookie, err := c.service.Login(r.Context(), modelmap.MapToLoginModel(&request, r.UserAgent(), clientIP))
		if err != nil {
			apierr := getAPIError(err)
			if apierr.Code == http.StatusInternalServerError {
				c.logger.Error(err.Error())
			}
			responser.MakeErrorResponseJSON(w, apierr)
			return
		}

		http.SetCookie(w, refreshCookie)
		responser.MakeResponseJSON(w, http.StatusOK, &response)
	}
}

func (c *authController) HandleGetCurrentUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response, err := c.service.GetCurrentUser(r.Context())
		if err != nil {
			apierr := getAPIError(err)
			if apierr.Code == http.StatusInternalServerError {
				c.logger.Error(err.Error())
			}
			responser.MakeErrorResponseJSON(w, apierr)
			return
		}

		responser.MakeResponseJSON(w, http.StatusOK, &response)
	}
}

func (c *authController) HandleRefresh() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accessToken := strings.TrimPrefix(r.Header.Get(middleware.AuthorizationHeader), "Bearer ")
		if accessToken == "" {
			responser.MakeErrorResponseJSON(w, dtomap.MapToErrorResponse(apierrors.ErrNoAccessToken, http.StatusBadRequest))
			return
		}
		refreshCookie, err := r.Cookie(tokenizer.RefreshCookieName)
		if err != nil {
			responser.MakeErrorResponseJSON(w, dtomap.MapToErrorResponse(apierrors.ErrNoRefreshToken, http.StatusBadRequest))
			return
		}

		clientIP := getClientIPFromRequest(r)

		response, refreshCookie, err := c.service.Refresh(r.Context(), modelmap.MapToRefreshModel(accessToken, refreshCookie.Value, r.UserAgent(), clientIP))
		if err != nil {
			apierr := getAPIError(err)
			if apierr.Code == http.StatusInternalServerError {
				c.logger.Error(err.Error())
			}
			responser.MakeErrorResponseJSON(w, apierr)
			return
		}

		http.SetCookie(w, refreshCookie)
		responser.MakeResponseJSON(w, http.StatusOK, &response)
	}
}

func (c *authController) HandleLogout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := c.service.Logout(r.Context(), r.UserAgent())
		if err != nil {
			apierr := getAPIError(err)
			if apierr.Code == http.StatusInternalServerError {
				c.logger.Error(err.Error())
			}
			responser.MakeErrorResponseJSON(w, apierr)
			return
		}

		responser.MakeResponseJSON(w, http.StatusOK, nil)
	}
}

package middleware

import (
	"auth-service/internal/controller/apierrors"
	"auth-service/internal/controller/responser"
	"auth-service/internal/mappers/dtomap"
	"auth-service/internal/tokenizer"
	"context"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

const AuthorizationHeader = "Authorization"

const (
	UserGUIDKey contextKey = iota
)

type contextKey int8

type Middleware interface {
	Auth() mux.MiddlewareFunc
}

type middlewareHub struct {
	tokenizer tokenizer.Tokenizer
}

func New(tok tokenizer.Tokenizer) Middleware {
	return &middlewareHub{
		tokenizer: tok,
	}
}

func (hub *middlewareHub) Auth() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			accessToken := strings.TrimPrefix(r.Header.Get(AuthorizationHeader), "Bearer ")

			claims, err := hub.tokenizer.VerifyAccessTokenJWT(accessToken, false)
			if err != nil {
				responser.MakeErrorResponseJSON(w, dtomap.MapToErrorResponse(apierrors.ErrAuthenticationFailed, http.StatusForbidden))
				return
			}

			userGUID, err := claims.GetSubject()
			if err != nil {
				responser.MakeErrorResponseJSON(w, dtomap.MapToErrorResponse(apierrors.ErrAuthenticationFailed, http.StatusForbidden))
				return
			}

			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), UserGUIDKey, userGUID)))
		})
	}
}

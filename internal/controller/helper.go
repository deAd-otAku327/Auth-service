package controller

import (
	"auth-service/internal/controller/apierrors"
	"auth-service/internal/mappers/dtomap"
	"auth-service/internal/service/serverrors"
	"auth-service/internal/types/dto"
	"errors"
	"net/http"
	"strings"
)

const XRealIPHeader = "X-Real-IP"

func getAPIError(err error) *dto.ErrorResponse {
	if errors.Is(err, serverrors.ErrSessionAlreadyExists) {
		return dtomap.MapToErrorResponse(apierrors.ErrAuthorizeNotNeeded, http.StatusBadRequest)
	} else if errors.Is(err, serverrors.ErrUserGUIDInvalid) {
		return dtomap.MapToErrorResponse(apierrors.ErrInvalidRequestData, http.StatusBadRequest)
	} else if errors.Is(err, serverrors.ErrNoRefreshSession) {
		return dtomap.MapToErrorResponse(apierrors.ErrRefreshUnavalible, http.StatusUnauthorized)
	} else if errors.Is(err, serverrors.ErrOldAccessTokenInvalid) ||
		errors.Is(err, serverrors.ErrRefreshTokenInvalid) ||
		errors.Is(err, serverrors.ErrTokenPairInvalid) {

		return dtomap.MapToErrorResponse(apierrors.ErrAuthenticationFailed, http.StatusForbidden)
	} else {
		return dtomap.MapToErrorResponse(apierrors.ErrSomethingWentWrong, http.StatusInternalServerError)
	}
}

func getClientIPFromRequest(r *http.Request) string {
	splittedRemoteAddress := strings.Split(r.RemoteAddr, ":")

	clientIP := ""
	if len(splittedRemoteAddress) == 2 {
		clientIP = splittedRemoteAddress[0]
	}

	xrealIP := r.Header.Get(XRealIPHeader)
	if xrealIP != "" {
		clientIP = xrealIP
	}

	return clientIP
}

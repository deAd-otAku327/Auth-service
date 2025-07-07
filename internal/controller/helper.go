package controller

import (
	"auth-service/internal/controller/apierrors"
	"auth-service/internal/mappers/dtomap"
	"auth-service/internal/service/serverrors"
	"auth-service/internal/types/dto"
	"errors"
	"net/http"
)

func getAPIError(err error) *dto.ErrorResponse {
	if errors.Is(err, serverrors.ErrSessionAlreadyExists) {
		return dtomap.MapToErrorResponse(apierrors.ErrAuthorizeNotNeeded, http.StatusFound)
	} else if errors.Is(err, serverrors.ErrUserGUIDInvalid) {
		return dtomap.MapToErrorResponse(apierrors.ErrInvalidRequestData, http.StatusBadRequest)
	} else {
		return dtomap.MapToErrorResponse(apierrors.ErrSomethingWentWrong, http.StatusInternalServerError)
	}
}

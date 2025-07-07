package responser

import (
	"auth-service/internal/types/dto"
	"encoding/json"
	"net/http"
)

const (
	contentTypeHeader = "Content-Type"
	contentTypeJSON   = "application/json"
)

func MakeResponseJSON(w http.ResponseWriter, code int, data any) {
	w.Header().Set(contentTypeHeader, contentTypeJSON)
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func MakeErrorResponseJSON(w http.ResponseWriter, apierr *dto.ErrorResponse) {
	MakeResponseJSON(w, apierr.Code, apierr)
}

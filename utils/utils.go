package utils

import (
	"encoding/json"
	"net/http"
)

type jwtClaimKey string

const UserContextKey jwtClaimKey = "USER"

type ErrorResponse struct {
	Status       int    `json:"status"`
	Message      string `json:"message"`
	RedirectPath string `json:"redirect_path"`
}

func ConstructErrorResponse(w http.ResponseWriter, e ErrorResponse) {
	w.WriteHeader(e.Status)
	res, _ := json.Marshal(e)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(res)
}

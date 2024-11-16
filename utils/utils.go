package utils

import (
	"encoding/json"
	"net/http"
)

var TwoDaysInSeconds = 172800

type jwtClaimKey string

const UserContextKey jwtClaimKey = "USER"

type ErrorResponse struct {
	Status       int    `json:"status"`
	Message      string `json:"message"`
	RedirectPath string `json:"redirect_path"`
}

type CookieParam struct {
	CookieName   string
	CookieDomain string
	CookieSecure bool
}

func ConstructErrorResponse(w http.ResponseWriter, e ErrorResponse) {
	w.WriteHeader(e.Status)
	res, _ := json.Marshal(e)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(res)
}

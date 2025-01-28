package utils

import (
	"encoding/json"
	"net/http"
)

var TwoDaysInSeconds = 172800

type jwtClaimKey string

const UserContextKey jwtClaimKey = "USER"

type Error struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type CookieParam struct {
	CookieName   string
	CookieDomain string
	CookieSecure bool
	SameSite     http.SameSite
}

func ErrorResponse(w http.ResponseWriter, err error) {
	obj := Error{
		Status:  errorStatus(err),
		Message: err.Error(),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(obj.Status)
	res, _ := json.Marshal(obj)
	_, _ = w.Write(res)
}

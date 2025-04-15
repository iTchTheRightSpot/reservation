package utils

import (
	"net/http"
	"time"
)

var TimeFormat = time.RFC3339

var TwoDaysInSeconds = 172800

type jwtClaimKey string

const UserContextKey jwtClaimKey = "USER"

type CookieParam struct {
	CookieName   string
	CookieDomain string
	CookieSecure bool
	SameSite     http.SameSite
}
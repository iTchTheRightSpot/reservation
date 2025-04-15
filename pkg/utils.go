package pkg

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/iTchTheRightSpot/reservation/utils"
	log "github.com/iTchTheRightSpot/utility/utils"
	"io"
	"net/http"
	"time"
)

type Db interface {
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

func WriteCookie(w http.ResponseWriter, param *utils.CookieParam, token string, expireAt time.Time) {
	http.SetCookie(
		w,
		&http.Cookie{
			Name:     param.CookieName,
			Value:    token,
			Expires:  expireAt,
			HttpOnly: true,
			SameSite: param.SameSite,
			Secure:   param.CookieSecure,
			Domain:   param.CookieDomain,
			Path:     "/",
		},
	)
}

func ReadBody[T any](lg log.ILogger, r *http.Request) (*T, error) {
	if r.Body == nil {
		return nil, &log.ServerError{Message: "request body cannot be nil"}
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		lg.Error(r.Context(), err.Error())
		return nil, &log.ServerError{Message: "failed to read request body"}
	}

	var payload T
	if err = json.Unmarshal(bodyBytes, &payload); err != nil {
		lg.Error(r.Context(), err.Error())
		return nil, &log.ServerError{Message: "failed to unmarshal request body"}
	}

	return &payload, nil
}
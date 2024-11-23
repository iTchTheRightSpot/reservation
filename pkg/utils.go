package pkg

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/iTchTheRightSpot/erp-golang/utils"
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
			SameSite: http.SameSiteLaxMode,
			Secure:   param.CookieSecure,
			Domain:   param.CookieDomain,
			Path:     "/",
		},
	)
}

func ReadBody[T any](r *http.Request) (*T, error) {
	if r.Body == nil {
		return nil, fmt.Errorf("request body cannot be nil")
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read request body: %w", err)
	}

	var payload T
	if err := json.Unmarshal(bodyBytes, &payload); err != nil {
		return nil, fmt.Errorf("failed to unmarshal request body: %w", err)
	}

	return &payload, nil
}

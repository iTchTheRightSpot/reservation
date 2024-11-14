package middleware

import (
	"context"
	"fmt"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models"
	"github.com/iTchTheRightSpot/erp-golang/pkg/services/auth"
	"github.com/iTchTheRightSpot/erp-golang/utils"
	"net/http"
	"reflect"
	"slices"
	"strings"
	"time"
)

type Middleware[T any] struct {
	Logger     utils.ILogger
	Auth       auth.IJwtService
	CookieName string
}

type middlewareFunc func(http.Handler) http.Handler

func (dep *Middleware[T]) Initialize(router *http.ServeMux) http.Handler {
	stack := dep.createStack(dep.timeout, dep.logging)
	return stack(router)
}

func (dep *Middleware[T]) createStack(m ...middlewareFunc) middlewareFunc {
	return func(next http.Handler) http.Handler {
		for i := len(m) - 1; i >= 0; i-- {
			next = m[i](next)
		}
		return next
	}
}

// https://stackoverflow.com/questions/27234861/correct-way-of-getting-clients-ip-addresses-from-http-request
func (dep *Middleware[T]) requestIP(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip != "" {
		// the header can contain multiple IPs, so take the first one
		return strings.Split(ip, ",")[0]
	}
	// fallback to RemoteAddr
	return r.RemoteAddr
}

func (dep *Middleware[T]) timeout(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dur := time.Second * time.Duration(5)
		ctx, cancel := context.WithTimeout(r.Context(), dur)
		defer cancel()

		// create a new request with the timeout context
		r = r.WithContext(ctx)

		// call the next handler
		next.ServeHTTP(w, r)
	})
}

func (dep *Middleware[T]) logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := dep.requestIP(r)
		dep.Logger.Log(fmt.Sprintf("[Request] IP: %s | Method: %s | Path: %s", ip, r.Method, r.URL.Path))

		type wrappedWriter struct {
			http.ResponseWriter
			statusCode int
		}

		wrapped := &wrappedWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(wrapped, r)
		dep.Logger.Log(fmt.Sprintf("[Request] IP: %s | [Response] Status: %v | Method: %s | Path: %s", ip, wrapped.statusCode, r.Method, r.URL.Path))
	})
}

func (dep *Middleware[T]) Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(dep.CookieName)
		if err != nil {
			dep.Logger.Error(err)
			utils.ConstructErrorResponse(
				w,
				utils.ErrorResponse{
					Status:  http.StatusUnauthorized,
					Message: fmt.Sprintf("unauthorized cookie not present"),
				},
			)
			return
		}

		u, err := dep.Auth.ValidateJwt(cookie.Value)
		if err != nil {
			dep.Logger.Error(err)
			utils.ConstructErrorResponse(
				w,
				utils.ErrorResponse{
					Status:  http.StatusUnauthorized,
					Message: fmt.Sprintf("%s", err),
				},
			)
			return
		}

		// TODO refresh token if about to expire

		// add the user to the context
		ctx := r.Context()
		ctx = context.WithValue(ctx, utils.UserContextKey, u)
		r = r.WithContext(ctx)

		// call the function if the token is valid
		next.ServeHTTP(w, r)
	})
}

func (dep *Middleware[T]) Permission(next http.Handler, roles ...models.RoleEnum) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		obj := r.Context().Value(utils.UserContextKey).(*models.StaffJwtObj)

		var contains bool
		for i := 0; i < len(roles); i++ {
			contains = slices.ContainsFunc(obj.Roles, func(role models.RoleEnum) bool {
				return reflect.DeepEqual(role, roles[i])
			})
			if contains {
				break
			}
		}

		if !contains {
			dep.Logger.Error("access denied")
			utils.ConstructErrorResponse(
				w,
				utils.ErrorResponse{
					Status:  http.StatusForbidden,
					Message: "access denied",
				},
			)
			return
		}

		next.ServeHTTP(w, r)
	})
}

package middleware

import (
	"context"
	"fmt"
	"github.com/iTchTheRightSpot/erp-golang/pkg"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models"
	"github.com/iTchTheRightSpot/erp-golang/pkg/services/auth"
	"github.com/iTchTheRightSpot/erp-golang/utils"
	"net/http"
	"reflect"
	"slices"
	"strings"
	"time"
)

type wrappedWriter struct {
	http.ResponseWriter
	status int
}

func (w *wrappedWriter) WriteHeader(statusCode int) {
	w.status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

type Middleware struct {
	Logger utils.ILogger
	Auth   auth.IJwtService
	Param  *utils.CookieParam
}

type middlewareFunc func(http.Handler) http.Handler

func (dep *Middleware) Initialize(router *http.ServeMux) http.Handler {
	stack := dep.createStack(dep.timeout, dep.logging)
	return stack(router)
}

func (dep *Middleware) createStack(m ...middlewareFunc) middlewareFunc {
	return func(next http.Handler) http.Handler {
		for i := len(m) - 1; i >= 0; i-- {
			next = m[i](next)
		}
		return next
	}
}

// https://stackoverflow.com/questions/27234861/correct-way-of-getting-clients-ip-addresses-from-http-request
func (dep *Middleware) requestIP(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip != "" {
		// the header can contain multiple IPs, so take the first one
		return strings.Split(ip, ",")[0]
	}
	// fallback to RemoteAddr
	return r.RemoteAddr
}

func (dep *Middleware) timeout(next http.Handler) http.Handler {
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

func (dep *Middleware) logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := dep.requestIP(r)
		dep.Logger.Log(fmt.Sprintf(
			"[Request] IP: %s | Method: %s | Path: %s", ip, r.Method, r.URL.Path,
		))
		obj := &wrappedWriter{
			ResponseWriter: w,
			status:         http.StatusOK,
		}
		next.ServeHTTP(obj, r)
		dep.Logger.Log(fmt.Sprintf(
			"[Response] IP: %s | Status: %d | Method: %s | Path: %s",
			ip, obj.status, r.Method, r.URL.Path,
		))
	})
}

func (dep *Middleware) Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allowed, err := pkg.CasbinEnforcer.Enforce(models.STAFF, r.URL.Path, r.Method)

		if err != nil {
			dep.Logger.Error(fmt.Sprintf("CASBIN ENFORCER ERROR %s", err))
		} else {
			dep.Logger.Log(fmt.Sprintf("CASBIN ENFORCER LOG %v", allowed))
		}

		if r.Cookies() == nil || len(r.Cookies()) == 0 {
			dep.Logger.Error("no cookie present")
			utils.ConstructErrorResponse(
				w,
				utils.ErrorResponse{
					Status:  http.StatusUnauthorized,
					Message: fmt.Sprintf("unauthorized cookie not present"),
				},
			)
			return
		}

		cookie, err := r.Cookie(dep.Param.CookieName)
		if err != nil || cookie == nil {
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

		obj, err := dep.Auth.ValidateJwt(cookie.Value)
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

		// validate if token is about to expire
		isTokenExpiringSoon := func(now time.Time, t time.Time, expirationInSeconds int) bool {
			return t.Before(now.Add(time.Duration(expirationInSeconds) * time.Second))
		}

		b := !strings.HasSuffix(r.URL.Path, "/logout")
		if b && isTokenExpiringSoon(dep.Logger.Date(), obj.ExpireAt, utils.TwoDaysInSeconds) {
			if o, err := dep.Auth.GenerateJwt(obj, utils.TwoDaysInSeconds); err != nil {
				dep.Logger.Error(fmt.Sprintf("failed to refresh token %s", err))
			} else {
				cookie.Value = o.Token
				cookie.Expires = o.ExpireAt
				http.SetCookie(w, cookie)
				dep.Logger.Log("Refreshed jwt")
			}
		}

		// add the user to the context
		ctx := r.Context()
		ctx = context.WithValue(ctx, utils.UserContextKey, obj)
		r = r.WithContext(ctx)

		// call the function if the token is valid
		next.ServeHTTP(w, r)
	})
}

func (dep *Middleware) Permission(next http.Handler, roles ...models.RoleEnum) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		obj := r.Context().Value(utils.UserContextKey).(*models.JwtObj)

		contains := false
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

func (dep *Middleware) ChainAuth(next http.Handler, roles ...models.RoleEnum) http.Handler {
	return dep.Authentication(dep.Permission(next, roles...))
}

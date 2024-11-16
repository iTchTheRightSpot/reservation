package middleware

import (
	"context"
	"fmt"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models"
	"github.com/iTchTheRightSpot/erp-golang/pkg/services/auth"
	"github.com/iTchTheRightSpot/erp-golang/utils"
	"net/http"
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
	return r.RemoteAddr
}

func (dep *Middleware) timeout(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dur := time.Second * time.Duration(5)
		ctx, cancel := context.WithTimeout(r.Context(), dur)
		defer cancel()
		next.ServeHTTP(w, r.WithContext(ctx))
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
		if b && isTokenExpiringSoon(dep.Logger.Date(), *obj.ExpireAt, utils.TwoDaysInSeconds) {
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

func (dep *Middleware) HasRole(next http.Handler, role *models.RoleEnum) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if role == nil {
			dep.Logger.Error("HasRole: role cannot be nil")
			utils.ConstructErrorResponse(
				w,
				utils.ErrorResponse{
					Status:  http.StatusInternalServerError,
					Message: "internal server error",
				},
			)
			return
		}

		obj, ok := r.Context().Value(utils.UserContextKey).(*models.JwtObj)
		if !ok || obj == nil {
			dep.Logger.Error("HasRoleAndPermissions: invalid user context")
			utils.ConstructErrorResponse(w, utils.ErrorResponse{
				Status:  http.StatusUnauthorized,
				Message: "unauthorized",
			})
			return
		}

		contains := slices.ContainsFunc(obj.AccessControls, func(access models.RolePermission) bool {
			return access.Role == *role
		})
		if !contains {
			dep.Logger.Error(fmt.Sprintf("access denied request role does not match %v", role))
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

func (dep *Middleware) HasRoleAndPermissions(next http.Handler, cred *models.RolePermission) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if cred == nil {
			dep.Logger.Error("HasRoleAndPermissions: cred cannot be nil")
			utils.ConstructErrorResponse(w, utils.ErrorResponse{
				Status:  http.StatusInternalServerError,
				Message: "internal server error",
			})
			return
		}

		obj, ok := r.Context().Value(utils.UserContextKey).(*models.JwtObj)
		if !ok || obj == nil {
			dep.Logger.Error("HasRoleAndPermissions: invalid user context")
			utils.ConstructErrorResponse(w, utils.ErrorResponse{
				Status:  http.StatusUnauthorized,
				Message: "unauthorized",
			})
			return
		}

		if !dep.validateRoleAndPermissions(obj.AccessControls, cred) {
			dep.Logger.Error("HasRoleAndPermissions: insufficient roles or permissions")
			utils.ConstructErrorResponse(w, utils.ErrorResponse{
				Status:  http.StatusForbidden,
				Message: "access denied",
			})
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (dep *Middleware) validateRoleAndPermissions(arr []models.RolePermission, param *models.RolePermission) bool {
	return slices.ContainsFunc(arr, func(obj models.RolePermission) bool {
		if obj.Role == param.Role {
			if len(param.Permissions) < 1 {
				return false
			}
			for _, permission := range param.Permissions {
				if !slices.Contains(obj.Permissions, permission) {
					return false
				}
			}
			return true
		}
		return false
	})
}

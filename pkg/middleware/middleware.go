package middleware

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
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

func (dep *Middleware) Initialize(router *http.ServeMux) http.Handler {
	return dep.logging(dep.timeout(router))
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
		start := dep.Logger.Date()
		ip := dep.requestIP(r)
		id := uuid.NewString()
		dep.Logger.Log(fmt.Sprintf(
			"[Request] ID: %s | IP: %s | Method: %s | Path: %s", id, ip, r.Method, r.URL.Path,
		))
		obj := &wrappedWriter{
			ResponseWriter: w,
			status:         http.StatusOK,
		}
		next.ServeHTTP(obj, r)
		end := dep.Logger.Date()
		dep.Logger.Log(fmt.Sprintf(
			"[Response] ID: %s | IP: %s | Status: %d | Method: %s | Path: %s | Duration: %v seconds",
			id, ip, obj.status, r.Method, r.URL.Path, end.Sub(start).Seconds(),
		))
	})
}

func (dep *Middleware) Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Cookies() == nil || len(r.Cookies()) == 0 {
			dep.Logger.Error("no cookie present")
			utils.ErrorResponse(w, &utils.AuthenticationError{})
			return
		}

		cookie, err := r.Cookie(dep.Param.CookieName)
		if err != nil || cookie == nil {
			dep.Logger.Error(err)
			utils.ErrorResponse(w, &utils.AuthenticationError{})
			return
		}

		obj, err := dep.Auth.Decode(cookie.Value)
		if err != nil {
			dep.Logger.Error(err)
			utils.ErrorResponse(w, &utils.AuthenticationError{})
			return
		}

		isTokenExpiringSoon := func(now time.Time, t time.Time, expirationInSeconds int) bool {
			return t.Before(now.Add(time.Duration(expirationInSeconds) * time.Second))
		}

		isLogout := strings.HasSuffix(r.URL.Path, "/logout")
		if !isLogout && isTokenExpiringSoon(dep.Logger.Date(), *obj.ExpireAt, utils.TwoDaysInSeconds) {
			if o, err := dep.Auth.Encode(obj, utils.TwoDaysInSeconds); err != nil {
				dep.Logger.Error(fmt.Sprintf("failed to refresh token %s", err))
			} else {
				cookie.Value = o.Token
				cookie.Expires = o.ExpireAt
				http.SetCookie(w, cookie)
				dep.Logger.Log("refreshed jwt")
			}
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), utils.UserContextKey, obj)))
	})
}

func (dep *Middleware) HasRole(next http.Handler, role *models.RoleEnum) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if role == nil {
			dep.Logger.Error("HasRole: role cannot be nil")
			utils.ErrorResponse(w, errors.New("internal server error"))
			return
		}

		obj, ok := r.Context().Value(utils.UserContextKey).(*models.JwtObj)
		if !ok || obj == nil {
			dep.Logger.Error("HasRoleAndPermissions: invalid user context")
			utils.ErrorResponse(w, &utils.AuthenticationError{})
			return
		}

		contains := slices.ContainsFunc(obj.AccessControls, func(access models.RolePermissionEnum) bool {
			return access.Role == *role
		})
		if !contains {
			dep.Logger.Error(fmt.Sprintf("access denied request role does not match %v", role))
			utils.ErrorResponse(w, &utils.AccessDeniedError{})
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (dep *Middleware) HasPermission(next http.Handler, permission *models.PermissionEnum) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if permission == nil {
			dep.Logger.Error("HasPermission: permission cannot be nil")
			utils.ErrorResponse(w, errors.New("internal server error"))
			return
		}

		obj, ok := r.Context().Value(utils.UserContextKey).(*models.JwtObj)
		if !ok || obj == nil {
			dep.Logger.Error("HasRoleAndPermissions: invalid user context")
			utils.ErrorResponse(w, &utils.AuthenticationError{})
			return
		}

		contains := slices.ContainsFunc(obj.AccessControls, func(rp models.RolePermissionEnum) bool {
			return slices.ContainsFunc(rp.Permissions, func(enum models.PermissionEnum) bool {
				return reflect.DeepEqual(enum, *permission)
			})
		})

		if !contains {
			dep.Logger.Error(fmt.Sprintf("access denied request permission does not match %v", permission))
			utils.ErrorResponse(w, &utils.AccessDeniedError{})
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (dep *Middleware) HasRoleAndPermissions(next http.Handler, cred *models.RolePermissionEnum) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if cred == nil {
			dep.Logger.Error("HasRoleAndPermissions: cred cannot be nil")
			utils.ErrorResponse(w, errors.New("internal server error"))
			return
		}

		obj, ok := r.Context().Value(utils.UserContextKey).(*models.JwtObj)
		if !ok || obj == nil {
			dep.Logger.Error("HasRoleAndPermissions: invalid user context")
			utils.ErrorResponse(w, &utils.AuthenticationError{})
			return
		}

		if !dep.validateRoleAndPermissions(obj.AccessControls, cred) {
			dep.Logger.Error("HasRoleAndPermissions: insufficient roles or permissions")
			utils.ErrorResponse(w, &utils.AccessDeniedError{})
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (dep *Middleware) validateRoleAndPermissions(arr []models.RolePermissionEnum, param *models.RolePermissionEnum) bool {
	return slices.ContainsFunc(arr, func(obj models.RolePermissionEnum) bool {
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

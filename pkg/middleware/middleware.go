package middleware

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/iTchTheRightSpot/reservation/pkg/models"
	"github.com/iTchTheRightSpot/reservation/pkg/services/auth"
	"github.com/iTchTheRightSpot/reservation/utils"
	"github.com/iTchTheRightSpot/utility/middleware"
	log "github.com/iTchTheRightSpot/utility/utils"
	"net/http"
	"reflect"
	"slices"
	"strings"
	"time"
)

type Middleware struct {
	Logger    log.ILogger
	Auth      auth.IJwtService
	Param     *utils.CookieParam
	ApiPrefix string
	FS        http.FileSystem
	Validator *validator.Validate
}

func (dep *Middleware) Initialize(mux *http.ServeMux) http.Handler {
	m := middleware.Middleware{Logger: dep.Logger, Fs: dep.FS, ApiPrefix: dep.ApiPrefix}
	return m.Log(m.SPA(mux))
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

func (dep *Middleware) Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Cookies() == nil || len(r.Cookies()) == 0 {
			dep.Logger.Error(r.Context(), "no cookie present")
			log.ErrorResponse(w, &log.AuthenticationError{})
			return
		}

		cookie, err := r.Cookie(dep.Param.CookieName)
		if err != nil || cookie == nil {
			dep.Logger.Error(r.Context(), err)
			log.ErrorResponse(w, &log.AuthenticationError{})
			return
		}

		obj, err := dep.Auth.Decode(r.Context(), cookie.Value)
		if err != nil {
			dep.Logger.Error(r.Context(), err.Error())
			log.ErrorResponse(w, &log.AuthenticationError{})
			return
		}

		isTokenExpiringSoon := func(now time.Time, t time.Time, expirationInSeconds int) bool {
			return t.Before(now.Add(time.Duration(expirationInSeconds) * time.Second))
		}

		isLogout := strings.HasSuffix(r.URL.Path, "/logout")
		if !isLogout && isTokenExpiringSoon(dep.Logger.Date(), *obj.ExpireAt, utils.TwoDaysInSeconds) {
			if o, err := dep.Auth.Encode(r.Context(), obj, utils.TwoDaysInSeconds); err != nil {
				dep.Logger.Error(r.Context(), "failed to refresh token", err.Error())
			} else {
				cookie.Domain = dep.Param.CookieDomain
				cookie.Value = o.Token
				cookie.Expires = o.ExpireAt
				cookie.Path = "/"
				cookie.HttpOnly = true
				cookie.SameSite = dep.Param.SameSite
				cookie.Secure = dep.Param.CookieSecure
				http.SetCookie(w, cookie)
				dep.Logger.Log(r.Context(), "refreshed jwt")
			}
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), utils.UserContextKey, obj)))
	})
}

func (dep *Middleware) HasRole(next http.Handler, role *models.RoleEnum) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if role == nil {
			dep.Logger.Error(r.Context(), "HasRole: role cannot be nil")
			log.ErrorResponse(w, &log.ServerError{})
			return
		}

		obj, ok := r.Context().Value(utils.UserContextKey).(*models.JwtObj)
		if !ok || obj == nil {
			dep.Logger.Error(r.Context(), "HasRoleAndPermissions: invalid user context")
			log.ErrorResponse(w, &log.AuthenticationError{})
			return
		}

		contains := slices.ContainsFunc(obj.AccessControls, func(access models.RolePermissionEnum) bool {
			return access.Role == *role
		})
		if !contains {
			dep.Logger.Error(r.Context(), fmt.Sprintf("access denied request role does not match %v", role))
			log.ErrorResponse(w, &log.AccessDeniedError{})
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (dep *Middleware) HasPermission(next http.Handler, permission *models.PermissionEnum) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if permission == nil {
			dep.Logger.Error(r.Context(), "HasPermission: permission cannot be nil")
			log.ErrorResponse(w, errors.New("internal server error"))
			return
		}

		obj, ok := r.Context().Value(utils.UserContextKey).(*models.JwtObj)
		if !ok || obj == nil {
			dep.Logger.Error(r.Context(), "HasRoleAndPermissions: invalid user context")
			log.ErrorResponse(w, &log.AuthenticationError{})
			return
		}

		contains := slices.ContainsFunc(obj.AccessControls, func(rp models.RolePermissionEnum) bool {
			return slices.ContainsFunc(rp.Permissions, func(enum models.PermissionEnum) bool {
				return reflect.DeepEqual(enum, *permission)
			})
		})

		if !contains {
			dep.Logger.Error(r.Context(), fmt.Sprintf("access denied request permission does not match %v", permission))
			log.ErrorResponse(w, &log.AccessDeniedError{})
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (dep *Middleware) HasRoleAndPermissions(next http.Handler, cred *models.RolePermissionEnum) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if cred == nil {
			dep.Logger.Error(r.Context(), "HasRoleAndPermissions: cred cannot be nil")
			log.ErrorResponse(w, &log.ServerError{})
			return
		}

		obj, ok := r.Context().Value(utils.UserContextKey).(*models.JwtObj)
		if !ok || obj == nil {
			dep.Logger.Error(r.Context(), "HasRoleAndPermissions: invalid user context")
			log.ErrorResponse(w, &log.AuthenticationError{})
			return
		}

		if !dep.validateRoleAndPermissions(obj.AccessControls, cred) {
			dep.Logger.Error(r.Context(), "HasRoleAndPermissions: insufficient roles or permissions")
			log.ErrorResponse(w, &log.AccessDeniedError{})
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
package middleware

import (
	"context"
	"fmt"
	"github.com/iTchTheRightSpot/erp-golang/config"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models"
	"github.com/iTchTheRightSpot/erp-golang/pkg/services/auth"
	"github.com/iTchTheRightSpot/erp-golang/utils"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

func TestMiddleware(t *testing.T) {
	if err := os.Chdir("../../"); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}

	t.Parallel()

	sv := &config.SecretVariables{}
	env, err := sv.Config()
	if err != nil {
		t.Fatalf("%s", err)
	}

	logger := utils.NewMockLogger()

	t.Run("Authentication middleware", func(t *testing.T) {
		t.Parallel()

		t.Run("reject request cookie is nil", func(t *testing.T) {
			t.Parallel()

			mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			req := httptest.NewRequest(http.MethodGet, "/path", nil)
			resp := httptest.NewRecorder()

			// method to test
			jwtService := &auth.MockJwtService{}
			middleware := Middleware{Logger: logger, Auth: jwtService}
			middleware.Authentication(mockHandler).ServeHTTP(resp, req)

			if resp.Code != http.StatusUnauthorized {
				t.Errorf("expected status code %d, got %d", http.StatusUnauthorized, resp.Code)
			}

			if jwtService.ValidateJwtCalled {
				t.Errorf("expected ValidateJwtCalled to be false but is true")
			}
		})

		t.Run("reject request cookie not present", func(t *testing.T) {
			t.Parallel()

			mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			req := httptest.NewRequest(http.MethodGet, "/path", nil)
			req.AddCookie(&http.Cookie{Name: "NOT_NAME", Value: "jwt"})
			resp := httptest.NewRecorder()

			// method to test
			jwtService := &auth.MockJwtService{}
			middleware := Middleware{Logger: logger, Auth: jwtService, Param: env.CookieParam}
			middleware.Authentication(mockHandler).ServeHTTP(resp, req)

			if resp.Code != http.StatusUnauthorized {
				t.Errorf("expected status code %d, got %d", http.StatusUnauthorized, resp.Code)
			}

			if jwtService.ValidateJwtCalled {
				t.Errorf("expected ValidateJwtCalled to be false but is true")
			}
		})

		t.Run("reject request invalid jwt", func(t *testing.T) {
			t.Parallel()

			mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			req := httptest.NewRequest(http.MethodGet, "/path", nil)
			req.AddCookie(&http.Cookie{Name: env.CookieParam.CookieName, Value: "jwt"})
			resp := httptest.NewRecorder()

			// method to test
			jwtService := &auth.MockJwtService{ValidateJwtError: fmt.Errorf("")}
			middleware := Middleware{Logger: logger, Auth: jwtService, Param: env.CookieParam}
			middleware.Authentication(mockHandler).ServeHTTP(resp, req)

			if resp.Code != http.StatusUnauthorized {
				t.Errorf("expected status code %d, got %d", http.StatusUnauthorized, resp.Code)
			}

			if !jwtService.ValidateJwtCalled {
				t.Errorf("expected ValidateJwtCalled to be true but is false")
			}
		})

		t.Run("accept request valid jwt", func(t *testing.T) {
			t.Parallel()

			mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			req := httptest.NewRequest(http.MethodGet, "/path", nil)
			req.AddCookie(&http.Cookie{Name: env.CookieParam.CookieName, Value: "jwt"})
			resp := httptest.NewRecorder()

			// method to test
			add := time.Now().Add(180 * time.Hour)
			obj := &models.JwtObj{ExpireAt: &add}
			jwtService := &auth.MockJwtService{StaffJwtObj: obj}
			middleware := Middleware{Logger: logger, Auth: jwtService, Param: env.CookieParam}
			middleware.Authentication(mockHandler).ServeHTTP(resp, req)

			if resp.Code != http.StatusOK {
				t.Errorf("expected status code %d, got %d", http.StatusOK, resp.Code)
			}

			if !jwtService.ValidateJwtCalled {
				t.Errorf("expected ValidateJwtCalled to be true but is false")
			}
		})

		t.Run("refresh token is called", func(t *testing.T) {
			t.Parallel()

			mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			req := httptest.NewRequest(http.MethodGet, "/path", nil)
			req.AddCookie(&http.Cookie{Name: env.CookieParam.CookieName, Value: "jwt"})
			resp := httptest.NewRecorder()

			// method to test
			add := time.Now().Add(8 * time.Hour)
			obj := &models.JwtObj{ExpireAt: &add}
			jwtRes := &models.JwtResponse{Token: "new-token", ExpireAt: time.Now().Add(80 * time.Hour)}
			jwtService := &auth.MockJwtService{StaffJwtObj: obj, JwtResponse: jwtRes}
			middleware := Middleware{Logger: logger, Auth: jwtService, Param: env.CookieParam}
			middleware.Authentication(mockHandler).ServeHTTP(resp, req)

			if resp.Code != http.StatusOK {
				t.Errorf("expected status code %d, got %d", http.StatusOK, resp.Code)
			}

			if !jwtService.ValidateJwtCalled {
				t.Errorf("expected ValidateJwtCalled to be true but is false")
			}

			if !jwtService.GenerateJwtCalled {
				t.Errorf("expected GenerateJwtCalled to be true but is false")
			}

			newToken := resp.Result().Header.Get("Set-Cookie")
			if len(newToken) == 0 {
				t.Errorf("header is empty")
			}

			if !strings.Contains(newToken, "new-token") {
				t.Errorf("expected new-token to be in cookie, isn't %s", newToken)
			}
		})

		t.Run("should not call refresh token as route is logout", func(t *testing.T) {
			t.Parallel()

			mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			req := httptest.NewRequest(http.MethodGet, "/path/logout", nil)
			req.AddCookie(&http.Cookie{Name: env.CookieParam.CookieName, Value: "jwt"})
			resp := httptest.NewRecorder()

			// method to test
			add := time.Now().Add(8 * time.Hour)
			obj := &models.JwtObj{ExpireAt: &add}
			jwtRes := &models.JwtResponse{Token: "new-token", ExpireAt: time.Now().Add(80 * time.Hour)}
			jwtService := &auth.MockJwtService{StaffJwtObj: obj, JwtResponse: jwtRes}
			middleware := Middleware{Logger: logger, Auth: jwtService, Param: env.CookieParam}
			middleware.Authentication(mockHandler).ServeHTTP(resp, req)

			if resp.Code != http.StatusOK {
				t.Errorf("expected status code %d, got %d", http.StatusOK, resp.Code)
			}

			if !jwtService.ValidateJwtCalled {
				t.Errorf("expected ValidateJwtCalled to be true but is false")
			}

			if jwtService.GenerateJwtCalled {
				t.Errorf("expected GenerateJwtCalled to be false but is true")
			}
		})
	})

	t.Run("HasRole middleware", func(t *testing.T) {
		t.Parallel()

		t.Run(fmt.Sprintf("reject request. not of role %s", models.DEVELOPER), func(t *testing.T) {
			t.Parallel()

			// given
			cred := make([]models.RolePermission, 2)
			cred[0] = models.RolePermission{
				Role:        models.STAFF,
				Permissions: []models.PermissionEnum{models.READ, models.DELETE},
			}
			obj := &models.JwtObj{
				AccessControls: cred,
				UserUUID:       "staff-id",
			}

			mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			req := httptest.NewRequest(http.MethodGet, "/path", nil)
			ctx := req.Context()
			ctx = context.WithValue(ctx, utils.UserContextKey, obj)
			req = req.WithContext(ctx)
			resp := httptest.NewRecorder()

			middleware := Middleware{Logger: logger}

			// method to test
			role := models.DEVELOPER
			middleware.HasRole(mockHandler, &role).ServeHTTP(resp, req)

			if resp.Code != http.StatusForbidden {
				t.Errorf("expected status code %d, got %d", http.StatusForbidden, resp.Code)
			}
		})

		t.Run(fmt.Sprintf("accept request. role matches %s", models.STAFF), func(t *testing.T) {
			t.Parallel()

			// given
			cred := make([]models.RolePermission, 2)
			cred[0] = models.RolePermission{
				Role:        models.STAFF,
				Permissions: []models.PermissionEnum{models.READ, models.DELETE, models.WRITE},
			}
			obj := &models.JwtObj{
				AccessControls: cred,
				UserUUID:       "staff-id",
			}

			mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			req := httptest.NewRequest(http.MethodGet, "/path", nil)
			ctx := req.Context()
			ctx = context.WithValue(ctx, utils.UserContextKey, obj)
			req = req.WithContext(ctx)
			resp := httptest.NewRecorder()

			middleware := Middleware{Logger: logger}

			// method to test
			staff := models.STAFF
			middleware.HasRole(mockHandler, &staff).ServeHTTP(resp, req)

			if resp.Code != http.StatusOK {
				t.Errorf("expected status code %d, got %d", http.StatusOK, resp.Code)
			}
		})
	})

	t.Run("HasRoleAndPermissions middleware", func(t *testing.T) {
		t.Parallel()
	})

}

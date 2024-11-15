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
			obj := &models.JwtObj{ExpireAt: time.Now().Add(180 * time.Hour)}
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
			obj := &models.JwtObj{ExpireAt: time.Now().Add(8 * time.Hour)}
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
			obj := &models.JwtObj{ExpireAt: time.Now().Add(8 * time.Hour)}
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

	t.Run("Permission middleware", func(t *testing.T) {
		t.Parallel()

		t.Run("reject request access denied. need role Developer", func(t *testing.T) {
			t.Parallel()

			// given
			obj := &models.JwtObj{
				Roles:    []models.RoleEnum{models.USER, models.STAFF},
				UserUUID: "staff-id",
			}

			mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			req := httptest.NewRequest(http.MethodGet, "/path", nil)
			ctx := req.Context()
			ctx = context.WithValue(ctx, utils.UserContextKey, obj)
			req = req.WithContext(ctx)
			resp := httptest.NewRecorder()

			// method to test
			middleware := Middleware{Logger: logger}
			middleware.Permission(mockHandler, models.DEVELOPER).ServeHTTP(resp, req)

			if resp.Code != http.StatusForbidden {
				t.Errorf("expected status code %d, got %d", http.StatusForbidden, resp.Code)
			}
		})

		t.Run("reject request Permission roles is empty in ", func(t *testing.T) {
			t.Parallel()

			// given
			obj := &models.JwtObj{
				Roles:    []models.RoleEnum{models.USER, models.STAFF, models.DEVELOPER},
				UserUUID: "staff-id",
			}

			mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			req := httptest.NewRequest(http.MethodGet, "/path", nil)
			ctx := req.Context()
			ctx = context.WithValue(ctx, utils.UserContextKey, obj)
			req = req.WithContext(ctx)
			resp := httptest.NewRecorder()

			// method to test
			middleware := Middleware{Logger: logger}
			middleware.Permission(mockHandler).ServeHTTP(resp, req)

			if resp.Code != http.StatusForbidden {
				t.Errorf("expected status code %d, got %d", http.StatusForbidden, resp.Code)
			}
		})

		t.Run("accept request roles exist", func(t *testing.T) {
			t.Parallel()

			// given
			obj := &models.JwtObj{
				Roles:    []models.RoleEnum{models.USER, models.STAFF, models.DEVELOPER},
				UserUUID: "staff-id",
			}

			mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			req := httptest.NewRequest(http.MethodGet, "/path", nil)
			ctx := req.Context()
			ctx = context.WithValue(ctx, utils.UserContextKey, obj)
			req = req.WithContext(ctx)
			resp := httptest.NewRecorder()

			// method to test
			middleware := Middleware{Logger: logger}
			middleware.Permission(mockHandler, models.DEVELOPER).ServeHTTP(resp, req)

			if resp.Code != http.StatusOK {
				t.Errorf("expected status code %d, got %d", http.StatusOK, resp.Code)
			}
		})
	})

	t.Run("ChainAuth middleware", func(t *testing.T) {
		t.Parallel()

		t.Run("should reject request invalid jwt", func(t *testing.T) {
			t.Parallel()

			mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			req := httptest.NewRequest(http.MethodGet, "/chain", nil)
			req.AddCookie(&http.Cookie{Name: env.CookieParam.CookieName, Value: "jwt"})
			resp := httptest.NewRecorder()

			// method to test
			jwtService := &auth.MockJwtService{ValidateJwtError: fmt.Errorf("")}
			middleware := Middleware{Logger: logger, Auth: jwtService, Param: env.CookieParam}
			middleware.ChainAuth(mockHandler).ServeHTTP(resp, req)

			if resp.Code != http.StatusUnauthorized {
				t.Errorf("expected status code %d, got %d", http.StatusUnauthorized, resp.Code)
			}

			if !jwtService.ValidateJwtCalled {
				t.Errorf("expected ValidateJwtCalled to be true but is false")
			}
		})

		t.Run("should reject request access denied although valid jwt", func(t *testing.T) {
			t.Parallel()

			mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			jwtService := auth.NewJwtService(logger, env)

			o := &models.JwtObj{UserUUID: "uuid", Roles: []models.RoleEnum{models.USER, models.DEVELOPER}}
			obj, err := jwtService.GenerateJwt(o, utils.TwoDaysInSeconds)
			if err != nil {
				t.Errorf("%s", err)
			}

			req := httptest.NewRequest(http.MethodGet, "/chain", nil)
			req.AddCookie(&http.Cookie{Name: env.CookieParam.CookieName, Value: obj.Token})
			resp := httptest.NewRecorder()

			// method to test
			middleware := Middleware{Logger: logger, Auth: jwtService, Param: env.CookieParam}
			middleware.ChainAuth(mockHandler, models.STAFF).ServeHTTP(resp, req)

			if resp.Code != http.StatusForbidden {
				t.Errorf("expected status code %d, got %d", http.StatusForbidden, resp.Code)
			}
		})

		t.Run("should accept request valid jwt & matching roles", func(t *testing.T) {
			t.Parallel()

			mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			jwtService := auth.NewJwtService(logger, env)

			o := &models.JwtObj{UserUUID: "uuid", Roles: []models.RoleEnum{models.USER, models.DEVELOPER}}
			obj, err := jwtService.GenerateJwt(o, utils.TwoDaysInSeconds)
			if err != nil {
				t.Errorf("%s", err)
			}

			req := httptest.NewRequest(http.MethodGet, "/chain", nil)
			req.AddCookie(&http.Cookie{Name: env.CookieParam.CookieName, Value: obj.Token})
			resp := httptest.NewRecorder()

			// method to test
			middleware := Middleware{Logger: logger, Auth: jwtService, Param: env.CookieParam}
			middleware.ChainAuth(mockHandler, models.USER, models.DEVELOPER, models.STAFF).ServeHTTP(resp, req)

			if resp.Code != http.StatusOK {
				t.Errorf("expected status code %d, got %d", http.StatusOK, resp.Code)
			}
		})
	})
}

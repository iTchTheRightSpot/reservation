package middleware

import (
	"bytes"
	"encoding/json"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models"
	"github.com/iTchTheRightSpot/erp-golang/utils"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestBodyMiddleware(t *testing.T) {
	t.Parallel()

	logger := utils.NewDevLogger()

	t.Run("should accept request valid request body", func(t *testing.T) {
		t.Parallel()

		// given
		type obj struct {
			Id string `json:"id" validate:"required"`
		}

		dto, err := json.Marshal(obj{Id: "staff-id"})
		if err != nil {
			t.Errorf("failed to marshal obj: %s", err)
		}

		mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/path", bytes.NewBuffer(dto))
		rr := httptest.NewRecorder()

		// method to test
		b := RequestBodyMiddleware[obj]{Logger: logger}

		b.RequestBody(mockHandler).ServeHTTP(rr, req)

		// assert
		if rr.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
		}
	})

	t.Run("should reject request invalid request body", func(t *testing.T) {
		t.Parallel()

		// given
		type nest struct {
			Role models.RoleEnum `json:"role" validate:"required"`
		}
		type obj struct {
			Id   string  `json:"id" validate:"required"`
			Nest *[]nest `json:"nest" validate:"required,dive,required"`
		}

		dto, err := json.Marshal(obj{Id: "staff-id"})
		if err != nil {
			t.Errorf("failed to marshal obj: %s", err)
		}

		mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/path", bytes.NewBuffer(dto))
		rr := httptest.NewRecorder()

		// method to test
		b := RequestBodyMiddleware[obj]{Logger: logger}

		b.RequestBody(mockHandler).ServeHTTP(rr, req)

		// assert
		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("should reject request empty roles", func(t *testing.T) {
		t.Parallel()

		// given
		type obj struct {
			Id    string            `json:"id" validate:"required"`
			Roles []models.RoleEnum `json:"roles" validate:"required,min=1,dive,required"`
		}

		dto, err := json.Marshal(obj{Id: "staff-id", Roles: []models.RoleEnum{}})
		if err != nil {
			t.Errorf("failed to marshal obj: %s", err)
		}

		mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/path", bytes.NewBuffer(dto))
		rr := httptest.NewRecorder()

		// method to test
		b := RequestBodyMiddleware[obj]{Logger: logger}

		b.RequestBody(mockHandler).ServeHTTP(rr, req)

		// assert
		if rr.Code == http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})
}

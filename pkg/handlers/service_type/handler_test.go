package service_type

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/iTchTheRightSpot/erp-golang/config"
	"github.com/iTchTheRightSpot/erp-golang/database"
	"github.com/iTchTheRightSpot/erp-golang/pkg/handlers"
	"github.com/iTchTheRightSpot/erp-golang/pkg/middleware"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/profile"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/service"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/staff"
	"github.com/iTchTheRightSpot/erp-golang/pkg/services/auth"
	"github.com/iTchTheRightSpot/erp-golang/pkg/services/service_type"
	"github.com/iTchTheRightSpot/erp-golang/pkg/stores"
	"github.com/iTchTheRightSpot/erp-golang/utils"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

var db *sql.DB
var env *config.SecretVariables

func TestMain(m *testing.M) {
	if err := os.Chdir("../../../"); err != nil {
		log.Fatalf("failed to change directory: %v", err)
	}

	secret := config.SecretVariables{}
	env = secret.Config()

	d, err := database.ConnectToPostgres(env.DbConnectionString)
	if err != nil {
		log.Fatal(err)
	}
	db = d

	defer func(db *sql.DB) {
		if err := db.Close(); err != nil {
			log.Printf("db connection did not close after tests")
			return
		}
	}(db)

	m.Run()
}

func TestServiceHandler(t *testing.T) {
	t.Cleanup(func() {
		if err := handlers.DeleteAll(db); err != nil {
			t.Errorf(err.Error())
		}
	})

	// given
	logger := utils.NewMockLogger()
	mux := http.NewServeMux()
	prov := stores.NewTransactionProvider(logger, db)
	adapters := stores.NewAdapters(logger, db, prov)
	jwtSer := auth.NewJwtService(logger, env)
	m := &middleware.Middleware{Logger: logger, Auth: jwtSer, Param: env.CookieParam}
	s := service_type.NewServiceImpl(logger, adapters)

	// register all routes
	NewServiceHandler(mux, logger, s, m).Register()

	t.Run("should create service", func(t *testing.T) {
		arr := []models.RolePermission{
			{
				Role:        models.STAFF,
				Permissions: []models.PermissionEnum{models.WRITE},
			},
		}

		obj, err := jwtSer.GenerateJwt(
			&models.JwtObj{
				UserId:         "staff-uuid",
				AccessControls: arr,
			},
			utils.TwoDaysInSeconds,
		)

		dtoBytes, err := json.Marshal(service.ServiceTypePayload{
			Name:        uuid.New().String(),
			Price:       15.97,
			Duration:    3600,
			CleanUpTime: 1800,
			IsVisible:   true,
		})
		if err != nil {
			t.Errorf("failed to marshal SchedulePayload: %s", err)
		}

		req := httptest.NewRequest(http.MethodPost, "/service", bytes.NewBuffer(dtoBytes))
		req.AddCookie(&http.Cookie{Name: env.CookieParam.CookieName, Value: obj.Token})
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, req)

		// assert
		if rr.Code != http.StatusCreated {
			t.Errorf("expected status code %d, got %d", http.StatusCreated, rr.Code)
		}
	})

	t.Run("should retrieve all services", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/service", nil)
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)

		// assert
		if rr.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
		}

		type ui struct {
			Name     string  `json:"name"`
			Price    float64 `json:"price"`
			Duration int     `json:"duration"`
		}

		var resp []*ui
		if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
			t.Errorf(err.Error())
		}

		// assert
		size := len(resp)
		if size != 1 {
			t.Errorf("expected size 1, got %d", size)
			t.Errorf("%v", resp)
		}

		for _, obj := range resp {
			if len(obj.Name) < 1 {
				t.Error("expect name to not be empty")
			}

			if obj.Price < 1 {
				t.Error("expect Price to be greater than zero")
			}

			if obj.Duration < 1 {
				t.Error("expect Duration to be greater than zero")
			}
		}
	})

	t.Run("should retrieve staffs by services", func(t *testing.T) {
		t.Run("reject missing request param(s)", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/service/staffs", nil)
			rr := httptest.NewRecorder()
			mux.ServeHTTP(rr, req)

			// assert
			if rr.Code != http.StatusBadRequest {
				t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
			}

			con := strings.Contains(rr.Body.String(), "bad request, missing services type(s)")
			if !con {
				t.Errorf("expect to contain bad request, missing services type(s)")
				t.Errorf("%s", rr.Body.String())
			}
		})

		t.Run("success", func(t *testing.T) {
			// pre-save
			ctx := context.Background()
			serv1 := preSave(t, adapters, ctx)
			serv2 := preSave(t, adapters, ctx)

			url := fmt.Sprintf("/service/staffs?service=%s&service=%s", serv1.Name, serv2.Name)
			req := httptest.NewRequest(http.MethodGet, url, nil)
			rr := httptest.NewRecorder()
			mux.ServeHTTP(rr, req)

			// assert
			if rr.Code != http.StatusOK {
				t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
				t.Errorf("%s", rr.Body.String())
			}

			type ui struct {
				StaffId  string  `json:"staff_id"`
				Name     string  `json:"name"`
				ImageKey *string `json:"image_key"`
				Bio      *string `json:"bio"`
			}

			var resp []*ui
			if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
				t.Errorf(err.Error())
			}

			// assert
			size := len(resp)
			if size != 2 {
				t.Errorf("expected size 1, got %d", size)
				t.Errorf("%v", resp)
			}

			for _, obj := range resp {
				if len(obj.StaffId) < 1 {
					t.Error("expect staffId to not be empty")
				}

				if len(obj.Name) < 1 {
					t.Error("expect name to not be empty")
				}

				if obj.Bio == nil || len(*obj.Bio) < 1 {
					t.Error("expect bio to not be nil or empty")
					t.Error(obj.Bio)
				}
			}
		})
	})
}

func preSave(t *testing.T, adapters *stores.Adapters, ctx context.Context) service.ServiceTypeEntity {
	prof := profile.Profile{Firstname: "f", Lastname: "l", Email: uuid.NewString()}
	if err := adapters.ProfileStore.Save(ctx, &prof); err != nil {
		t.Error(err.Error())
	}

	staf := staff.Staff{ProfileId: &prof.ProfileId, UUID: uuid.New()}
	if err := adapters.StaffStore.Save(ctx, &staf); err != nil {
		t.Error(err.Error())
	}

	serv := service.ServiceTypeEntity{
		Name:        uuid.NewString(),
		Price:       20,
		IsVisible:   true,
		Duration:    3600,
		CleanUpTime: 30 * 60,
	}
	if err := adapters.ServiceStore.Save(ctx, &serv); err != nil {
		t.Error(err.Error())
	}

	if err := adapters.StaffServiceStore.Save(ctx, &staff.StaffServiceEntity{
		StaffId: staf.StaffId, ServiceId: serv.ServiceId,
	}); err != nil {
		t.Error(err.Error())
	}
	return serv
}

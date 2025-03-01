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
	serviceModel "github.com/iTchTheRightSpot/erp-golang/pkg/models/service_type"
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

func TestServiceTypeHandler(t *testing.T) {
	t.Cleanup(func() {
		if err := handlers.DeleteAll(db); err != nil {
			t.Error(err.Error())
		}
	})

	// given
	logger := utils.NewDevLogger()
	mux := http.NewServeMux()
	prov := stores.NewTransactionProvider(logger, db)
	adapters := stores.NewAdapters(logger, db, prov)
	jwtSer := auth.NewJwtServiceAsymmetric(logger, env)
	m := &middleware.Middleware{Logger: logger, Auth: jwtSer, Param: env.CookieParam}
	s := service_type.NewServiceImpl(logger, adapters)

	// register all routes
	NewServiceTypeHandler(mux, logger, s, m).Register()

	t.Run("should retrieve all services. empty", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/service", nil)
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)

		// assert
		if rr.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
		}

		if strings.TrimSpace(rr.Body.String()) != "[]" {
			t.Errorf("expect [], given %s", rr.Body.String())
		}
	})

	t.Run("crm", func(t *testing.T) {
		serviceName := uuid.NewString()
		save := preSave(t, adapters, context.Background())
		cred := []models.RolePermissionEnum{
			{
				Role:        models.STAFF,
				Permissions: []models.PermissionEnum{models.WRITE, models.READ},
			},
		}
		jwtObj, _ := jwtSer.Encode(
			&models.JwtObj{
				UserId:         "staff-uuid",
				AccessControls: cred,
			},
			utils.TwoDaysInSeconds,
		)

		t.Run("should create service", func(t *testing.T) {
			dtoBytes, err := json.Marshal(serviceModel.ServiceTypePayload{
				Name:        serviceName,
				Price:       15.97,
				Duration:    3600,
				CleanUpTime: 1800,
				IsVisible:   false,
			})
			if err != nil {
				t.Errorf("failed to marshal SchedulePayload: %s", err)
			}

			req := httptest.NewRequest(http.MethodPost, "/crm/service", bytes.NewBuffer(dtoBytes))
			req.AddCookie(&http.Cookie{Name: env.CookieParam.CookieName, Value: jwtObj.Token})
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()

			mux.ServeHTTP(rr, req)

			// assert
			if rr.Code != http.StatusCreated {
				t.Errorf("expected status code %d, got %d", http.StatusCreated, rr.Code)
			}
		})

		t.Run("should retrieve all services", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/crm/services", nil)
			req.AddCookie(&http.Cookie{Name: env.CookieParam.CookieName, Value: jwtObj.Token})
			rr := httptest.NewRecorder()
			mux.ServeHTTP(rr, req)

			// assert
			if rr.Code != http.StatusOK {
				t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
			}

			type ui struct {
				Name      string  `json:"name"`
				Price     float64 `json:"price"`
				IsVisible bool    `json:"is_visible"`
				Duration  int     `json:"duration"`
				CU        int     `json:"clean_up_time"`
			}

			var resp []*ui
			if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
				t.Error(err.Error())
			}

			// assert
			size := len(resp)
			if size != 2 {
				t.Errorf("expected size 2, got %d", size)
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

				if obj.CU < 1 {
					t.Error("expect Clean up time to be greater than zero", obj.CU)
				}
			}
		})

		t.Run("should link service to staff & also reject", func(t *testing.T) {
			obj, err := jwtSer.Encode(
				&models.JwtObj{
					UserId:         save.staff.UUID.String(),
					AccessControls: cred,
				},
				utils.TwoDaysInSeconds,
			)
			if err != nil {
				t.Error(err.Error())
			}

			pl, err := json.Marshal(serviceModel.LinkServiceTypeToStaffPayload{Service: serviceName, StaffUUID: save.staff.UUID.String()})
			if err != nil {
				t.Errorf("failed to marshal SchedulePayload: %s", err)
			}

			for i := 0; i < 2; i++ {
				req := httptest.NewRequest(http.MethodPost, "/crm/service/staff", bytes.NewBuffer(pl))
				req.AddCookie(&http.Cookie{Name: env.CookieParam.CookieName, Value: obj.Token})
				req.Header.Set("Content-Type", "application/json")
				rr := httptest.NewRecorder()

				mux.ServeHTTP(rr, req)

				// assert
				if i == 0 {
					if rr.Code != http.StatusCreated {
						t.Errorf("index %v expected status code %d, got %d", i, http.StatusCreated, rr.Code)
					}
				} else {
					if rr.Code != http.StatusConflict {
						t.Errorf("index %v expected status code %d, got %d", i, http.StatusBadRequest, rr.Code)
					}
				}
			}
		})

		t.Run("should update service type", func(t *testing.T) {
			dtoBytes, err := json.Marshal(serviceModel.ServiceTypePayload{
				Name:        serviceName,
				Price:       10.97,
				Duration:    3600,
				CleanUpTime: 1800,
				IsVisible:   false,
			})
			if err != nil {
				t.Errorf("failed to marshal SchedulePayload: %s", err)
			}

			req := httptest.NewRequest(http.MethodPut, "/crm/service", bytes.NewBuffer(dtoBytes))
			req.AddCookie(&http.Cookie{Name: env.CookieParam.CookieName, Value: jwtObj.Token})
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()

			mux.ServeHTTP(rr, req)

			// assert
			if rr.Code != http.StatusNoContent {
				t.Errorf("expected status code %d, got %d", http.StatusNoContent, rr.Code)
			}
		})

		t.Run("should retrieve service types by staff uuid", func(t *testing.T) {
			url := fmt.Sprintf("/crm/services/staff?staff_id=%s", save.staff.UUID.String())
			req := httptest.NewRequest(http.MethodGet, url, nil)
			req.AddCookie(&http.Cookie{Name: env.CookieParam.CookieName, Value: jwtObj.Token})
			rr := httptest.NewRecorder()
			mux.ServeHTTP(rr, req)

			// assert
			if rr.Code != http.StatusOK {
				t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
				t.Log(rr.Body.String())
			}

			var arr []string
			if err := json.Unmarshal(rr.Body.Bytes(), &arr); err != nil {
				t.Error(err.Error())
			}

			if len(arr) != 2 {
				t.Errorf("expect size to be 2, given %v", len(arr))
				t.Log(arr)
			}
		})
	})

	t.Run("store front", func(t *testing.T) {
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

			var resp []ui
			if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
				t.Error(err.Error())
			}

			// assert
			size := len(resp)
			if size != 1 {
				t.Errorf("expected size 1, got %d", size)
				t.Errorf("%v", resp)
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

			// pre-save
			ctx := context.Background()
			serv1 := preSave(t, adapters, ctx)
			serv2 := preSave(t, adapters, ctx)

			type ui struct {
				StaffId  string  `json:"staff_id"`
				Name     string  `json:"name"`
				ImageKey *string `json:"image_key"`
				Bio      *string `json:"bio"`
			}

			t.Run("success. should be empty as all services are not for each staff", func(t *testing.T) {
				url := fmt.Sprintf("/service/staffs?name=%s&name=%s", serv1.service.Name, serv2.service.Name)
				req := httptest.NewRequest(http.MethodGet, url, nil)
				rr := httptest.NewRecorder()
				mux.ServeHTTP(rr, req)

				// assert
				if rr.Code != http.StatusOK {
					t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
					t.Errorf("%s", rr.Body.String())
				}

				var resp []*ui
				if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
					t.Error(err.Error())
				}

				// assert
				size := len(resp)
				if size >= 1 {
					t.Errorf("expected size 0, got %d", size)
					t.Errorf("%v", resp)
				}
			})

			t.Run("success", func(t *testing.T) {
				url := fmt.Sprintf("/service/staffs?name=%s", serv1.service.Name)
				req := httptest.NewRequest(http.MethodGet, url, nil)
				rr := httptest.NewRecorder()
				mux.ServeHTTP(rr, req)

				// assert
				if rr.Code != http.StatusOK {
					t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
					t.Errorf("%s", rr.Body.String())
				}

				var resp []*ui
				if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
					t.Error(err.Error())
				}

				// assert
				size := len(resp)
				if size != 1 {
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
	})
}

type testObj struct {
	service *serviceModel.ServiceTypeEntity
	profile *models.ProfileEntity
	staff   *staff.StaffEntity
}

func preSave(t *testing.T, adapters *stores.Adapters, ctx context.Context) *testObj {
	prof := models.ProfileEntity{Firstname: "f", Lastname: "l", Email: uuid.NewString(), Password: "password"}
	if err := adapters.ProfileStore.Save(ctx, &prof); err != nil {
		t.Error(err.Error())
	}

	staf := staff.StaffEntity{ProfileId: &prof.ProfileId, UUID: uuid.New()}
	if err := adapters.StaffStore.Save(ctx, &staf); err != nil {
		t.Error(err.Error())
	}

	serv := serviceModel.ServiceTypeEntity{
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

	return &testObj{profile: &prof, staff: &staf, service: &serv}
}

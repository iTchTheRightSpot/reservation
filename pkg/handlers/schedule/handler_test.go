package schedule

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/iTchTheRightSpot/erp-golang/config"
	"github.com/iTchTheRightSpot/erp-golang/database"
	"github.com/iTchTheRightSpot/erp-golang/pkg/middleware"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/profile"
	model "github.com/iTchTheRightSpot/erp-golang/pkg/models/schedule"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/staff"
	"github.com/iTchTheRightSpot/erp-golang/pkg/services/auth"
	"github.com/iTchTheRightSpot/erp-golang/pkg/services/schedule"
	"github.com/iTchTheRightSpot/erp-golang/pkg/stores"
	"github.com/iTchTheRightSpot/erp-golang/utils"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

var db *sql.DB
var env *config.SecretVariables

func TestMain(m *testing.M) {
	if err := os.Chdir("../../../"); err != nil {
		log.Fatalf("failed to change directory: %v", err)
	}

	secret := config.SecretVariables{}
	e, err := secret.Config()
	if err != nil {
		log.Fatal(err)
	}
	env = e

	db, err = database.ConnectToPostgres(e.DbConnectionString)
	if err != nil {
		log.Fatal(err)
	}

	defer func(db *sql.DB) {
		if err := db.Close(); err != nil {
			log.Printf("db connection did not close after tests")
			return
		}
	}(db)

	m.Run()
}

func setupTest(t *testing.T) (*sql.Tx, func()) {
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("failed to start transaction: %v", err)
	}

	return tx, func() {
		if err := tx.Rollback(); err != nil {
			return
		}
	}
}

func preSaveStaff(a *stores.Adapters) (*staff.Staff, error) {
	ctx := context.Background()

	p := profile.Profile{
		Firstname: "erp",
		Lastname:  "erp",
		Email:     "erp@email.com",
	}

	if _, err := a.ProfileStore.Save(ctx, &p); err != nil {
		return nil, err
	}

	s := staff.Staff{
		UUID:      uuid.New(),
		ProfileId: &p.ProfileId,
	}

	if _, err := a.StaffStore.Save(ctx, &s); err != nil {
		return nil, err
	}

	return &s, nil
}

func TestScheduleHandler(t *testing.T) {
	t.Parallel()

	logger := utils.NewMockLogger()

	t.Run("save & retrieve schedules handlers", func(t *testing.T) {
		t.Parallel()

		t.Run("reject request duplicate schedule", func(t *testing.T) {
			t.Parallel()

			tx, fn := setupTest(t)
			defer fn()

			// given
			mux := http.NewServeMux()
			prov := stores.MockLiveTransactionProvider(logger, tx)
			adapters := stores.NewAdapters(logger, tx, prov)
			jwtSer := auth.NewJwtService(logger, env)
			m := &middleware.Middleware{Logger: logger, Auth: jwtSer, Param: env.CookieParam}
			s := schedule.NewScheduleService(logger, adapters)

			save, err := preSaveStaff(adapters)
			if err != nil {
				t.Error(err)
			}

			cred := []models.RolePermission{
				{
					Role:        models.STAFF,
					Permissions: []models.PermissionEnum{models.WRITE},
				},
			}

			obj, err := jwtSer.GenerateJwt(
				&models.JwtObj{
					UserId:         save.UUID.String(),
					AccessControls: cred,
				},
				utils.TwoDaysInSeconds,
			)

			dto := model.SchedulePayload{
				StaffId: save.UUID.String(),
				Times: &[]model.ScheduleSegmentPayload{
					{
						IsVisible:     true,
						IsReoccurring: false,
						Start:         logger.Date().Add(time.Duration(1) * time.Hour).Format(utils.TimeFormat),
						Duration:      3600,
					},
					{
						IsVisible:     false,
						IsReoccurring: true,
						Start:         logger.Date().Add(time.Duration(2) * time.Hour).Format(utils.TimeFormat),
						Duration:      3600,
					},
				},
			}

			dtoBytes, err := json.Marshal(dto)
			if err != nil {
				t.Errorf("failed to marshal SchedulePayload: %s", err)
			}

			// handler to test
			NewScheduleHandler(mux, m, logger, s).Register()

			for i := 0; i < 2; i++ {
				req := httptest.NewRequest(http.MethodPost, "/schedule", bytes.NewBuffer(dtoBytes))
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
					if rr.Code != http.StatusBadRequest {
						t.Errorf("index %v expected status code %d, got %d", i, http.StatusBadRequest, rr.Code)
					}
				}
			}
		})

		t.Run("reject creation schedule bleeds into the next day", func(t *testing.T) {
			t.Parallel()

			tx, fn := setupTest(t)
			defer fn()

			// given
			mux := http.NewServeMux()
			prov := stores.MockLiveTransactionProvider(logger, tx)
			adapters := stores.NewAdapters(logger, tx, prov)
			jwtSer := auth.NewJwtService(logger, env)
			m := &middleware.Middleware{Logger: logger, Auth: jwtSer, Param: env.CookieParam}
			s := schedule.NewScheduleService(logger, adapters)

			save, err := preSaveStaff(adapters)
			if err != nil {
				t.Error(err)
			}

			cred := []models.RolePermission{
				{
					Role:        models.STAFF,
					Permissions: []models.PermissionEnum{models.WRITE},
				},
			}

			obj, err := jwtSer.GenerateJwt(
				&models.JwtObj{
					UserId:         save.UUID.String(),
					AccessControls: cred,
				},
				utils.TwoDaysInSeconds,
			)

			d := logger.Date()
			date := time.Date(d.Year(), d.Month(), d.Day(), 23, 0, 0, d.Nanosecond(), logger.Timezone())
			dto := model.SchedulePayload{
				StaffId: save.UUID.String(),
				Times: &[]model.ScheduleSegmentPayload{
					{
						IsVisible:     true,
						IsReoccurring: false,
						Start:         date.Format(utils.TimeFormat),
						Duration:      2 * 60 * 60,
					},
				},
			}

			dtoBytes, err := json.Marshal(dto)
			if err != nil {
				t.Errorf("failed to marshal SchedulePayload: %s", err)
			}

			// handler to test
			NewScheduleHandler(mux, m, logger, s).Register()

			req := httptest.NewRequest(http.MethodPost, "/schedule", bytes.NewBuffer(dtoBytes))
			req.AddCookie(&http.Cookie{Name: env.CookieParam.CookieName, Value: obj.Token})
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()

			mux.ServeHTTP(rr, req)

			if rr.Code != http.StatusBadRequest {
				t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
			}
		})

		t.Run("success saving schedule & retrieving schedules", func(t *testing.T) {
			t.Parallel()

			tx, fn := setupTest(t)
			defer fn()

			// given
			mux := http.NewServeMux()
			prov := stores.MockLiveTransactionProvider(logger, tx)
			adapters := stores.NewAdapters(logger, tx, prov)
			jwtSer := auth.NewJwtService(logger, env)
			m := &middleware.Middleware{Logger: logger, Auth: jwtSer, Param: env.CookieParam}
			s := schedule.NewScheduleService(logger, adapters)

			save, err := preSaveStaff(adapters)
			if err != nil {
				t.Error(err)
			}

			cred := []models.RolePermission{
				{
					Role:        models.STAFF,
					Permissions: []models.PermissionEnum{models.WRITE},
				},
			}

			obj, err := jwtSer.GenerateJwt(
				&models.JwtObj{
					UserId:         save.UUID.String(),
					AccessControls: cred,
				},
				utils.TwoDaysInSeconds,
			)

			dto := model.SchedulePayload{
				StaffId: save.UUID.String(),
				Times: &[]model.ScheduleSegmentPayload{
					{
						IsVisible:     true,
						IsReoccurring: false,
						Start:         logger.Date().Add(time.Duration(1) * time.Hour).Format(utils.TimeFormat),
						Duration:      3600,
					},
					{
						IsVisible:     false,
						IsReoccurring: true,
						Start:         logger.Date().Add(time.Duration(2) * time.Hour).Format(utils.TimeFormat),
						Duration:      3600,
					},
				},
			}

			dtoBytes, err := json.Marshal(dto)
			if err != nil {
				t.Errorf("failed to marshal SchedulePayload: %s", err)
			}

			// handler to test
			NewScheduleHandler(mux, m, logger, s).Register()

			// route to test
			req := httptest.NewRequest(http.MethodPost, "/schedule", bytes.NewBuffer(dtoBytes))
			req.AddCookie(&http.Cookie{Name: env.CookieParam.CookieName, Value: obj.Token})
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()

			mux.ServeHTTP(rr, req)

			// assert
			if rr.Code != http.StatusCreated {
				t.Errorf("expected status code %d, got %d", http.StatusCreated, rr.Code)
			}

			// route to test
			url := fmt.Sprintf("/schedule?month=%v&year=%v&timezone=%s", int(logger.Date().Month()), logger.Date().Year(), logger.Timezone().String())
			req = httptest.NewRequest(http.MethodGet, url, nil)
			req.AddCookie(&http.Cookie{Name: env.CookieParam.CookieName, Value: obj.Token})
			req.Header.Set("Content-Type", "application/json")

			rr = httptest.NewRecorder()

			mux.ServeHTTP(rr, req)

			// assert
			if rr.Code != http.StatusOK {
				t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
			}
		})
	})
}

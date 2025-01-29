package schedule

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/iTchTheRightSpot/erp-golang/config"
	"github.com/iTchTheRightSpot/erp-golang/database"
	"github.com/iTchTheRightSpot/erp-golang/pkg/handlers"
	"github.com/iTchTheRightSpot/erp-golang/pkg/middleware"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models"
	model "github.com/iTchTheRightSpot/erp-golang/pkg/models/schedule"
	"github.com/iTchTheRightSpot/erp-golang/pkg/services/auth"
	"github.com/iTchTheRightSpot/erp-golang/pkg/services/schedule"
	"github.com/iTchTheRightSpot/erp-golang/pkg/stores"
	"github.com/iTchTheRightSpot/erp-golang/utils"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
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

func TestScheduleHandler(t *testing.T) {
	del := func() {
		if err := handlers.DeleteAll(db); err != nil {
			t.Log("failed to delete all from db after test ", err.Error())
		}
	}

	t.Cleanup(del)

	// given
	mux := http.NewServeMux()
	logger := utils.NewMockLogger()
	prov := stores.NewTransactionProvider(logger, db)
	adapters := stores.NewAdapters(logger, db, prov)
	jwtSer := auth.NewJwtService(logger, env)
	m := &middleware.Middleware{Logger: logger, Auth: jwtSer, Param: env.CookieParam}
	s := schedule.NewScheduleService(logger, adapters)

	ctx := context.Background()
	staff1, err := handlers.PreSaveStaff(ctx, adapters)
	if err != nil {
		t.Error(err)
	}

	staff2, err := handlers.PreSaveStaff(ctx, adapters)
	if err != nil {
		t.Error(err)
	}

	NewScheduleHandler(mux, m, logger, s).Register()

	date := logger.Date()
	d := time.Date(date.Year(), date.Month(), date.Day(), 9, 0, 0, 0, date.Location()).Add(24 * time.Hour)

	t.Run("flow of application", func(t *testing.T) {

		t.Run("success. create schedule", func(t *testing.T) {
			cred := []models.RolePermissionEnum{
				{Role: models.STAFF, Permissions: []models.PermissionEnum{models.WRITE}},
			}

			obj, err := jwtSer.Encode(
				&models.JwtObj{
					UserId:         staff1.UUID.String(),
					AccessControls: cred,
				},
				utils.TwoDaysInSeconds,
			)

			dto := model.SchedulePayload{
				StaffId: staff1.UUID.String(),
				Times: &[]model.ScheduleSegmentPayload{
					{
						IsVisible:     true,
						IsReoccurring: false,
						Start:         d.Add(time.Hour).Format(utils.TimeFormat),
						Duration:      3600,
					},
				},
			}

			dtoBytes, err := json.Marshal(dto)
			if err != nil {
				t.Errorf("failed to marshal SchedulePayload: %s", err.Error())
				return
			}

			// create schedule route to test
			req := httptest.NewRequest(http.MethodPost, "/schedule", bytes.NewBuffer(dtoBytes))
			req.AddCookie(&http.Cookie{Name: env.CookieParam.CookieName, Value: obj.Token})
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()

			mux.ServeHTTP(rr, req)

			// assert
			if rr.Code != http.StatusCreated {
				t.Errorf("expected status code %d, got %d", http.StatusCreated, rr.Code)
			}

			// staff2
			dto.StaffId = staff2.UUID.String()

			dtoBytes, err = json.Marshal(dto)
			if err != nil {
				t.Errorf("failed to marshal SchedulePayload: %s", err)
				return
			}

			// create schedule route to test
			req = httptest.NewRequest(http.MethodPost, "/schedule", bytes.NewBuffer(dtoBytes))
			req.AddCookie(&http.Cookie{Name: env.CookieParam.CookieName, Value: obj.Token})
			req.Header.Set("Content-Type", "application/json")

			rr = httptest.NewRecorder()

			mux.ServeHTTP(rr, req)

			// assert
			if rr.Code != http.StatusCreated {
				t.Errorf("expected status code %d, got %d", http.StatusCreated, rr.Code)
			}
		})

		t.Run("success. staff accessing his/her schedule", func(t *testing.T) {
			cred := []models.RolePermissionEnum{
				{Role: models.STAFF, Permissions: []models.PermissionEnum{models.READ}},
			}

			obj, _ := jwtSer.Encode(
				&models.JwtObj{
					UserId:         staff1.UUID.String(),
					AccessControls: cred,
				},
				utils.TwoDaysInSeconds,
			)

			url := fmt.Sprintf("/schedule?month=%v&year=%v&timezone=%s", int(d.Month()), d.Year(), logger.Timezone().String())
			req := httptest.NewRequest(http.MethodGet, url, nil)
			req.AddCookie(&http.Cookie{Name: env.CookieParam.CookieName, Value: obj.Token})
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()

			mux.ServeHTTP(rr, req)

			// assert
			if rr.Code != http.StatusOK {
				t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
			}

			var arr []model.ScheduleResponse
			if err = json.NewDecoder(rr.Body).Decode(&arr); err != nil {
				t.Errorf(err.Error())
			}

			if len(arr) < 1 {
				t.Errorf("expect length > 0 len %v", len(arr))
				t.Errorf("%v", arr)
			}
		})

		t.Run("reject. staff without WRITE permission trying to view another staffs schedule", func(t *testing.T) {
			cred := []models.RolePermissionEnum{
				{Role: models.STAFF, Permissions: []models.PermissionEnum{models.READ}},
			}

			obj, _ := jwtSer.Encode(
				&models.JwtObj{
					UserId:         staff1.UUID.String(),
					AccessControls: cred,
				},
				utils.TwoDaysInSeconds,
			)

			url := fmt.Sprintf("/schedule/staff?month=%v&year=%v&timezone=%s", int(d.Month()), d.Year(), logger.Timezone().String())
			req := httptest.NewRequest(http.MethodGet, url, nil)
			req.AddCookie(&http.Cookie{Name: env.CookieParam.CookieName, Value: obj.Token})
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()

			mux.ServeHTTP(rr, req)

			// assert
			if rr.Code != http.StatusForbidden {
				t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
			}
		})

		t.Run("success. staff with WRITE permission trying to see another staffs schedule", func(t *testing.T) {
			cred := []models.RolePermissionEnum{
				{Role: models.STAFF, Permissions: []models.PermissionEnum{models.WRITE}},
			}

			obj, _ := jwtSer.Encode(
				&models.JwtObj{
					UserId:         staff1.UUID.String(),
					AccessControls: cred,
				},
				utils.TwoDaysInSeconds,
			)

			url := fmt.Sprintf("/schedule/staff?month=%v&year=%v&timezone=%s", int(d.Month()), d.Year(), logger.Timezone().String())
			req := httptest.NewRequest(http.MethodGet, url, nil)
			req.AddCookie(&http.Cookie{Name: env.CookieParam.CookieName, Value: obj.Token})
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()

			mux.ServeHTTP(rr, req)

			// assert
			if rr.Code != http.StatusOK {
				t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
			}

			var arr []model.ScheduleResponse
			if err = json.NewDecoder(rr.Body).Decode(&arr); err != nil {
				t.Errorf(err.Error())
			}

			if len(arr) < 1 {
				t.Errorf("expect length > 0 len %v", len(arr))
				t.Errorf("%v", arr)
			}
		})
	})

	t.Run("reject. concurrent CREATE schedule request. duplicate", func(t *testing.T) {
		t.Cleanup(del)

		var wg sync.WaitGroup
		var mu sync.Mutex
		var statusArr []int
		var errArr []string

		cred := []models.RolePermissionEnum{
			{Role: models.STAFF, Permissions: []models.PermissionEnum{models.WRITE}},
		}

		obj, err := jwtSer.Encode(
			&models.JwtObj{
				UserId:         staff1.UUID.String(),
				AccessControls: cred,
			},
			utils.TwoDaysInSeconds,
		)

		nt := d.Add(48 * time.Hour)
		dto := model.SchedulePayload{
			StaffId: staff1.UUID.String(),
			Times: &[]model.ScheduleSegmentPayload{
				{
					IsVisible:     true,
					IsReoccurring: false,
					Start:         nt.Add(time.Hour).Format(utils.TimeFormat),
					Duration:      3600,
				},
				{
					IsVisible:     false,
					IsReoccurring: true,
					Start:         nt.Add(3 * time.Hour).Format(utils.TimeFormat),
					Duration:      3600,
				},
			},
		}

		dtoBytes, err := json.Marshal(dto)
		if err != nil {
			t.Errorf("failed to marshal SchedulePayload: %s", err)
			return
		}

		randNum := rand.Intn(10-2) + 2

		for i := 0; i < randNum; i++ {
			wg.Add(1)

			go func(barr []byte) {
				defer wg.Done()

				req := httptest.NewRequest(http.MethodPost, "/schedule", bytes.NewBuffer(barr))
				req.AddCookie(&http.Cookie{Name: env.CookieParam.CookieName, Value: obj.Token})
				req.Header.Set("Content-Type", "application/json")

				rr := httptest.NewRecorder()
				mux.ServeHTTP(rr, req)

				mu.Lock()
				statusArr = append(statusArr, rr.Code)
				errArr = append(errArr, rr.Body.String())
				mu.Unlock()
			}(dtoBytes)
		}

		wg.Wait()

		// assert
		num := handlers.CountResponseStatus(statusArr, 201)
		if num != 1 {
			t.Errorf("expect 1 given %v", num)
			t.Errorf("%v", statusArr)
			t.Errorf("%v", errArr)
		}

		num = handlers.CountResponseStatus(statusArr, 409)
		if num != (randNum - 1) {
			t.Errorf("expect %v given %v", randNum-1, num)
			t.Errorf("%v", statusArr)
			t.Errorf("%v", errArr)
		}
	})

	t.Run("reject. CREATE schedule bleeds into the next day", func(t *testing.T) {
		staff1, _ = handlers.PreSaveStaff(ctx, adapters)

		cred := []models.RolePermissionEnum{
			{Role: models.STAFF, Permissions: []models.PermissionEnum{models.WRITE}},
		}

		obj, _ := jwtSer.Encode(
			&models.JwtObj{
				UserId:         staff1.UUID.String(),
				AccessControls: cred,
			},
			utils.TwoDaysInSeconds,
		)

		d = d.Add(240 * time.Hour) // add 10 days
		dat := time.Date(d.Year(), d.Month(), d.Day(), 23, 0, 0, d.Nanosecond(), logger.Timezone())
		dto := model.SchedulePayload{
			StaffId: staff1.UUID.String(),
			Times: &[]model.ScheduleSegmentPayload{
				{
					IsVisible:     true,
					IsReoccurring: false,
					Start:         dat.Format(utils.TimeFormat),
					Duration:      2 * 60 * 60,
				},
			},
		}

		dtoBytes, err := json.Marshal(dto)
		if err != nil {
			t.Errorf("failed to marshal SchedulePayload: %s", err)
		}

		req := httptest.NewRequest(http.MethodPost, "/schedule", bytes.NewBuffer(dtoBytes))
		req.AddCookie(&http.Cookie{Name: env.CookieParam.CookieName, Value: obj.Token})
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})
}

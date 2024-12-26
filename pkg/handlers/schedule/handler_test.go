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

func setupTest(t *testing.T) (*sql.Tx, func()) {
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("failed to start transaction: %v", err)
	}

	return tx, func() {
		if err = tx.Rollback(); err != nil {
			return
		}
	}
}

func TestScheduleHandler(t *testing.T) {
	logger := utils.NewMockLogger()

	del := func() {
		if err := handlers.DeleteAll(db); err != nil {
			t.Errorf(err.Error())
		}
	}

	t.Run("reject request duplicate schedule", func(t *testing.T) {
		var wg sync.WaitGroup
		var mu sync.Mutex
		var statusArr []int
		var errArr []string

		ctx := context.Background()
		mux := http.NewServeMux()
		prov := stores.NewTransactionProvider(logger, db)
		adapters := stores.NewAdapters(logger, db, prov)
		jwtSer := auth.NewJwtService(logger, env)
		ware := &middleware.Middleware{Logger: logger, Auth: jwtSer, Param: env.CookieParam}
		s := schedule.NewScheduleService(logger, adapters)

		// setup dependencies
		staf, err := handlers.PreSaveStaff(ctx, adapters)
		if err != nil {
			t.Errorf("preSaveStaff failed: %v", err)
			return
		}

		cred := []models.RolePermission{
			{Role: models.STAFF, Permissions: []models.PermissionEnum{models.WRITE}},
		}

		obj, err := jwtSer.GenerateJwt(
			&models.JwtObj{
				UserId:         staf.UUID.String(),
				AccessControls: cred,
			},
			utils.TwoDaysInSeconds,
		)

		date := logger.Date()
		d := time.Date(date.Year(), date.Month(), date.Day(), 9, 0, 0, 0, date.Location()).Add(24 * time.Hour)

		dto := model.SchedulePayload{
			StaffId: staf.UUID.String(),
			Times: &[]model.ScheduleSegmentPayload{
				{
					IsVisible:     true,
					IsReoccurring: false,
					Start:         d.Add(time.Hour).Format(utils.TimeFormat),
					Duration:      3600,
				},
				{
					IsVisible:     false,
					IsReoccurring: true,
					Start:         d.Add(3 * time.Hour).Format(utils.TimeFormat),
					Duration:      3600,
				},
			},
		}

		dtoBytes, err := json.Marshal(dto)
		if err != nil {
			t.Errorf("failed to marshal SchedulePayload: %s", err)
			return
		}

		// initialize routes
		NewScheduleHandler(mux, ware, logger, s).Register()

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

		t.Cleanup(del)
	})

	t.Run("reject creation schedule bleeds into the next day", func(t *testing.T) {
		tx, fn := setupTest(t)
		defer fn()

		// given
		ctx := context.Background()
		mux := http.NewServeMux()
		prov := stores.MockLiveTransactionProvider(logger, tx)
		adapters := stores.NewAdapters(logger, tx, prov)
		jwtSer := auth.NewJwtService(logger, env)
		m := &middleware.Middleware{Logger: logger, Auth: jwtSer, Param: env.CookieParam}
		s := schedule.NewScheduleService(logger, adapters)

		save, err := handlers.PreSaveStaff(ctx, adapters)
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
		tx, fn := setupTest(t)
		defer fn()

		// given
		ctx := context.Background()
		mux := http.NewServeMux()
		prov := stores.MockLiveTransactionProvider(logger, tx)
		adapters := stores.NewAdapters(logger, tx, prov)
		jwtSer := auth.NewJwtService(logger, env)
		m := &middleware.Middleware{Logger: logger, Auth: jwtSer, Param: env.CookieParam}
		s := schedule.NewScheduleService(logger, adapters)

		save, err := handlers.PreSaveStaff(ctx, adapters)
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

		date := logger.Date()
		d := time.Date(date.Year(), date.Month(), date.Day(), 9, 0, 0, 0, date.Location())
		d = d.Add(time.Duration(48) * time.Hour)

		dto := model.SchedulePayload{
			StaffId: save.UUID.String(),
			Times: &[]model.ScheduleSegmentPayload{
				{
					IsVisible:     true,
					IsReoccurring: false,
					Start:         d.Add(time.Duration(1) * time.Hour).Format(utils.TimeFormat),
					Duration:      3600,
				},
				{
					IsVisible:     false,
					IsReoccurring: true,
					Start:         d.Add(time.Duration(3) * time.Hour).Format(utils.TimeFormat),
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
}

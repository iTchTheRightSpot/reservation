package reservation

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/iTchTheRightSpot/erp-golang/config"
	"github.com/iTchTheRightSpot/erp-golang/database"
	scheduleHandler "github.com/iTchTheRightSpot/erp-golang/pkg/handlers/schedule"
	"github.com/iTchTheRightSpot/erp-golang/pkg/middleware"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/profile"
	model "github.com/iTchTheRightSpot/erp-golang/pkg/models/reservation"
	scheduleModel "github.com/iTchTheRightSpot/erp-golang/pkg/models/schedule"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/service"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/staff"
	pkg "github.com/iTchTheRightSpot/erp-golang/pkg/services"
	"github.com/iTchTheRightSpot/erp-golang/pkg/services/auth"
	"github.com/iTchTheRightSpot/erp-golang/pkg/services/mail"
	"github.com/iTchTheRightSpot/erp-golang/pkg/services/reservation"
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

func TestReservationHandler(t *testing.T) {
	logger := utils.NewMockLogger()

	baseDate := logger.Date()
	startTime := time.Date(baseDate.Year(), baseDate.Month(), baseDate.Day(), 9, 0, 0, 0, logger.Timezone())
	newTime := startTime.Add(24 * time.Hour)

	t.Run("single request make reservation", func(t *testing.T) {
		tx, fn := setupTest(t)
		defer fn()

		// setup dependencies
		prov := stores.MockLiveTransactionProvider(logger, tx)
		adapters := stores.NewAdapters(logger, tx, prov)
		mux, savedStaff, saveService, req, rr, payload := reservationFlow(t, logger, adapters, newTime)

		// create reservation
		createBody := model.ReservationPayload{
			StaffId:  savedStaff.UUID.String(),
			Name:     "user-name",
			Email:    "user@email.com",
			Address:  "123 transylvania",
			Phone:    "0123456789",
			Services: []string{saveService.Name},
			Time:     payload[0].Times[1],
		}

		createBodyBytes, err := json.Marshal(createBody)
		if err != nil {
			t.Fatalf("failed to marshal SchedulePayload: %s", err)
		}

		req = httptest.NewRequest(http.MethodPost, "/reservation", bytes.NewBuffer(createBodyBytes))
		req.Header.Set("Content-Type", "application/json")
		rr = httptest.NewRecorder()
		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusCreated {
			t.Errorf("expected status code %d, got %d", http.StatusCreated, rr.Code)
		}
	})

	t.Run("concurrent request to reserve the sametime", func(t *testing.T) {
		defer func() {
			if err := deleteAll(); err != nil {
				t.Errorf(err.Error())
			}
		}()

		// setup dependencies
		prov := stores.NewTransactionProvider(logger, db)
		adapters := stores.NewAdapters(logger, db, prov)
		mux, savedStaff, saveService, _, _, payload := reservationFlow(t, logger, adapters, newTime)

		createBody := model.ReservationPayload{
			StaffId:  savedStaff.UUID.String(),
			Name:     "user-name",
			Email:    fmt.Sprintf("%s@email.com", uuid.NewString()),
			Address:  "123 transylvania",
			Phone:    "0123456789",
			Services: []string{saveService.Name},
			Time:     payload[0].Times[0],
		}

		createBodyBytes, _ := json.Marshal(createBody)

		// create reservation
		var wg sync.WaitGroup
		var mu sync.Mutex
		var statusArr []int
		var errArr []string

		randNum := rand.Intn(5-2) + 2

		for idx := 0; idx < randNum; idx++ {
			wg.Add(1)

			go func(i int) {
				defer wg.Done()

				req := httptest.NewRequest(http.MethodPost, "/reservation", bytes.NewBuffer(createBodyBytes))
				req.Header.Set("Content-Type", "application/json")
				rr := httptest.NewRecorder()
				mux.ServeHTTP(rr, req)

				mu.Lock()
				statusArr = append(statusArr, rr.Code)
				errArr = append(errArr, rr.Body.String())
				mu.Unlock()
			}(idx)
		}

		wg.Wait()

		// assert
		num := count(statusArr, 201)
		if num != 1 {
			t.Errorf("expect 1 given %v", num)
			t.Errorf("%v", statusArr)
			t.Errorf("%v", errArr)
		}

		num = inRange(statusArr)
		if num != (randNum - 1) {
			t.Errorf("expect %v given %v", randNum-1, num)
			t.Errorf("%v", statusArr)
			t.Errorf("%v", errArr)
		}
	})
}

func preSaveStaff(a *stores.Adapters) (*staff.Staff, error) {
	ctx := context.Background()

	p := profile.Profile{
		Firstname: "erp",
		Lastname:  "erp",
		Email:     fmt.Sprintf("%s@email.com", uuid.NewString()),
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

func preSaveService(a *stores.Adapters) (*service.ServiceEntity, error) {
	return a.ServiceStore.Save(context.Background(), &service.ServiceEntity{
		Name:        uuid.New().String(),
		Price:       19.56,
		Duration:    3600,
		CleanUpTime: 30 * 60,
	})
}

func linkServiceToStaff(a *stores.Adapters, sta *staff.Staff, staSer *service.ServiceEntity) error {
	_, err := a.StaffServiceStore.Save(context.Background(), &staff.StaffServiceEntity{
		StaffId:   sta.StaffId,
		ServiceId: staSer.ServiceId,
	})
	return err
}

func deleteAll() error {
	if _, err := db.Exec("TRUNCATE schedule, staff, profile, service, staff_service, reservation, reservation_service CASCADE"); err != nil {
		return err
	}
	return nil
}

func count(arr []int, status int) int {
	var count int
	for _, num := range arr {
		if num == status {
			count += 1
		}
	}
	return count
}

func inRange(arr []int) int {
	var n int
	for _, num := range arr {
		if num >= 400 && num <= 500 {
			n += 1
		}
	}
	return n
}

func reservationFlow(t *testing.T, logger utils.ILogger, adapters *stores.Adapters, newTime time.Time) (*http.ServeMux, *staff.Staff, *service.ServiceEntity, *http.Request, *httptest.ResponseRecorder, []model.ReservationTimeSlots) {
	mux := http.NewServeMux()
	jwtService := auth.NewJwtService(logger, env)
	ware := &middleware.Middleware{Logger: logger, Auth: jwtService, Param: env.CookieParam}
	scheduleService := schedule.NewScheduleService(logger, adapters)
	cache := pkg.NewInMemoryCache[string, []model.ReservationTimeSlots](logger, 30, 30)
	mailService := &mail.MockMailService{}
	reservationService := reservation.NewReservationService(logger, adapters, cache, mailService)

	savedStaff, err := preSaveStaff(adapters)
	if err != nil {
		t.Fatalf("preSaveStaff failed: %v", err)
	}

	saveService, err := preSaveService(adapters)
	if err != nil {
		t.Fatalf("preSaveService failed: %v", err)
	}

	if err = linkServiceToStaff(adapters, savedStaff, saveService); err != nil {
		t.Fatalf(err.Error())
	}

	cred := []models.RolePermission{
		{Role: models.STAFF, Permissions: []models.PermissionEnum{models.WRITE}},
	}

	obj, err := jwtService.GenerateJwt(
		&models.JwtObj{
			UserId:         savedStaff.UUID.String(),
			AccessControls: cred,
		},
		utils.TwoDaysInSeconds,
	)

	dto := scheduleModel.SchedulePayload{
		StaffId: savedStaff.UUID.String(),
		Times: &[]scheduleModel.ScheduleSegmentPayload{
			{
				IsVisible:     true,
				IsReoccurring: false,
				Start:         newTime.Format(utils.TimeFormat),
				Duration:      8 * 60 * 60,
			},
		},
	}

	dtoBytes, err := json.Marshal(dto)
	if err != nil {
		t.Fatalf("failed to marshal SchedulePayload: %s", err)
	}

	// handler to test
	scheduleHandler.NewScheduleHandler(mux, ware, logger, scheduleService).Register()

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

	// retrieve valid reservation times
	NewReservationHandler(mux, logger, reservationService).Register()

	url := fmt.Sprintf("/reservation?day=%v&month=%v&year=%v&staff_id=%s&service=%s", 1, int(newTime.Month()), newTime.Year(), savedStaff.UUID.String(), saveService.Name)
	req = httptest.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	var payload []model.ReservationTimeSlots
	if err = json.Unmarshal(rr.Body.Bytes(), &payload); err != nil {
		t.Errorf(err.Error())
	}
	return mux, savedStaff, saveService, req, rr, payload
}

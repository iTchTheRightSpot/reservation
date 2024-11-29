package reservation

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
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
	"github.com/iTchTheRightSpot/erp-golang/pkg/services/auth"
	"github.com/iTchTheRightSpot/erp-golang/pkg/services/reservation"
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

	db, err = database.ConnectToPostgre(e.DbConnectionString)
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

func preSaveService(a *stores.Adapters) (*service.ServiceEntity, error) {
	return a.ServiceStore.Save(context.Background(), &service.ServiceEntity{
		Name:        "erp",
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

func TestReservationHandler(t *testing.T) {
	t.Parallel()

	logger := utils.NewMockLogger()

	baseDate := logger.Date()
	startTime := time.Date(baseDate.Year(), baseDate.Month(), baseDate.Day(), 9, 0, 0, 0, logger.Timezone())
	newTime := startTime.Add(24 * time.Hour)

	t.Run("describe creating a reservation", func(t *testing.T) {
		t.Parallel()

		t.Run("reject request invalid service", func(t *testing.T) {
			t.Parallel()

			tx, fn := setupTest(t)
			defer fn()

			// given
			mux := http.NewServeMux()
			prov := stores.MockLiveTransactionProvider(logger, tx)
			adapters := stores.NewAdapters(logger, tx, prov)
			s := reservation.NewReservationService(logger, adapters)

			savedStaff, err := preSaveStaff(adapters)
			if err != nil {
				t.Error(err)
			}

			erp := "erp"
			payload := model.ReservationPayload{
				StaffId:  savedStaff.UUID.String(),
				Name:     "temp user",
				Email:    "temp-user@email.com",
				Address:  "123 transylvania",
				Phone:    "0123456789",
				Services: []*string{&erp},
				Time:     "19800",
			}

			dtoBytes, err := json.Marshal(payload)
			if err != nil {
				t.Errorf("failed to marshal SchedulePayload: %s", err)
			}

			// handler to test
			NewReservationHandler(mux, logger, s).Register()

			// route to test
			req := httptest.NewRequest(http.MethodPost, "/reservation", bytes.NewBuffer(dtoBytes))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()

			mux.ServeHTTP(rr, req)

			// assert
			if rr.Code != http.StatusBadRequest {
				t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
			}

			var errBody utils.ErrorResponse
			if err = json.Unmarshal(rr.Body.Bytes(), &errBody); err != nil {
				t.Fatalf("failed to unmarshal response: %v", err)
			}

			if errBody.Message != "1 or more services were not found for selected staff" {
				t.Errorf("unexpected response: %+v", errBody)
			}
		})
	})

	t.Run("flow to create a reservation", func(t *testing.T) {
		t.Parallel()

		t.Run("single request same timezone", func(t *testing.T) {
			t.Parallel()

			tx, fn := setupTest(t)
			defer fn()

			// setup dependencies
			mux := http.NewServeMux()
			prov := stores.MockLiveTransactionProvider(logger, tx)
			adapters := stores.NewAdapters(logger, tx, prov)
			jwtService := auth.NewJwtService(logger, env)
			ware := &middleware.Middleware{Logger: logger, Auth: jwtService, Param: env.CookieParam}
			scheduleService := schedule.NewScheduleService(logger, adapters)
			reservationService := reservation.NewReservationService(logger, adapters)

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

			req = httptest.NewRequest(http.MethodGet, "/reservation", nil)
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

			// create reservation
			createBody := model.ReservationPayload{
				StaffId:  savedStaff.UUID.String(),
				Name:     "user-name",
				Email:    "user@email.com",
				Address:  "123 transylvania",
				Phone:    "0123456789",
				Services: []*string{&saveService.Name},
				Time:     payload[0].Times[0],
			}

			createBodyBytes, err := json.Marshal(createBody)
			if err != nil {
				t.Fatalf("failed to marshal SchedulePayload: %s", err)
			}

			req = httptest.NewRequest(http.MethodGet, "/reservation", bytes.NewBuffer(createBodyBytes))
			req.Header.Set("Content-Type", "application/json")
			rr = httptest.NewRecorder()
			mux.ServeHTTP(rr, req)

			if rr.Code != http.StatusCreated {
				t.Errorf("expected status code %d, got %d", http.StatusCreated, rr.Code)
			}
		})
	})
}

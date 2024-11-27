package reservation

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/iTchTheRightSpot/erp-golang/config"
	"github.com/iTchTheRightSpot/erp-golang/database"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/profile"
	model "github.com/iTchTheRightSpot/erp-golang/pkg/models/reservation"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/service"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/staff"
	"github.com/iTchTheRightSpot/erp-golang/pkg/services/reservation"
	"github.com/iTchTheRightSpot/erp-golang/pkg/stores"
	"github.com/iTchTheRightSpot/erp-golang/utils"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
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

func preSaveService(a *stores.Adapters) (*service.ServiceEntity, error) {
	return a.ServiceStore.Save(context.Background(), &service.ServiceEntity{
		Name:        "erp",
		Price:       19.56,
		Duration:    3600,
		CleanUpTime: 30 * 60,
	})
}

func TestReservationHandler(t *testing.T) {
	t.Parallel()

	logger := utils.NewMockLogger()

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
}

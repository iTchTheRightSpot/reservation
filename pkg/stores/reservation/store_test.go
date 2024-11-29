package reservation

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/iTchTheRightSpot/erp-golang/config"
	"github.com/iTchTheRightSpot/erp-golang/database"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/profile"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/reservation"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/service"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/staff"
	profileStore "github.com/iTchTheRightSpot/erp-golang/pkg/stores/profile"
	serviceStore "github.com/iTchTheRightSpot/erp-golang/pkg/stores/service"
	staffStore "github.com/iTchTheRightSpot/erp-golang/pkg/stores/staff"
	"github.com/iTchTheRightSpot/erp-golang/utils"
	"log"
	"testing"
	"time"
)

var db *sql.DB

func TestMain(m *testing.M) {
	secret := config.SecretVariables{}
	env, err := secret.Config()
	if err != nil {
		log.Fatal(err)
	}

	db, err = database.ConnectToPostgre(env.DbConnectionString)
	if err != nil {
		log.Fatal(err)
	}

	// close db connection
	defer func(db *sql.DB) {
		if err := db.Close(); err != nil {
			log.Printf("db connection did not close after tests")
			return
		}
	}(db)

	// run tests
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

func preSaveStaff(ps profileStore.IProfileStore, st staffStore.IStaffStore) (*staff.Staff, error) {
	ctx := context.Background()

	p := profile.Profile{
		Firstname: "erp",
		Lastname:  "erp",
		Email:     "erp@email.com",
	}

	if _, err := ps.Save(ctx, &p); err != nil {
		return nil, err
	}

	s := staff.Staff{
		UUID:      uuid.New(),
		ProfileId: &p.ProfileId,
	}

	if _, err := st.Save(ctx, &s); err != nil {
		return nil, err
	}

	return &s, nil
}

func preSaveService(a serviceStore.IServiceStore) (*service.ServiceEntity, error) {
	return a.Save(context.Background(), &service.ServiceEntity{
		Name:        "erp",
		Price:       19.56,
		Duration:    3600,
		CleanUpTime: 30 * 60,
	})
}

func TestReservationStore(t *testing.T) {
	mockLogger := utils.NewMockLogger()

	t.Run("should save reservation && service_reservation && count reservations for staff by time & statuses", func(t *testing.T) {
		tx, fn := setupTest(t)
		defer fn()

		// given
		ctx := context.Background()

		store := NewReservationStore(mockLogger, tx)
		reservationServStore := NewReservationServiceStore(mockLogger, tx)
		erpStaff, err := preSaveStaff(profileStore.NewProfileStore(mockLogger, tx), staffStore.NewStaffStore(mockLogger, tx))
		if err != nil {
			t.Errorf(err.Error())
		}

		erp, err := preSaveService(serviceStore.NewServiceStore(mockLogger, tx))
		if err != nil {
			t.Errorf(err.Error())
		}

		start := mockLogger.Date().Add(1 * time.Hour)
		end := start.Add(1 * time.Hour)

		r := &reservation.Reservation{
			StaffId:      erpStaff.StaffId,
			Name:         "user",
			Email:        "email@example.com",
			Price:        25.65,
			Status:       reservation.CONFIRMED,
			CreatedAt:    mockLogger.Date(),
			ScheduledFor: start,
			ExpireAt:     end,
		}

		// method to test & assert
		if err = store.SelectForUpdateSave(ctx, r, reservation.CONFIRMED); err != nil {
			t.Errorf(err.Error())
		}

		if r.ReservationId < 1 {
			t.Errorf("expect reservation id to be greater than 1. given %v", r.ReservationId)
		}

		// given
		re := &reservation.ReservationServiceEntity{
			ReservationId: r.ReservationId,
			ServiceId:     erp.ServiceId,
		}

		// method to test & assert
		if err = reservationServStore.Save(ctx, re); err != nil {
			t.Errorf(err.Error())
		}

		if re.JunctionId < 1 {
			t.Errorf("expect junction id to be greater than 1. given %v", r.ReservationId)
		}
	})
}

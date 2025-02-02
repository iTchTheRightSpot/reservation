package reservation

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/iTchTheRightSpot/erp-golang/config"
	"github.com/iTchTheRightSpot/erp-golang/database"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/reservation"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/service_type"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/staff"
	profileStore "github.com/iTchTheRightSpot/erp-golang/pkg/stores/profile"
	serviceStore "github.com/iTchTheRightSpot/erp-golang/pkg/stores/service_type"
	staffStore "github.com/iTchTheRightSpot/erp-golang/pkg/stores/staff"
	"github.com/iTchTheRightSpot/erp-golang/utils"
	"log"
	"testing"
	"time"
)

var db *sql.DB

func TestMain(m *testing.M) {
	secret := config.SecretVariables{}
	env := secret.Config()

	d, err := database.ConnectToPostgres(env.DbConnectionString)
	if err != nil {
		log.Fatal(err)
	}
	db = d

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

func preSaveStaff(ps profileStore.IProfileStore, st staffStore.IStaffStore) (*staff.StaffEntity, error) {
	ctx := context.Background()

	p := models.ProfileEntity{
		Password:  "password",
		Firstname: "erp",
		Lastname:  "erp",
		Email:     "erp@email.com",
	}

	if err := ps.Save(ctx, &p); err != nil {
		return nil, err
	}

	s := staff.StaffEntity{
		UUID:      uuid.New(),
		ProfileId: &p.ProfileId,
	}

	if err := st.Save(ctx, &s); err != nil {
		return nil, err
	}

	return &s, nil
}

func preSaveService(a serviceStore.IServiceTypeStore) (*service_type.ServiceTypeEntity, error) {
	s := service_type.ServiceTypeEntity{
		Name:        "erp",
		Price:       19.56,
		Duration:    3600,
		CleanUpTime: 30 * 60,
	}
	err := a.Save(context.Background(), &s)
	return &s, err
}

func TestReservationStore(t *testing.T) {
	mockLog := utils.NewMockLogger()

	t.Run("should save reservation & count reservations in range", func(t *testing.T) {
		tx, fn := setupTest(t)
		defer fn()

		// given
		ctx := context.Background()

		store := NewReservationStore(mockLog, tx)
		staf, err := preSaveStaff(profileStore.NewProfileStore(mockLog, tx), staffStore.NewStaffStore(mockLog, tx))
		if err != nil {
			t.Errorf(err.Error())
		}

		start := time.Date(mockLog.Date().Year(), mockLog.Date().Month(), 1, 9, 0, 0, 0, mockLog.Timezone())
		r1 := &reservation.Reservation{
			StaffId:      staf.StaffId,
			Name:         "user",
			Email:        "email@example.com",
			Price:        25.65,
			Status:       reservation.CONFIRMED,
			CreatedAt:    mockLog.Date(),
			ScheduledFor: start,
			ExpireAt:     start.Add(1 * time.Hour),
		}

		r2 := &reservation.Reservation{
			StaffId:      staf.StaffId,
			Name:         "user",
			Email:        "email@example.com",
			Price:        25.65,
			Status:       reservation.CONFIRMED,
			CreatedAt:    mockLog.Date(),
			ScheduledFor: start.Add(2 * time.Hour),
			ExpireAt:     start.Add(3 * time.Hour),
		}

		// method to test & assert
		if err = store.Save(ctx, r1); err != nil {
			t.Errorf(err.Error())
		}
		_ = store.Save(ctx, r2)

		// method to test
		count, err := store.CountReservationsInRange(ctx, staf.StaffId, start, start.Add(1*time.Minute), reservation.CONFIRMED)
		if err != nil {
			t.Error(err.Error())
		}

		if count != 1 {
			t.Errorf("expect 1 given %v", count)
		}

		count, _ = store.CountReservationsInRange(ctx, staf.StaffId, start, start.Add(2*time.Hour), reservation.CONFIRMED)
		if count != 2 {
			t.Errorf("expect 2 given %v", count)
		}

		count, _ = store.CountReservationsInRange(ctx, staf.StaffId, start.Add(4*time.Hour), start.Add(6*time.Hour), reservation.CONFIRMED)
		if count != 0 {
			t.Errorf("expect 0 given %v", count)
		}
	})

	t.Run("should save reservation && service_reservation && count reservations for staff by time & statuses", func(t *testing.T) {
		tx, fn := setupTest(t)
		defer fn()

		// given
		ctx := context.Background()

		store := NewReservationStore(mockLog, tx)
		reservationServStore := NewReservationServiceStore(mockLog, tx)
		erpStaff, err := preSaveStaff(profileStore.NewProfileStore(mockLog, tx), staffStore.NewStaffStore(mockLog, tx))
		if err != nil {
			t.Errorf(err.Error())
		}

		erp, err := preSaveService(serviceStore.NewServiceTypeStore(mockLog, tx))
		if err != nil {
			t.Errorf(err.Error())
		}

		start := mockLog.Date().Add(1 * time.Hour)
		end := start.Add(1 * time.Hour)

		r := &reservation.Reservation{
			StaffId:      erpStaff.StaffId,
			Name:         "user",
			Email:        "email@example.com",
			Price:        25.65,
			Status:       reservation.CONFIRMED,
			CreatedAt:    mockLog.Date(),
			ScheduledFor: start,
			ExpireAt:     end,
		}

		// method to test & assert
		if err = store.Save(ctx, r); err != nil {
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

	t.Run("reject saving reservation conflict", func(t *testing.T) {
		tx, fn := setupTest(t)
		defer fn()

		// given
		ctx := context.Background()

		store := NewReservationStore(mockLog, tx)
		erpStaff, err := preSaveStaff(profileStore.NewProfileStore(mockLog, tx), staffStore.NewStaffStore(mockLog, tx))
		if err != nil {
			t.Errorf(err.Error())
		}

		start := mockLog.Date().Add(1 * time.Hour)
		r := &reservation.Reservation{
			StaffId:      erpStaff.StaffId,
			Name:         "user",
			Email:        "email@example.com",
			Price:        25.65,
			Status:       reservation.CONFIRMED,
			CreatedAt:    mockLog.Date(),
			ScheduledFor: start,
			ExpireAt:     start.Add(1 * time.Hour),
		}

		// method to test & assert
		if err = store.Save(ctx, r); err != nil {
			t.Errorf(err.Error())
		}

		if err = store.Save(ctx, r); err == nil {
			t.Errorf("expect %v given nil", err.Error())
		}

		r.ScheduledFor = start.Add(10 * time.Minute)
		if err = store.Save(ctx, r); err == nil {
			t.Errorf("expect %v given nil", err.Error())
		}
	})
}

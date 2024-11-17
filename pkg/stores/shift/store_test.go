package shift

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/iTchTheRightSpot/erp-golang/config"
	"github.com/iTchTheRightSpot/erp-golang/database"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/shift"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/staff"
	staffStore "github.com/iTchTheRightSpot/erp-golang/pkg/stores/staff"
	"github.com/iTchTheRightSpot/erp-golang/utils"
	"log"
	"reflect"
	"testing"
	"time"
)

var dbInstance *sql.DB

func TestMain(m *testing.M) {
	secret := config.SecretVariables{}
	env, err := secret.Config()
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.ConnectToPostgres(env.DbConnectionString)
	if err != nil {
		log.Fatal(err)
	}

	dbInstance = db

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
	tx, err := dbInstance.Begin()
	if err != nil {
		t.Fatalf("failed to start transaction: %v", err)
	}

	rollback := func() {
		if err := tx.Rollback(); err != nil {
			return
		}
	}

	return tx, rollback
}

func TestShiftStore(t *testing.T) {
	t.Parallel()

	mockLog := utils.NewMockLogger()

	t.Run("should save staff shift", func(t *testing.T) {
		t.Parallel()

		tx, fn := setupTest(t)
		defer fn()

		ctx := context.Background()

		// given
		staffObj := staff.Staff{UUID: uuid.New()}
		if _, err := staffStore.NewStaffStore(mockLog, tx).Save(ctx, &staffObj); err != nil {
			t.Errorf("%s", err)
		}

		s := &shift.Shift{
			StaffId: staffObj.StaffId,
			Start:   mockLog.Date(),
			End:     mockLog.Date().Add(time.Duration(8) * time.Hour),
		}

		// method to test
		save, err := NewShiftStore(mockLog, tx).Save(ctx, s)
		if err != nil {
			t.Errorf("shift not saved")
		}

		// assert
		if save.ShiftId < 1 {
			t.Errorf("shift not save as ShiftId is less than 1")
		}

		if !reflect.DeepEqual(save, s) {
			t.Errorf("expect %v to equal %v", s, save)
		}
	})

	t.Run("shift count should be greater than zero for staff", func(t *testing.T) {
		t.Parallel()

		tx, fn := setupTest(t)
		defer fn()

		store := NewShiftStore(mockLog, tx)
		ctx := context.Background()

		// given
		staffObj := staff.Staff{UUID: uuid.New()}
		if _, err := staffStore.NewStaffStore(mockLog, tx).Save(ctx, &staffObj); err != nil {
			t.Errorf("%s", err)
		}

		date := mockLog.Date()
		s := &shift.Shift{
			StaffId: staffObj.StaffId,
			Start:   date,
			End:     mockLog.Date().Add(time.Duration(8) * time.Hour),
		}

		if _, err := store.Save(ctx, s); err != nil {
			t.Errorf("shift not saved")
		}

		// method to test
		count, err := store.CountExistingShiftsForStaff(ctx, staffObj.StaffId, s.Start, s.End)
		if err != nil {
			t.Error(err)
		}

		if count != 1 {
			t.Errorf("expect 1 given %v", count)
		}
	})

	t.Run("count should be zero for shifts for staff", func(t *testing.T) {
		t.Parallel()

		tx, fn := setupTest(t)
		defer fn()

		store := NewShiftStore(mockLog, tx)
		ctx := context.Background()

		staffObj := staff.Staff{UUID: uuid.New()}
		if _, err := staffStore.NewStaffStore(mockLog, tx).Save(ctx, &staffObj); err != nil {
			t.Errorf("%s", err)
		}

		date := mockLog.Date()
		s := &shift.Shift{
			StaffId: staffObj.StaffId,
			Start:   date,
			End:     mockLog.Date().Add(time.Duration(8) * time.Hour),
		}

		if _, err := store.Save(ctx, s); err != nil {
			t.Errorf("shift not saved")
		}

		// method to test
		param1 := s.End.Add(time.Duration(1) * time.Second)
		param2 := s.End.Add(time.Duration(8) * time.Hour)
		count, err := store.CountExistingShiftsForStaff(ctx, staffObj.StaffId, param1, param2)
		if err != nil {
			t.Error(err)
		}

		if count != 0 {
			t.Errorf("expect 0 given %v", count)
		}
	})
}

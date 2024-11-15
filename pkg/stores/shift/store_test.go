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
	mockLog := utils.NewMockLogger()

	t.Run("should save staff shift", func(t *testing.T) {
		con, fn := setupTest(t)
		defer fn()

		ctx := context.Background()

		// given
		staffObj := staff.Staff{StaffUUID: uuid.New()}
		if _, err := staffStore.NewStaffStore(mockLog, con).Save(ctx, &staffObj); err != nil {
			t.Errorf("%s", err)
		}

		s := &shift.Shift{
			StaffId: staffObj.StaffId,
			Start:   mockLog.Date(),
			End:     mockLog.Date().Add(time.Duration(8) * time.Hour),
		}

		// method to test
		save, err := NewShiftStore(mockLog, con).Save(ctx, s)
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
}

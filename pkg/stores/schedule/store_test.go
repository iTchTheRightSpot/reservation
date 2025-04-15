package schedule

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/iTchTheRightSpot/erp-golang/config"
	"github.com/iTchTheRightSpot/erp-golang/database"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/schedule"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/staff"
	staffStore "github.com/iTchTheRightSpot/erp-golang/pkg/stores/staff"
	"github.com/iTchTheRightSpot/utility/utils"
	"log"
	"reflect"
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

	defer func(db *sql.DB) {
		if err := db.Close(); err != nil {
			log.Printf("db connection did not close after tests")
		}
	}(d)

	m.Run()
}

func setupTest(t *testing.T) (*sql.Tx, func()) {
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("failed to start transaction: %v", err)
	}

	rollback := func() {
		if err = tx.Rollback(); err != nil {
			return
		}
	}

	return tx, rollback
}

func TestScheduleStore(t *testing.T) {
	t.Parallel()

	mockLog := utils.DevLogger("UTC")

	t.Run("should save staff schedule", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		tx, fn := setupTest(t)
		defer fn()

		// given
		staffObj := staff.StaffEntity{UUID: uuid.New()}
		if err := staffStore.NewStaffStore(mockLog, tx).Save(ctx, &staffObj); err != nil {
			t.Errorf("%s", err)
		}

		s := &schedule.ScheduleEntity{
			StaffId: staffObj.StaffId,
			Start:   mockLog.Date(),
			End:     mockLog.Date().Add(time.Duration(8) * time.Hour),
		}

		// method to test
		err := NewScheduleStore(mockLog, tx).Save(ctx, s)
		if err != nil {
			t.Errorf("schedule not saved")
		}

		// assert
		if s.ScheduleId < 1 {
			t.Errorf("schedule not save as ScheduleId is less than 1")
		}
	})

	t.Run("schedule count should be greater than zero for staff", func(t *testing.T) {
		t.Parallel()

		tx, fn := setupTest(t)
		defer fn()

		store := NewScheduleStore(mockLog, tx)
		ctx := context.Background()

		// given
		staffObj := staff.StaffEntity{UUID: uuid.New()}
		if err := staffStore.NewStaffStore(mockLog, tx).Save(ctx, &staffObj); err != nil {
			t.Errorf("%s", err)
		}

		date := mockLog.Date()
		s := &schedule.ScheduleEntity{
			StaffId: staffObj.StaffId,
			Start:   date,
			End:     mockLog.Date().Add(time.Duration(8) * time.Hour),
		}

		if err := store.Save(ctx, s); err != nil {
			t.Errorf("schedule not saved")
		}

		// method to test
		count, err := store.CountExistingSchedulesForStaff(ctx, staffObj.StaffId, s.Start, s.End)
		if err != nil {
			t.Error(err)
		}

		if count != 1 {
			t.Errorf("expect 1 given %v", count)
		}
	})

	t.Run("count should be zero for schedules for staff", func(t *testing.T) {
		t.Parallel()

		tx, fn := setupTest(t)
		defer fn()

		store := NewScheduleStore(mockLog, tx)
		ctx := context.Background()

		staffObj := staff.StaffEntity{UUID: uuid.New()}
		if err := staffStore.NewStaffStore(mockLog, tx).Save(ctx, &staffObj); err != nil {
			t.Errorf("%s", err)
		}

		date := mockLog.Date()
		s := &schedule.ScheduleEntity{
			StaffId: staffObj.StaffId,
			Start:   date,
			End:     mockLog.Date().Add(time.Duration(8) * time.Hour),
		}

		if err := store.Save(ctx, s); err != nil {
			t.Errorf("schedule not saved")
		}

		// method to test
		param1 := s.End.Add(time.Duration(1) * time.Second)
		param2 := s.End.Add(time.Duration(8) * time.Hour)
		count, err := store.CountExistingSchedulesForStaff(ctx, staffObj.StaffId, param1, param2)
		if err != nil {
			t.Error(err)
		}

		if count != 0 {
			t.Errorf("expect 0 given %v", count)
		}
	})

	t.Run("should return schedules in range", func(t *testing.T) {
		t.Parallel()

		tx, fn := setupTest(t)
		defer fn()

		store := NewScheduleStore(mockLog, tx)

		ctx := context.Background()

		// given
		staffObj := staff.StaffEntity{UUID: uuid.New()}
		if err := staffStore.NewStaffStore(mockLog, tx).Save(ctx, &staffObj); err != nil {
			t.Errorf("%s", err)
		}

		start := mockLog.Date()
		s := schedule.ScheduleEntity{
			StaffId: staffObj.StaffId,
			Start:   start,
			End:     start.Add(time.Duration(8) * time.Hour),
		}
		if err := store.Save(ctx, &s); err != nil {
			t.Error(err.Error())
		}

		// method to test
		start = time.Date(mockLog.Date().Year(), mockLog.Date().Month(), 1, 0, 0, 0, 0, mockLog.Timezone())
		end := time.Date(start.Year(), start.Month()+1, 0, 23, 59, 59, 999999999, start.Location())
		schs, err := store.SchedulesInRange(ctx, staffObj.StaffId, start, end)

		// assert
		if err != nil {
			t.Error(err.Error())
		}

		if s.ScheduleId < 1 {
			t.Errorf("schedule not save as ScheduleId is less than 1")
		}

		if !reflect.DeepEqual(s, *schs[0]) {
			t.Errorf("expect %v to equal given %v", s, schs[0])
		}
	})

	t.Run("should count schedules in range and visibility", func(t *testing.T) {
		t.Parallel()

		tx, fn := setupTest(t)
		defer fn()

		store := NewScheduleStore(mockLog, tx)

		ctx := context.Background()

		// given
		staffObj := staff.StaffEntity{UUID: uuid.New()}
		if err := staffStore.NewStaffStore(mockLog, tx).Save(ctx, &staffObj); err != nil {
			t.Errorf("%s", err)
		}

		start := mockLog.Date()
		err := store.Save(ctx, &schedule.ScheduleEntity{
			StaffId: staffObj.StaffId,
			Start:   start,
			End:     start.Add(time.Duration(8) * time.Hour),
		})
		if err != nil {
			t.Error(err.Error())
		}

		firstDayOfMonth := time.Date(mockLog.Date().Year(), mockLog.Date().Month(), 1, 0, 0, 0, 0, mockLog.Timezone())
		lastDayOfMonth := firstDayOfMonth.AddDate(0, 1, -1)

		// method to test
		count, err := store.CountSchedulesInRangeAndVisibility(ctx, staffObj.StaffId, firstDayOfMonth, lastDayOfMonth, false)

		// assert
		if err != nil {
			t.Error(err.Error())
		}

		if count != 0 {
			t.Errorf("expect 1 given %v", count)
		}

		count, err = store.CountSchedulesInRangeAndVisibility(ctx, staffObj.StaffId, start, start.Add(1*time.Hour), false)

		// assert
		if err != nil {
			t.Error(err.Error())
		}

		if count != 1 {
			t.Errorf("expect 1 given %v", count)
		}
	})

	t.Run("should return schedules within timeframe", func(t *testing.T) {
		t.Parallel()

		tx, fn := setupTest(t)
		defer fn()

		store := NewScheduleStore(mockLog, tx)

		ctx := context.Background()

		// given
		staffObj := staff.StaffEntity{UUID: uuid.New()}
		if err := staffStore.NewStaffStore(mockLog, tx).Save(ctx, &staffObj); err != nil {
			t.Errorf("%s", err)
		}

		// pre-save schedules
		for i := 0; i < 5; i++ {
			start := time.Date(mockLog.Date().Year(), mockLog.Date().Month(), i+1, 0, 0, 0, 0, mockLog.Timezone())
			err := store.Save(ctx, &schedule.ScheduleEntity{
				StaffId:   staffObj.StaffId,
				Start:     start,
				End:       start.Add(time.Duration(8) * time.Hour),
				IsVisible: true,
			})
			if err != nil {
				t.Error(err.Error())
				break
			}
		}

		firstDayOfMonth := time.Date(mockLog.Date().Year(), mockLog.Date().Month(), 1, 0, 0, 0, 0, mockLog.Timezone())
		lastDayOfMonth := firstDayOfMonth.AddDate(0, 1, -1)

		// method to test
		schs, err := store.SchedulesWithinTimeframe(ctx, staffObj.StaffId, firstDayOfMonth, lastDayOfMonth, true, 5*60*60)

		// assert
		if err != nil {
			t.Error(err.Error())
		}

		if len(schs) != 5 {
			t.Errorf("expect 5 given %v", len(schs))
			t.Errorf("%v", schs)
		}
	})
}
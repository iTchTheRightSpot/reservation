package service_type

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/iTchTheRightSpot/erp-golang/config"
	"github.com/iTchTheRightSpot/erp-golang/database"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/service_type"
	model "github.com/iTchTheRightSpot/erp-golang/pkg/models/staff"
	"github.com/iTchTheRightSpot/erp-golang/pkg/stores/staff"
	"github.com/iTchTheRightSpot/erp-golang/utils"
	"log"
	"reflect"
	"testing"
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
		if err = db.Close(); err != nil {
			log.Printf("db connection did not close after tests")
			return
		}
	}(d)

	// run tests
	m.Run()
}

func setupTest(t *testing.T) (*sql.Tx, func()) {
	tx, err := db.Begin()
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

func TestServiceStore(t *testing.T) {
	t.Parallel()

	logger := utils.NewDevLogger()

	t.Run("should save and return service by name", func(t *testing.T) {
		t.Parallel()

		tx, fn := setupTest(t)
		defer fn()

		store := NewServiceTypeStore(logger, tx)
		ctx := context.Background()

		// given
		s := service_type.ServiceTypeEntity{
			Name:        "name",
			Price:       1000.95,
			Duration:    3600,
			CleanUpTime: 1800,
		}

		// method to test
		if err := store.Save(ctx, &s); err != nil {
			t.Error(err)
		}

		// assert
		if s.ServiceId < 1 {
			t.Errorf("expect id to be >= 1 given %v", s.ServiceId)
		}

		find, err := store.ServiceTypeByName(ctx, "name")
		if err != nil {
			t.Error(err)
		}

		// assert
		if !reflect.DeepEqual(s, *find) {
			t.Errorf("expect %v to equal %v", s, find)
		}
	})

	t.Run("should save service and return by services staff id", func(t *testing.T) {
		t.Parallel()

		tx, fn := setupTest(t)
		defer fn()
		ctx := context.Background()

		store := NewServiceTypeStore(logger, tx)
		staffStore := staff.NewStaffStore(logger, tx)
		ss := staff.NewStaffServiceStore(logger, tx)

		// given
		s := service_type.ServiceTypeEntity{
			Name:        "name",
			Price:       1000.95,
			Duration:    3600,
			CleanUpTime: 1800,
		}

		if err := store.Save(ctx, &s); err != nil {
			t.Error(err.Error())
		}

		se := model.StaffEntity{UUID: uuid.New()}
		if err := staffStore.Save(ctx, &se); err != nil {
			t.Error(err.Error())
		}

		if err := ss.Save(ctx, &model.StaffServiceEntity{ServiceId: s.ServiceId, StaffId: se.StaffId}); err != nil {
			t.Error(err.Error())
		}

		// method to test
		arr, err := store.ServiceTypesByStaffId(ctx, se.StaffId, false)

		// assert
		if err != nil {
			t.Error(err.Error())
		}

		if len(arr) == 0 || !reflect.DeepEqual(*arr[0], s) {
			t.Errorf("expect %v to equal given %v", s, *arr[0])
		}
	})

	t.Run("reject saving. price is greater than DECIMAL(6, 2)", func(t *testing.T) {
		t.Parallel()

		tx, fn := setupTest(t)
		defer fn()

		store := NewServiceTypeStore(logger, tx)
		ctx := context.Background()

		// given
		s := service_type.ServiceTypeEntity{
			Name:  "name",
			Price: 10001.95,
		}

		// method to test
		if err := store.Save(ctx, &s); err == nil {
			t.Error(err)
		}

		// method to test
		s.Name = "higher"
		s.Price = 1000.959
		if err := store.Save(ctx, &s); err == nil {
			t.Error(err)
		}
	})
}

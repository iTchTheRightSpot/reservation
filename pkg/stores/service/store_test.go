package service

import (
	"context"
	"database/sql"
	"github.com/iTchTheRightSpot/erp-golang/config"
	"github.com/iTchTheRightSpot/erp-golang/database"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/service"
	"github.com/iTchTheRightSpot/erp-golang/utils"
	"log"
	"reflect"
	"testing"
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

func TestServiceStore(t *testing.T) {
	t.Parallel()

	logger := utils.NewMockLogger()

	t.Run("should save and return service by name", func(t *testing.T) {
		t.Parallel()

		tx, fn := setupTest(t)
		defer fn()

		store := NewServiceStore(logger, tx)
		ctx := context.Background()

		// given
		s := service.Service{
			Name:        "name",
			Price:       1000.95,
			Duration:    3600,
			CleanUpTime: 1800,
		}

		// method to test
		save, err := store.Save(ctx, &s)
		if err != nil {
			t.Error(err)
		}

		// assert
		if !reflect.DeepEqual(s, *save) {
			t.Errorf("expect %v to equal %v", s, &save)
		}

		find, err := store.ServiceByName(ctx, "name")
		if err != nil {
			t.Error(err)
		}

		// assert
		if !reflect.DeepEqual(*save, *find) {
			t.Errorf("expect %v to equal %v", save, find)
		}
	})

	t.Run("reject saving. price is greater than DECIMAL(6, 2)", func(t *testing.T) {
		t.Parallel()

		tx, fn := setupTest(t)
		defer fn()

		store := NewServiceStore(logger, tx)
		ctx := context.Background()

		// given
		s := service.Service{
			Name:  "name",
			Price: 10001.95,
		}

		// method to test
		if _, err := store.Save(ctx, &s); err == nil {
			t.Error(err)
		}

		// method to test
		s.Name = "higher"
		s.Price = 1000.959
		if _, err := store.Save(ctx, &s); err == nil {
			t.Error(err)
		}
	})
}

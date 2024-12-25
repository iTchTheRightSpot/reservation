package stores

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/iTchTheRightSpot/erp-golang/config"
	"github.com/iTchTheRightSpot/erp-golang/database"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/staff"
	"github.com/iTchTheRightSpot/erp-golang/utils"
	"log"
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
		if err := db.Close(); err != nil {
			log.Printf("db connection did not close after tests")
			return
		}
	}(db)

	// run tests
	m.Run()
}

func TestTransactionProvider(t *testing.T) {
	logger := utils.NewMockLogger()
	del := func(query string, args ...interface{}) (sql.Result, error) {
		exec, err := db.Exec(query, args...)
		if err != nil {
			return nil, err
		}
		return exec, nil
	}

	t.Run("should commit transaction", func(t *testing.T) {
		adap := NewAdapters(logger, db, NewTransactionProvider(logger, db))
		var staffId uint64
		// method to test
		ctx := context.Background()
		err := adap.Transaction.RunInTransaction(func(adapters *Adapters) error {
			save, err := adapters.StaffStore.Save(ctx, &staff.Staff{UUID: uuid.UUID{}})
			if err != nil {
				return err
			}
			staffId = save.StaffId
			return nil
		})

		defer func(id uint64) {
			if _, err = del("DELETE FROM staff WHERE staff_id = $1", id); err != nil {
				t.Errorf("%s", err)
			}
		}(staffId)

		// assert
		if err != nil {
			t.Errorf("transaction should commit successfully but go err: %s", err)
		}
	})

	t.Run("should rollback transaction", func(t *testing.T) {
		adap := NewAdapters(logger, db, NewTransactionProvider(logger, db))
		var staffId uint64
		// method to test
		ctx := context.Background()
		err := adap.Transaction.RunInTransaction(func(adapters *Adapters) error {
			uu := uuid.UUID{}
			save, err := adapters.StaffStore.Save(ctx, &staff.Staff{UUID: uu})
			if err != nil {
				return err
			}
			staffId = save.StaffId

			// save role with profile id that does not exist
			_, err = adapters.RoleStore.Save(ctx, &models.Role{Role: models.STAFF, ProfileId: 0})
			if err != nil {
				return err
			}

			return nil
		})

		// assert
		if err == nil {
			t.Error("transaction should rollback but commited")
		}

		exec, err := del("DELETE FROM staff WHERE staff_id = $1", staffId)
		if err != nil {
			t.Errorf("%s", err)
		}

		rowsAffected, err := exec.RowsAffected()
		if err != nil {
			t.Errorf("%s", err)
		} else if rowsAffected > 0 {
			t.Errorf("transaction was not rolled back")
		}
	})
}

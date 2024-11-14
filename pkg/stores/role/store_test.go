package role

import (
	"context"
	"database/sql"
	"github.com/iTchTheRightSpot/erp-golang/config"
	"github.com/iTchTheRightSpot/erp-golang/pkg/database"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/profile"
	profileStore "github.com/iTchTheRightSpot/erp-golang/pkg/stores/profile"
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

func TestRoleStore(t *testing.T) {
	mockLog := utils.NewMockLogger()

	t.Run("should save role", func(t *testing.T) {
		con, fn := setupTest(t)
		defer fn()

		ctx := context.Background()
		repo := profileStore.NewProfileStore(mockLog, con)

		// given
		p := profile.Profile{
			Firstname: "frog",
			Lastname:  "lastname",
			Email:     "frog@email.com",
		}

		if _, err := repo.Save(ctx, &p); err != nil {
			t.Errorf("%s", err)
		}

		// method to test
		r := models.Role{Role: models.STAFF, ProfileId: p.ProfileId}

		save, err := NewRoleStore(mockLog, con).Save(ctx, &r)
		if err != nil {
			t.Errorf("%s", err)
		}

		if save.RoleId < 1 {
			t.Errorf("role not saved. Expected RoleId > 0, got %d", save.RoleId)
		}

		if !reflect.DeepEqual(&r, save) {
			t.Errorf("role not saved correctly. expected: %+v, Got: %+v", r, save)
		}
	})
}

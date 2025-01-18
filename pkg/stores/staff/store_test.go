package staff

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/iTchTheRightSpot/erp-golang/config"
	"github.com/iTchTheRightSpot/erp-golang/database"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/profile"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/staff"
	profileStore "github.com/iTchTheRightSpot/erp-golang/pkg/stores/profile"
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

func TestStaffStore(t *testing.T) {
	t.Parallel()

	mockLog := utils.NewMockLogger()

	t.Run("should save staff and find by uuid", func(t *testing.T) {
		t.Parallel()

		con, fn := setupTest(t)
		defer fn()

		ctx := context.Background()
		staffRepo := NewStaffStore(mockLog, con)
		profileRepo := profileStore.NewProfileStore(mockLog, con)

		// given
		p := profile.Profile{
			Firstname: "frog",
			Lastname:  "lastname",
			Email:     "frog@email.com",
		}

		if err := profileRepo.Save(ctx, &p); err != nil {
			t.Errorf("%s", err)
		}

		s := staff.Staff{
			UUID:      uuid.New(),
			ProfileId: &p.ProfileId,
		}

		// method to test
		if err := staffRepo.Save(ctx, &s); err != nil {
			t.Error(err.Error())
		}

		if s.StaffId < 1 {
			t.Errorf("staff not saved. Expected StaffId > 0, got %d", s.StaffId)
		}

		// method to test
		if _, err := staffRepo.StaffByUUID(ctx, uuid.New().String()); err == nil {
			t.Errorf("should not find staff that does not exist")
		}

		find, err := staffRepo.StaffByUUID(ctx, s.UUID.String())
		if err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(s, *find) {
			t.Errorf("expected: %+v, Got: %+v", s, find)
		}
	})
}

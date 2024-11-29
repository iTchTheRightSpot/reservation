package profile

import (
	"context"
	"database/sql"
	"github.com/iTchTheRightSpot/erp-golang/config"
	"github.com/iTchTheRightSpot/erp-golang/database"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/profile"
	"github.com/iTchTheRightSpot/erp-golang/utils"
	"log"
	"reflect"
	"testing"
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

func TestProfileStore(t *testing.T) {
	mockLog := utils.NewMockLogger()

	t.Run("should Save profile when image_key is not null", func(t *testing.T) {
		con, fn := setupTest(t)
		defer fn()

		repo := NewProfileStore(mockLog, con)

		// given
		key := "image-key"
		p := profile.Profile{
			Firstname: "frog",
			Lastname:  "lastname",
			Email:     "frog@email.com",
			ImageKey:  &key,
		}

		// method to test
		save, err := repo.Save(context.Background(), &p)
		if err != nil {
			t.Errorf("%s", err)
		}

		if save.ProfileId < 1 {
			t.Errorf("profile not saved. Expected ProfileId > 0, got %d", save.ProfileId)
		}

		if !reflect.DeepEqual(&p, save) {
			t.Errorf("staff not saved correctly. expected: %+v, Got: %+v", p, save)
		}
	})

	t.Run("should Save profile when image_key is null", func(t *testing.T) {
		con, fn := setupTest(t)
		defer fn()

		repo := NewProfileStore(mockLog, con)

		// given
		p := profile.Profile{
			Firstname: "frog",
			Lastname:  "lastname",
			Email:     "frog@email.com",
		}

		// method to test
		save, err := repo.Save(context.Background(), &p)
		if err != nil {
			t.Errorf("%s", err)
		}

		if save.ProfileId < 1 {
			t.Errorf("profile not saved. Expected ProfileId > 0, got %d", save.ProfileId)
		}

		if !reflect.DeepEqual(&p, save) {
			t.Errorf("staff not saved correctly. expected: %+v, Got: %+v", p, save)
		}
	})

	t.Run("test saving profile, role & permissions", func(t *testing.T) {
		con, fn := setupTest(t)
		defer fn()

		ctx := context.Background()
		profileRepo := NewProfileStore(mockLog, con)
		roleRepo := NewRoleStore(mockLog, con)
		permissionRepo := NewPermissionStore(mockLog, con)

		p := profile.Profile{
			Firstname: "frog",
			Lastname:  "lastname",
			Email:     "frog@email.com",
		}

		// save profile
		if _, err := profileRepo.Save(ctx, &p); err != nil {
			t.Errorf("%s", err)
		}

		// save role & assert it is saved
		r := models.Role{Role: models.STAFF, ProfileId: p.ProfileId}

		save, err := roleRepo.Save(ctx, &r)
		if err != nil {
			t.Errorf("%s", err)
		}

		if save.RoleId < 1 {
			t.Errorf("role not saved. Expected RoleId > 0, got %d", save.RoleId)
		}

		if !reflect.DeepEqual(&r, save) {
			t.Errorf("role not saved correctly. expected: %+v, Got: %+v", r, save)
		}

		// save permission & assert
		per := models.Permission{Permission: models.WRITE, RoleId: r.RoleId}

		savePer, err := permissionRepo.Save(ctx, &per)
		if err != nil {
			t.Errorf("%s", err)
		}

		if savePer.PermissionId < 1 {
			t.Errorf("permission not saved. Expected PermissionId > 0, got %d", savePer.PermissionId)
		}

		if !reflect.DeepEqual(&per, savePer) {
			t.Errorf("permission not saved correctly. expected: %+v, Got: %+v", r, save)
		}
	})
}

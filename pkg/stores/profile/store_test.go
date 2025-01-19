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
			Password:  []byte("password"),
		}

		// method to test
		if err := repo.Save(context.Background(), &p); err != nil {
			t.Errorf("%s", err)
		}

		if p.ProfileId < 1 {
			t.Errorf("profile not saved. Expected ProfileId > 0, got %d", p.ProfileId)
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
			Password:  []byte("password"),
			Email:     "frog@email.com",
		}

		// method to test
		if err := repo.Save(context.Background(), &p); err != nil {
			t.Error(err.Error())
		}

		if p.ProfileId < 1 {
			t.Errorf("profile not saved. Expected ProfileId > 0, got %d", p.ProfileId)
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
			Password:  []byte("password"),
			Firstname: "frog",
			Lastname:  "lastname",
			Email:     "frog@email.com",
		}

		// save profile
		if err := profileRepo.Save(ctx, &p); err != nil {
			t.Errorf("%s", err)
		}

		// save role & assert it is saved
		r := models.Role{Role: models.STAFF, ProfileId: p.ProfileId}

		if err := roleRepo.Save(ctx, &r); err != nil {
			t.Errorf("%s", err)
		}

		if r.RoleId < 1 {
			t.Errorf("role not saved. Expected RoleId > 0, got %d", r.RoleId)
		}

		// save permission & assert
		per := models.Permission{Permission: models.WRITE, RoleId: r.RoleId}

		if err := permissionRepo.Save(ctx, &per); err != nil {
			t.Errorf("%s", err)
		}

		if per.PermissionId < 1 {
			t.Errorf("permission not saved. Expected PermissionId > 0, got %d", per.PermissionId)
		}
	})
}

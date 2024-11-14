package profile

import (
	"context"
	"database/sql"
	"github.com/iTchTheRightSpot/erp-golang/config"
	"github.com/iTchTheRightSpot/erp-golang/pkg/database"
	"github.com/iTchTheRightSpot/erp-golang/pkg/model/profile"
	"github.com/iTchTheRightSpot/erp-golang/utils"
	"log"
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

func TestProfile(t *testing.T) {
	mockLog := utils.NewMockLogger()

	t.Run("should save profile when image_key is not null", func(t *testing.T) {
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
		save, err := repo.save(context.Background(), &p)
		if err != nil {
			t.Errorf("%s", err)
		}

		if save.ProfileId < 1 {
			t.Errorf("profile not saved. Expected ProfileId > 0, got %d", save.ProfileId)
		}

		if save.Firstname != p.Firstname {
			t.Errorf("profile not saved. Expected Firstname '%s', got '%s'", p.Firstname, save.Firstname)
		}

		if save.Lastname != p.Lastname {
			t.Errorf("profile not saved. Expected Lastname '%s', got '%s'", p.Lastname, save.Lastname)
		}

		if save.Email != p.Email {
			t.Errorf("profile not saved. Expected Email '%s', got '%s'", p.Email, save.Email)
		}

		if save.ImageKey != p.ImageKey {
			t.Errorf("profile not saved. Expected ImageKey to be nil, got '%v'", save.ImageKey)
		}
	})

	t.Run("should save profile when image_key is null", func(t *testing.T) {
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
		save, err := repo.save(context.Background(), &p)
		if err != nil {
			t.Errorf("%s", err)
		}

		if save.ProfileId < 1 {
			t.Errorf("profile not saved. Expected ProfileId > 0, got %d", save.ProfileId)
		}

		if save.Firstname != p.Firstname {
			t.Errorf("profile not saved. Expected Firstname '%s', got '%s'", p.Firstname, save.Firstname)
		}

		if save.Lastname != p.Lastname {
			t.Errorf("profile not saved. Expected Lastname '%s', got '%s'", p.Lastname, save.Lastname)
		}

		if save.Email != p.Email {
			t.Errorf("profile not saved. Expected Email '%s', got '%s'", p.Email, save.Email)
		}

		if save.ImageKey != nil {
			t.Errorf("profile not saved. Expected ImageKey to be nil, got '%v'", save.ImageKey)
		}
	})
}

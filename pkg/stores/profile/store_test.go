package profile

import (
	"context"
	"database/sql"
	"github.com/iTchTheRightSpot/erp-golang/config"
	"github.com/iTchTheRightSpot/erp-golang/database"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models"
	"github.com/iTchTheRightSpot/utility/utils"
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
	mockLog := utils.DevLogger("UTC")

	t.Run("should Save profile when image_key is not null", func(t *testing.T) {
		con, fn := setupTest(t)
		defer fn()

		repo := NewProfileStore(mockLog, con)

		// given
		key := "image-key"
		p := models.ProfileEntity{
			Firstname: "frog",
			Lastname:  "lastname",
			Email:     "frog@email.com",
			ImageKey:  &key,
			Password:  "password",
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
		p := models.ProfileEntity{
			Firstname: "frog",
			Lastname:  "lastname",
			Password:  "password",
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

	t.Run("account flow", func(t *testing.T) {
		t.Run("saving profile, role & permissions", func(t *testing.T) {
			con, fn := setupTest(t)
			defer fn()

			ctx := context.Background()
			profileRepo := NewProfileStore(mockLog, con)
			roleRepo := NewRoleStore(mockLog, con)
			permissionRepo := NewPermissionStore(mockLog, con)

			p := models.ProfileEntity{
				Password:  "password",
				Firstname: "frog",
				Lastname:  "lastname",
				Email:     "frog@email.com",
			}

			// save profile
			if err := profileRepo.Save(ctx, &p); err != nil {
				t.Errorf("%s", err)
			}

			// save role & assert it is saved
			r := models.RoleEntity{Role: models.STAFF, ProfileId: p.ProfileId}

			if err := roleRepo.Save(ctx, &r); err != nil {
				t.Errorf("%s", err)
			}

			if r.RoleId < 1 {
				t.Errorf("role not saved. Expected RoleId > 0, got %d", r.RoleId)
			}

			// save permission & assert
			per := models.PermissionEntity{Permission: models.WRITE, RoleId: r.RoleId}

			if err := permissionRepo.Save(ctx, &per); err != nil {
				t.Errorf("%s", err)
			}

			if per.PermissionId < 1 {
				t.Errorf("permission not saved. Expected PermissionId > 0, got %d", per.PermissionId)
			}
		})

		t.Run("ProfileRolesAndPermissionByEmail", func(t *testing.T) {
			tx, fn := setupTest(t)
			defer fn()

			ctx := context.Background()
			profileRepo := NewProfileStore(mockLog, tx)
			roleRepo := NewRoleStore(mockLog, tx)
			permissionRepo := NewPermissionStore(mockLog, tx)

			// save profile
			p := models.ProfileEntity{
				Password:  "password",
				Firstname: "frog",
				Lastname:  "lastname",
				Email:     "frog@email.com",
			}
			_ = profileRepo.Save(ctx, &p)

			// save role1 & permission1
			r1 := models.RoleEntity{Role: models.STAFF, ProfileId: p.ProfileId}
			_ = roleRepo.Save(ctx, &r1)
			p1 := models.PermissionEntity{Permission: models.WRITE, RoleId: r1.RoleId}
			_ = permissionRepo.Save(ctx, &p1)

			// save role2 & permission2
			r2 := models.RoleEntity{Role: models.DEVELOPER, ProfileId: p.ProfileId}
			_ = roleRepo.Save(ctx, &r2)
			p2 := models.PermissionEntity{Permission: models.READ, RoleId: r2.RoleId}
			_ = permissionRepo.Save(ctx, &p2)

			// method to test & assert
			prsSave, err := profileRepo.ProfileRolesAndPermissionByEmail(ctx, p.Email)
			if err != nil {
				t.Error(err.Error())
			}

			prs := models.ProfileRolePermissionEntity{
				Profile: p,
				RolePermission: []models.RolePermissionEntity{
					{
						Role:        r1,
						Permissions: []models.PermissionEntity{p1},
					},
					{
						Role:        r2,
						Permissions: []models.PermissionEntity{p2},
					},
				},
			}

			if !reflect.DeepEqual(prs, *prsSave) {
				if prs.Profile != prsSave.Profile {
					t.Errorf("Profile mismatch: %+v != %+v", prs.Profile, prsSave.Profile)
				}

				for i, rp := range prs.RolePermission {
					if i >= len(prsSave.RolePermission) {
						t.Errorf("Extra role in expected: %+v", rp)
						continue
					}
					if !reflect.DeepEqual(rp, prsSave.RolePermission[i]) {
						t.Errorf("RolePermission[%d] mismatch: %+v != %+v", i, rp, prsSave.RolePermission[i])
					}
				}
				if len(prsSave.RolePermission) > len(prs.RolePermission) {
					t.Errorf("Extra roles in actual: %+v", prsSave.RolePermission[len(prs.RolePermission):])
				}
			}
		})
	})
}
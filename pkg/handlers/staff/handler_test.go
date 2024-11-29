package staff

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/iTchTheRightSpot/erp-golang/config"
	"github.com/iTchTheRightSpot/erp-golang/database"
	"github.com/iTchTheRightSpot/erp-golang/pkg/middleware"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/profile"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/service"
	model "github.com/iTchTheRightSpot/erp-golang/pkg/models/staff"
	"github.com/iTchTheRightSpot/erp-golang/pkg/services/auth"
	"github.com/iTchTheRightSpot/erp-golang/pkg/services/staff"
	"github.com/iTchTheRightSpot/erp-golang/pkg/stores"
	"github.com/iTchTheRightSpot/erp-golang/utils"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var db *sql.DB
var env *config.SecretVariables

func TestMain(m *testing.M) {
	if err := os.Chdir("../../../"); err != nil {
		log.Fatalf("failed to change directory: %v", err)
	}

	secret := config.SecretVariables{}
	e, err := secret.Config()
	if err != nil {
		log.Fatal(err)
	}
	env = e

	db, err = database.ConnectToPostgre(e.DbConnectionString)
	if err != nil {
		log.Fatal(err)
	}

	defer func(db *sql.DB) {
		if err := db.Close(); err != nil {
			log.Printf("db connection did not close after tests")
			return
		}
	}(db)

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

func preSaveStaff(a *stores.Adapters) (*model.Staff, error) {
	ctx := context.Background()

	p := profile.Profile{
		Firstname: "erp",
		Lastname:  "erp",
		Email:     "erp@email.com",
	}

	if _, err := a.ProfileStore.Save(ctx, &p); err != nil {
		return nil, err
	}

	s := model.Staff{
		UUID:      uuid.New(),
		ProfileId: &p.ProfileId,
	}

	if _, err := a.StaffStore.Save(ctx, &s); err != nil {
		return nil, err
	}

	return &s, nil
}

func preSaveService(a *stores.Adapters) (*service.ServiceEntity, error) {
	return a.ServiceStore.Save(context.Background(), &service.ServiceEntity{
		Name:        "erp",
		Price:       19.56,
		Duration:    3600,
		CleanUpTime: 30 * 60,
	})
}

func TestStaffHandler(t *testing.T) {
	t.Parallel()

	logger := utils.NewMockLogger()

	t.Run("should link service to staff & also reject", func(t *testing.T) {
		t.Parallel()

		tx, fn := setupTest(t)
		defer fn()

		// given
		mux := http.NewServeMux()
		prov := stores.MockLiveTransactionProvider(logger, tx)
		adapters := stores.NewAdapters(logger, tx, prov)
		jwtSer := auth.NewJwtService(logger, env)
		m := &middleware.Middleware{Logger: logger, Auth: jwtSer, Param: env.CookieParam}
		s := staff.NewStaffService(logger, adapters)

		saveStaff, err := preSaveStaff(adapters)
		if err != nil {
			t.Error(err)
		}

		saveService, err := preSaveService(adapters)
		if err != nil {
			t.Error(err)
		}

		cred := []models.RolePermission{
			{
				Role:        models.STAFF,
				Permissions: []models.PermissionEnum{models.WRITE},
			},
		}

		obj, err := jwtSer.GenerateJwt(
			&models.JwtObj{
				UserId:         saveStaff.UUID.String(),
				AccessControls: cred,
			},
			utils.TwoDaysInSeconds,
		)

		NewStaffHandler(mux, m, logger, s).Register()

		for i := 0; i < 2; i++ {
			url := fmt.Sprintf("/staff/service?service_name=%s", saveService.Name)
			req := httptest.NewRequest(http.MethodPost, url, nil)
			req.AddCookie(&http.Cookie{Name: env.CookieParam.CookieName, Value: obj.Token})
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			mux.ServeHTTP(rr, req)

			// assert
			if i == 0 {
				if rr.Code != http.StatusCreated {
					t.Errorf("index %v expected status code %d, got %d", i, http.StatusCreated, rr.Code)
				}
			} else {
				if rr.Code != http.StatusBadRequest {
					t.Errorf("index %v expected status code %d, got %d", i, http.StatusBadRequest, rr.Code)
				}
			}
		}
	})
}

package staff

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/iTchTheRightSpot/erp-golang/config"
	"github.com/iTchTheRightSpot/erp-golang/database"
	"github.com/iTchTheRightSpot/erp-golang/pkg/handlers"
	"github.com/iTchTheRightSpot/erp-golang/pkg/middleware"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models"
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
	env = secret.Config()

	d, err := database.ConnectToPostgres(env.DbConnectionString)
	if err != nil {
		log.Fatal(err)
	}
	db = d

	defer func(db *sql.DB) {
		if err = db.Close(); err != nil {
			log.Printf("db connection did not close after tests")
			return
		}
	}(db)

	m.Run()
}

func TestStaffHandler(t *testing.T) {
	t.Cleanup(func() {
		if err := handlers.DeleteAll(db); err != nil {
			t.Errorf(err.Error())
		}
	})

	mux := http.NewServeMux()
	logger := utils.NewMockLogger()
	prov := stores.NewTransactionProvider(logger, db)
	adp := stores.NewAdapters(logger, db, prov)
	jwtSer := auth.NewJwtService(logger, env)
	m := &middleware.Middleware{Logger: logger, Auth: jwtSer, Param: env.CookieParam}
	s := staff.NewStaffService(logger, adp)

	// register routes
	NewStaffHandler(mux, m, logger, s).Register()

	// pre-save
	saveStaff, err := preSaveStaff(adp)
	if err != nil {
		t.Error(err)
	}

	saveService, err := preSaveService(adp)
	if err != nil {
		t.Error(err)
	}

	t.Run("should link service to staff & also reject", func(t *testing.T) {
		cred := []models.RolePermissionEnum{
			{
				Role:        models.STAFF,
				Permissions: []models.PermissionEnum{models.WRITE},
			},
		}

		obj, _ := jwtSer.Encode(
			&models.JwtObj{
				UserId:         saveStaff.UUID.String(),
				AccessControls: cred,
			},
			utils.TwoDaysInSeconds,
		)

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
				if rr.Code != http.StatusConflict {
					t.Errorf("index %v expected status code %d, got %d", i, http.StatusBadRequest, rr.Code)
				}
			}
		}
	})
}

func preSaveStaff(a *stores.Adapters) (*model.Staff, error) {
	ctx := context.Background()

	p := models.ProfileEntity{
		Firstname: "erp",
		Lastname:  "erp",
		Email:     "erp@email.com",
		Password:  "password",
	}

	if err := a.ProfileStore.Save(ctx, &p); err != nil {
		return nil, err
	}

	s := model.Staff{
		UUID:      uuid.New(),
		ProfileId: &p.ProfileId,
	}

	if err := a.StaffStore.Save(ctx, &s); err != nil {
		return nil, err
	}

	return &s, nil
}

func preSaveService(a *stores.Adapters) (*service.ServiceTypeEntity, error) {
	s := service.ServiceTypeEntity{
		Name:        "erp",
		Price:       19.56,
		Duration:    3600,
		CleanUpTime: 30 * 60,
	}
	err := a.ServiceStore.Save(context.Background(), &s)
	return &s, err
}

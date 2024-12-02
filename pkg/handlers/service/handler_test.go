package service

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"github.com/iTchTheRightSpot/erp-golang/config"
	"github.com/iTchTheRightSpot/erp-golang/database"
	"github.com/iTchTheRightSpot/erp-golang/pkg/middleware"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models"
	model "github.com/iTchTheRightSpot/erp-golang/pkg/models/service"
	"github.com/iTchTheRightSpot/erp-golang/pkg/services/auth"
	"github.com/iTchTheRightSpot/erp-golang/pkg/services/service"
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

	db, err = database.ConnectToPostgres(e.DbConnectionString)
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

func TestServiceHandler(t *testing.T) {
	t.Parallel()

	logger := utils.NewMockLogger()

	t.Run("should create service", func(t *testing.T) {
		tx, fn := setupTest(t)
		defer fn()

		// given
		mux := http.NewServeMux()
		prov := stores.MockLiveTransactionProvider(logger, tx)
		adapters := stores.NewAdapters(logger, tx, prov)
		jwtSer := auth.NewJwtService(logger, env)
		m := &middleware.Middleware{Logger: logger, Auth: jwtSer, Param: env.CookieParam}
		s := service.NewServiceImpl(logger, adapters)

		arr := []models.RolePermission{
			{
				Role:        models.STAFF,
				Permissions: []models.PermissionEnum{models.WRITE},
			},
		}

		obj, err := jwtSer.GenerateJwt(
			&models.JwtObj{
				UserId:         "staff-uuid",
				AccessControls: arr,
			},
			utils.TwoDaysInSeconds,
		)

		dtoBytes, err := json.Marshal(model.ServicePayload{
			Name:        "erp",
			Price:       15.97,
			Duration:    3600,
			CleanUpTime: 1800,
		})
		if err != nil {
			t.Errorf("failed to marshal SchedulePayload: %s", err)
		}

		req := httptest.NewRequest(http.MethodPost, "/service", bytes.NewBuffer(dtoBytes))
		req.AddCookie(&http.Cookie{Name: env.CookieParam.CookieName, Value: obj.Token})
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		// handler to test
		NewServiceHandler(mux, logger, s, m).Register()

		mux.ServeHTTP(rr, req)

		// assert
		if rr.Code != http.StatusCreated {
			t.Errorf("expected status code %d, got %d", http.StatusCreated, rr.Code)
		}
	})
}

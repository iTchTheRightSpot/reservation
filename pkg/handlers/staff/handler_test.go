package staff

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/iTchTheRightSpot/erp-golang/config"
	"github.com/iTchTheRightSpot/erp-golang/database"
	"github.com/iTchTheRightSpot/erp-golang/pkg/handlers"
	"github.com/iTchTheRightSpot/erp-golang/pkg/middleware"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models"
	model "github.com/iTchTheRightSpot/erp-golang/pkg/models/staff"
	pkg "github.com/iTchTheRightSpot/erp-golang/pkg/services"
	"github.com/iTchTheRightSpot/erp-golang/pkg/services/auth"
	"github.com/iTchTheRightSpot/erp-golang/pkg/services/staff"
	"github.com/iTchTheRightSpot/erp-golang/pkg/stores"
	"github.com/iTchTheRightSpot/erp-golang/utils"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"slices"
	"strings"
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
	jwtSer := auth.NewJwtServiceAsymmetric(logger, env)
	m := &middleware.Middleware{Logger: logger, Auth: jwtSer, Param: env.CookieParam}
	cach := pkg.NewInMemoryCache[string, []*model.AllStaffsEntity](logger, 10, 10)
	s := staff.NewStaffService(logger, adp, cach)

	// register routes
	NewStaffHandler(mux, m, logger, s).Register()

	t.Run("should return all staffs", func(t *testing.T) {
		obj, _ := jwtSer.Encode(
			&models.JwtObj{
				UserId: "uuid",
				AccessControls: []models.RolePermissionEnum{
					{
						Role:        models.STAFF,
						Permissions: []models.PermissionEnum{models.WRITE},
					},
				},
			},
			utils.TwoDaysInSeconds,
		)

		t.Run("empty array", func(t *testing.T) {
			// route to test
			req := httptest.NewRequest(http.MethodGet, "/staffs", nil)
			req.AddCookie(&http.Cookie{Name: env.CookieParam.CookieName, Value: obj.Token})
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			mux.ServeHTTP(rr, req)

			// assert
			if rr.Code != http.StatusOK {
				t.Errorf("expected status code %d, received %d", http.StatusOK, rr.Code)
				t.Error(rr.Body.String())
			}

			if strings.TrimSpace(rr.Body.String()) != "[]" {
				t.Errorf("expect [], given %s", rr.Body.String())
			}
		})

		t.Run("not empty array", func(t *testing.T) {
			size := 10

			// pre-save
			for i := 0; i < size; i++ {
				_, err := preSaveStaff(context.Background(), adp)
				if err != nil {
					t.Error(err.Error())
				}
			}

			// route to test
			req := httptest.NewRequest(http.MethodGet, "/staffs", nil)
			req.AddCookie(&http.Cookie{Name: env.CookieParam.CookieName, Value: obj.Token})
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			mux.ServeHTTP(rr, req)

			// assert
			if rr.Code != http.StatusOK {
				t.Errorf("expected status code %d, received %d", http.StatusOK, rr.Code)
				t.Error(rr.Body.String())
			}

			var arr []model.AllStaffsEntity

			if err := json.NewDecoder(rr.Body).Decode(&arr); err != nil {
				t.Error(err.Error())
			}

			if len(arr) != size {
				t.Errorf("expect size %d given %v", size, len(arr))
			}

			var temp = make([]model.AllStaffsEntity, len(arr))
			for i, stf := range arr {
				b := slices.ContainsFunc(temp, func(entity model.AllStaffsEntity) bool {
					return stf.Email == entity.Email
				})

				if b {
					t.Error(" array contains duplicate")
					break
				}
				temp[i] = stf
			}
		})
	})
}

func preSaveStaff(ctx context.Context, a *stores.Adapters) (*model.StaffEntity, error) {
	id := uuid.New()
	p := models.ProfileEntity{
		Firstname: "firstname",
		Lastname:  "lastname",
		Email:     fmt.Sprintf("%s@email.com", id.String()),
		Password:  "password",
	}

	if err := a.ProfileStore.Save(ctx, &p); err != nil {
		return nil, err
	}

	r := models.RoleEntity{Role: models.STAFF, ProfileId: p.ProfileId}
	if err := a.RoleStore.Save(ctx, &r); err != nil {
		return nil, err
	}

	if err := a.PermissionStore.Save(ctx, &models.PermissionEntity{RoleId: r.RoleId, Permission: models.WRITE}); err != nil {
		return nil, err
	}

	s := model.StaffEntity{
		UUID:      id,
		ProfileId: &p.ProfileId,
	}

	if err := a.StaffStore.Save(ctx, &s); err != nil {
		return nil, err
	}

	return &s, nil
}

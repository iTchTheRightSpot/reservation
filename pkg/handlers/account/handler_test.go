package account

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/iTchTheRightSpot/reservation/config"
	"github.com/iTchTheRightSpot/reservation/database"
	"github.com/iTchTheRightSpot/reservation/pkg/handlers"
	"github.com/iTchTheRightSpot/reservation/pkg/middleware"
	"github.com/iTchTheRightSpot/reservation/pkg/models"
	"github.com/iTchTheRightSpot/reservation/pkg/models/staff"
	"github.com/iTchTheRightSpot/reservation/pkg/services/account"
	"github.com/iTchTheRightSpot/reservation/pkg/services/auth"
	"github.com/iTchTheRightSpot/reservation/pkg/stores"
	"github.com/iTchTheRightSpot/reservation/utils"
	"github.com/iTchTheRightSpot/utility/cache"
	logg "github.com/iTchTheRightSpot/utility/utils"
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

func TestAccountHandler(t *testing.T) {
	t.Cleanup(func() {
		if err := handlers.DeleteAll(db); err != nil {
			t.Error(err.Error())
		}
	})

	// given
	mux := http.NewServeMux()
	l := logg.DevLogger("UTC")
	prov := stores.NewTransactionProvider(l, db)
	adp := stores.NewAdapters(l, db, prov)
	jwtSer := auth.NewJwtServiceAsymmetric(l, env)
	ps := auth.NewPasswordService(l)
	acs := account.NewAccountService(l, adp, jwtSer, ps, cache.SyncMapInMemoryCache[string, []*staff.AllStaffsEntity](l, 10, 10))

	// register handler
	m := &middleware.Middleware{Logger: l, Auth: jwtSer, Param: env.CookieParam}
	NewAccountHandler(mux, m, l, env, ps, acs).Register()

	// given
	name := uuid.NewString()
	p := models.ProfilePayload{
		Firstname: name,
		Lastname:  name,
		Email:     name + "@email.com",
		Password:  "pa(ssworD123#",
	}

	t.Run("should register a user", func(t *testing.T) {
		jwtObj, _ := jwtSer.Encode(
			context.Background(),
			&models.JwtObj{
				UserId: p.Firstname,
				AccessControls: []models.RolePermissionEnum{
					{Role: models.STAFF, Permissions: []models.PermissionEnum{models.WRITE}},
				},
			},
			utils.TwoDaysInSeconds,
		)

		pl, err := json.Marshal(p)
		if err != nil {
			t.Errorf("failed to marshal ProfilePayload: %s", err)
			return
		}

		req := httptest.NewRequest(http.MethodPost, "/account/register", bytes.NewBuffer(pl))
		req.AddCookie(&http.Cookie{Name: env.CookieParam.CookieName, Value: jwtObj.Token})
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, req)

		// assert
		if rr.Code != http.StatusCreated {
			t.Errorf("expected %d, received %d", http.StatusCreated, rr.Code)
		}
	})

	t.Run("should login", func(t *testing.T) {
		pl, err := json.Marshal(models.Login{Email: p.Email, Password: p.Password})
		if err != nil {
			t.Errorf("failed to marshal Login: %s", err)
			return
		}

		req := httptest.NewRequest(http.MethodPost, "/account/login", bytes.NewBuffer(pl))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, req)

		// assert
		if rr.Code != http.StatusCreated {
			t.Errorf("expected %d, received %d", http.StatusCreated, rr.Code)
		}
	})

	t.Run("should return active user", func(t *testing.T) {
		ctx := context.Background()

		pr, err := adp.ProfileStore.ProfileByEmail(ctx, p.Email)
		if err != nil {
			t.Error(err.Error())
		}

		staf := staff.StaffEntity{ProfileId: &pr.ProfileId, UUID: uuid.New()}
		if err = adp.StaffStore.Save(ctx, &staf); err != nil {
			t.Error(err.Error())
		}

		cred := []models.RolePermissionEnum{
			{Role: models.STAFF, Permissions: []models.PermissionEnum{models.READ}},
		}

		obj, err := jwtSer.Encode(context.Background(),
			&models.JwtObj{
				UserId:         staf.UUID.String(),
				AccessControls: cred,
			},
			utils.TwoDaysInSeconds,
		)

		if err != nil {
			t.Error(err.Error())
			return
		}

		req := httptest.NewRequest(http.MethodGet, "/active", nil)
		req.AddCookie(&http.Cookie{Name: env.CookieParam.CookieName, Value: obj.Token})
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected %d, received %d", http.StatusOK, rr.Code)
			t.Error(rr.Body.String())
		}

		if len(rr.Body.String()) < 1 {
			t.Error("expected not empty body, received empty body")
		}
	})

	jwtObj, _ := jwtSer.Encode(context.Background(),
		&models.JwtObj{
			UserId: p.Firstname,
			AccessControls: []models.RolePermissionEnum{
				{Role: models.STAFF, Permissions: []models.PermissionEnum{models.WRITE}},
				{Role: models.DEVELOPER, Permissions: []models.PermissionEnum{models.WRITE}},
			},
		},
		utils.TwoDaysInSeconds,
	)

	pr, err := adp.ProfileStore.ProfileRolesAndPermissionByEmail(context.Background(), p.Email)
	if err != nil {
		t.Error(err.Error())
	}

	if pr == nil {
		t.Fatal("profile is nil")
	}

	staf, err := adp.StaffStore.StaffByProfileId(context.Background(), pr.Profile.ProfileId)
	if err != nil {
		t.Error(err.Error())
	}

	t.Run("should delete permission", func(t *testing.T) {
		url := fmt.Sprintf("/account/permission/%s/%s/%s",
			staf.UUID.String(),
			pr.RolePermission[0].Role.Role,
			pr.RolePermission[0].Permissions[0].Permission)
		req := httptest.NewRequest(http.MethodDelete, url, nil)
		req.AddCookie(&http.Cookie{Name: env.CookieParam.CookieName, Value: jwtObj.Token})
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, req)

		// assert
		if rr.Code != http.StatusNoContent {
			t.Errorf("expected %d, received %d", http.StatusNoContent, rr.Code)
			t.Log(rr.Body.String())
		}
	})

	t.Run("should delete role", func(t *testing.T) {
		url := fmt.Sprintf("/account/role/%s/%s", staf.UUID.String(), pr.RolePermission[0].Role.Role)
		req := httptest.NewRequest(http.MethodDelete, url, nil)
		req.AddCookie(&http.Cookie{Name: env.CookieParam.CookieName, Value: jwtObj.Token})
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, req)

		// assert
		if rr.Code != http.StatusNoContent {
			t.Errorf("expected %d, received %d", http.StatusNoContent, rr.Code)
			t.Log(rr.Body.String())
		}
	})

	t.Run("should add role to user", func(t *testing.T) {
		st, err := adp.StaffStore.StaffByProfileId(context.Background(), pr.Profile.ProfileId)
		if err != nil {
			t.Error(err.Error())
		}

		pl, err := json.Marshal(models.RoleAndPermissionPayload{
			UserId: st.UUID.String(),
			RolePermission: []models.RolePermissionEnum{
				{Role: models.DEVELOPER, Permissions: []models.PermissionEnum{models.DELETE}},
			},
		})
		if err != nil {
			t.Errorf("failed to marshal Login: %s", err)
			return
		}

		req := httptest.NewRequest(http.MethodPost, "/account/role-permission", bytes.NewBuffer(pl))
		req.AddCookie(&http.Cookie{Name: env.CookieParam.CookieName, Value: jwtObj.Token})
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, req)

		// assert
		if rr.Code != http.StatusNoContent {
			t.Errorf("expected %d, received %d", http.StatusNoContent, rr.Code)
		}
	})
}
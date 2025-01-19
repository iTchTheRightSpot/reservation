package account

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/iTchTheRightSpot/erp-golang/config"
	"github.com/iTchTheRightSpot/erp-golang/database"
	"github.com/iTchTheRightSpot/erp-golang/pkg/handlers"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/profile"
	"github.com/iTchTheRightSpot/erp-golang/pkg/services/account"
	"github.com/iTchTheRightSpot/erp-golang/pkg/services/auth"
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
		if err := db.Close(); err != nil {
			log.Printf("db connection did not close after tests")
			return
		}
	}(db)

	m.Run()
}

func TestAccountHandler(t *testing.T) {
	t.Cleanup(func() {
		if err := handlers.DeleteAll(db); err != nil {
			t.Errorf(err.Error())
		}
	})

	// given
	mux := http.NewServeMux()
	l := utils.NewMockLogger()
	prov := stores.NewTransactionProvider(l, db)
	adp := stores.NewAdapters(l, db, prov)
	ps := auth.NewPasswordService(l)
	acs := account.NewAccountService(l, adp, ps)

	// register handler
	NewAccountHandler(mux, l, ps, acs).Register()

	t.Run("should register a user", func(t *testing.T) {
		name := uuid.NewString()
		pl, err := json.Marshal(profile.ProfilePayload{
			Firstname: name,
			Lastname:  name,
			Email:     name + "@email.com",
			Password:  "pa(ssworD123#",
		})

		if err != nil {
			t.Errorf("failed to marshal SchedulePayload: %s", err)
			return
		}

		req := httptest.NewRequest(http.MethodPost, "/account/register", bytes.NewBuffer(pl))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, req)

		// assert
		if rr.Code != http.StatusCreated {
			t.Errorf("expected %d, received %d", http.StatusCreated, rr.Code)
		}
	})
}

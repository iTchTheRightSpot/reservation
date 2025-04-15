package reservation

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
	scheduleHandler "github.com/iTchTheRightSpot/reservation/pkg/handlers/schedule"
	"github.com/iTchTheRightSpot/reservation/pkg/middleware"
	"github.com/iTchTheRightSpot/reservation/pkg/models"
	model "github.com/iTchTheRightSpot/reservation/pkg/models/reservation"
	scheduleModel "github.com/iTchTheRightSpot/reservation/pkg/models/schedule"
	"github.com/iTchTheRightSpot/reservation/pkg/models/service_type"
	"github.com/iTchTheRightSpot/reservation/pkg/models/staff"
	"github.com/iTchTheRightSpot/reservation/pkg/services/auth"
	"github.com/iTchTheRightSpot/reservation/pkg/services/mail"
	"github.com/iTchTheRightSpot/reservation/pkg/services/reservation"
	"github.com/iTchTheRightSpot/reservation/pkg/services/schedule"
	"github.com/iTchTheRightSpot/reservation/pkg/stores"
	"github.com/iTchTheRightSpot/reservation/utils"
	"github.com/iTchTheRightSpot/utility/cache"
	logg "github.com/iTchTheRightSpot/utility/utils"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
	"time"
	"unicode"
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

func TestReservationHandler(t *testing.T) {
	t.Cleanup(func() {
		if err := handlers.DeleteAll(db); err != nil {
			t.Error(err.Error())
		}
	})

	// given
	mux := http.NewServeMux()
	logger := logg.DevLogger("UTC")
	prov := stores.NewTransactionProvider(logger, db)
	adapters := stores.NewAdapters(logger, db, prov)
	jwtSer := auth.NewJwtServiceAsymmetric(logger, env)
	ware := &middleware.Middleware{Logger: logger, Auth: jwtSer, Param: env.CookieParam}
	s := schedule.NewScheduleService(logger, adapters)
	c := cache.SyncMapInMemoryCache[string, []*model.ReservationTimeSlots](logger, 30, 30)
	mailService := &mail.MockMailService{}
	rs := reservation.NewReservationService(logger, adapters, c, mailService)

	ctx := context.Background()
	staff1, err := handlers.PreSaveStaff(ctx, adapters)
	if err != nil {
		t.Error(err.Error())
	}

	serviceType1, err := preSaveService(ctx, adapters)
	if err != nil {
		t.Fatalf("preSaveService type 1 failed: %v", err.Error())
	}

	if err = linkServiceToStaff(adapters, staff1, serviceType1); err != nil {
		t.Fatal(err.Error())
	}

	serviceType2, err := preSaveService(ctx, adapters)
	if err != nil {
		t.Fatalf("preSaveService type 2 failed: %v", err.Error())
	}

	if err = linkServiceToStaff(adapters, staff1, serviceType2); err != nil {
		t.Fatal(err.Error())
	}

	// register schedule handler
	scheduleHandler.NewScheduleHandler(mux, ware, logger, s).Register()

	// register reservation handler
	NewReservationHandler(mux, logger, ware, rs).Register()

	date := logger.Date()
	d := time.Date(date.Year(), date.Month(), date.Day(), 9, 0, 0, 0, logger.Timezone()).Add(24 * time.Hour)
	zone, err := randomTimezone()
	if err != nil {
		t.Errorf("error getting timezone %s", err.Error())
	}

	t.Run("flow to creating a reservation", func(t *testing.T) {
		// create staff schedule
		createSchedule(t, jwtSer, staff1, d, mux)

		t.Run("return valid reservation dates. validate response body is not nil", func(t *testing.T) {
			url := fmt.Sprintf(
				"/reservation?day=%v&month=%v&year=%v&staff_id=%s&service=%s&timezone=%s",
				1, int(d.Month()), d.Year()+1, staff1.UUID.String(), serviceType1.Name, zone.String(),
			)
			req := httptest.NewRequest(http.MethodGet, url, nil)
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()
			mux.ServeHTTP(rr, req)

			if rr.Code != http.StatusOK {
				t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
			}

			body := strings.TrimSpace(rr.Body.String())

			if body == "" {
				t.Error("expected non-empty body, got nil or empty body")
			}

			if strings.Compare(body, "[]") != 0 {
				t.Errorf("expected empty array, got: %s", body)
			}
		})

		t.Run("success. create a reservation single service", func(t *testing.T) {
			// retrieve reservation times
			types := []string{serviceType1.Name}
			reserves := reservationTimes(t, d, staff1, types, zone.String(), mux, err)

			// create reservation
			maxT := len(reserves[0].Times)
			createBody := model.ReservationPayload{
				StaffId:  staff1.UUID.String(),
				Name:     "user-name",
				Email:    uuid.NewString() + "@email.com",
				Phone:    "0123456789",
				Services: types,
				Time:     reserves[0].Times[rand.Intn(maxT-0)+0],
				Timezone: zone.String(),
			}

			bts, err := json.Marshal(createBody)
			if err != nil {
				t.Fatalf("failed to marshal SchedulePayload: %s", err)
			}

			req := httptest.NewRequest(http.MethodPost, "/reservation", bytes.NewBuffer(bts))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()
			mux.ServeHTTP(rr, req)

			if rr.Code != http.StatusCreated {
				t.Errorf("expected status code %d, got %d", http.StatusCreated, rr.Code)
				t.Error(rr.Body.String())
			}
		})

		t.Run("success. create a reservation multiple services", func(t *testing.T) {
			// retrieve reservation times
			types := []string{serviceType1.Name, serviceType2.Name}
			reserves := reservationTimes(t, d, staff1, types, zone.String(), mux, err)

			// create reservation
			maxT := len(reserves[0].Times)
			payload := model.ReservationPayload{
				StaffId:  staff1.UUID.String(),
				Name:     "user-name",
				Email:    uuid.NewString() + "@email.com",
				Phone:    "0123456789",
				Services: types,
				Time:     reserves[0].Times[rand.Intn(maxT-0)+0],
				Timezone: zone.String(),
			}

			bts, err := json.Marshal(payload)
			if err != nil {
				t.Fatalf("failed to marshal SchedulePayload: %s", err)
			}

			req := httptest.NewRequest(http.MethodPost, "/reservation", bytes.NewBuffer(bts))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()
			mux.ServeHTTP(rr, req)

			if rr.Code != http.StatusCreated {
				t.Errorf("expected status code %d, got %d", http.StatusCreated, rr.Code)
				t.Log(rr.Body.String())
			}
		})

		t.Run("concurrent request to reserve the sametime", func(t *testing.T) {
			ddate := d.Add(48 * time.Hour)
			createSchedule(t, jwtSer, staff1, ddate, mux)

			// retrieve reservation times
			types := []string{serviceType1.Name, serviceType2.Name}
			reserves := reservationTimes(t, ddate, staff1, types, zone.String(), mux, err)

			maxT := len(reserves[0].Times)
			payload := model.ReservationPayload{
				StaffId:  staff1.UUID.String(),
				Name:     "user-name",
				Email:    uuid.NewString() + "@email.com",
				Phone:    "0123456789",
				Services: types,
				Time:     reserves[0].Times[rand.Intn(maxT-0)+0],
				Timezone: zone.String(),
			}

			bbyt, _ := json.Marshal(payload)

			// create reservation
			var wg sync.WaitGroup
			var mu sync.Mutex
			var statusArr []int
			var errArr []string

			randNum := rand.Intn(10-2) + 2

			for idx := 0; idx < randNum; idx++ {
				wg.Add(1)

				go func(i int) {
					defer wg.Done()

					req := httptest.NewRequest(http.MethodPost, "/reservation", bytes.NewBuffer(bbyt))
					req.Header.Set("Content-Type", "application/json")
					rr := httptest.NewRecorder()
					mux.ServeHTTP(rr, req)

					mu.Lock()
					statusArr = append(statusArr, rr.Code)
					errArr = append(errArr, rr.Body.String())
					mu.Unlock()
				}(idx)
			}

			wg.Wait()

			// assert
			num := handlers.CountResponseStatus(statusArr, 201)
			if num != 1 {
				t.Errorf("expect 1 given %v", num)
				t.Logf("%v", statusArr)
				t.Logf("%v", errArr)
			}

			num = handlers.CountResponseStatus(statusArr, 409)
			if num != (randNum - 1) {
				t.Errorf("expect %v given %v", randNum-1, num)
				t.Logf("%v", statusArr)
				t.Logf("%v", errArr)
			}
		})

		t.Run("should cancel reservation", func(t *testing.T) {
			first, err := firstReservation(ctx)
			if err != nil {
				t.Error(err.Error())
			}

			t.Run("success.", func(t *testing.T) {
				url := fmt.Sprintf("/reservation/cancel/%v", first.ReservationId)
				req := httptest.NewRequest(http.MethodPost, url, nil)
				req.Header.Set("Content-Type", "application/json")
				rr := httptest.NewRecorder()
				mux.ServeHTTP(rr, req)

				if rr.Code != http.StatusNoContent {
					t.Errorf("expected status code %d, got %d", http.StatusNoContent, rr.Code)
				}
			})

			t.Run("reject. already cancelled reservation", func(t *testing.T) {
				url := fmt.Sprintf("/reservation/cancel/%v", first.ReservationId)
				req := httptest.NewRequest(http.MethodPost, url, nil)
				req.Header.Set("Content-Type", "application/json")
				rr := httptest.NewRecorder()
				mux.ServeHTTP(rr, req)

				if rr.Code != http.StatusBadRequest {
					t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
				}

				var obj logg.Error
				if err = json.NewDecoder(rr.Body).Decode(&obj); err != nil {
					t.Error(err.Error())
				}

				str := "reservation already cancelled"
				if !strings.Contains(str, obj.Message) {
					t.Errorf("expected to contain %s given %s", str, obj.Message)
				}
			})

			t.Run("reject. invalid reservation id", func(t *testing.T) {
				req := httptest.NewRequest(http.MethodPost, "/reservation/cancel/0", nil)
				req.Header.Set("Content-Type", "application/json")
				rr := httptest.NewRecorder()
				mux.ServeHTTP(rr, req)

				if rr.Code != http.StatusNotFound {
					t.Errorf("expected status code %d, got %d", http.StatusNotFound, rr.Code)
				}
			})
		})
	})

	t.Run("CRM", func(t *testing.T) {
		obj, _ := jwtSer.Encode(context.Background(),
			&models.JwtObj{
				UserId: staff1.UUID.String(),
				AccessControls: []models.RolePermissionEnum{
					{Role: models.STAFF, Permissions: []models.PermissionEnum{models.WRITE}},
				},
			},
			utils.TwoDaysInSeconds,
		)

		t.Run("should return bookings. empty", func(t *testing.T) {
			url := fmt.Sprintf(
				"/crm/reservation?&month=%v&year=%v&user_id=%s&timezone=%s",
				int(d.Month()), d.Year()+1, staff1.UUID.String(), zone.String(),
			)

			req := httptest.NewRequest(http.MethodGet, url, nil)
			req.AddCookie(&http.Cookie{Name: env.CookieParam.CookieName, Value: obj.Token})
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()
			mux.ServeHTTP(rr, req)

			if rr.Code != http.StatusOK {
				t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
			}

			body := strings.TrimSpace(rr.Body.String())

			if body == "" {
				t.Error("expected non-empty body, got nil or empty body")
			}

			if strings.Compare(body, "[]") != 0 {
				t.Errorf("expected empty array, got: %s", body)
			}
		})

		t.Run("should return bookings", func(t *testing.T) {
			url := fmt.Sprintf(
				"/crm/reservation?&month=%v&year=%v&user_id=%s&timezone=%s",
				int(d.Month()), d.Year(), staff1.UUID.String(), zone.String(),
			)

			req := httptest.NewRequest(http.MethodGet, url, nil)
			req.AddCookie(&http.Cookie{Name: env.CookieParam.CookieName, Value: obj.Token})
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()
			mux.ServeHTTP(rr, req)

			if rr.Code != http.StatusOK {
				t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
			}

			bo := strings.TrimSpace(rr.Body.String())
			if len(bo) < 2 {
				t.Errorf("expected non-empty body, got %s", bo)
			}
		})

		t.Run("should update booking status", func(t *testing.T) {
			t.Run("reject invalid reservation id", func(t *testing.T) {
				body := model.UpdateBookingPayload{
					ReservationId: 10000,
					Status:        model.CANCELLED,
				}

				bts, err := json.Marshal(body)
				if err != nil {
					t.Fatalf("failed to marshal UpdateBookingPayload: %s", err.Error())
				}

				req := httptest.NewRequest(http.MethodPut, "/crm/reservation", bytes.NewBuffer(bts))
				req.AddCookie(&http.Cookie{Name: env.CookieParam.CookieName, Value: obj.Token})
				req.Header.Set("Content-Type", "application/json")
				rr := httptest.NewRecorder()

				mux.ServeHTTP(rr, req)

				if rr.Code != http.StatusNotFound {
					t.Errorf("expected status code %d, got %d", http.StatusNotFound, rr.Code)
				}
			})

			t.Run("success", func(t *testing.T) {
				first, err := firstReservation(ctx)
				if err != nil {
					t.Error(err.Error())
				}

				body := model.UpdateBookingPayload{
					ReservationId: first.ReservationId,
					Status:        model.CANCELLED,
				}

				bts, err := json.Marshal(body)
				if err != nil {
					t.Fatalf("failed to marshal UpdateBookingPayload: %s", err.Error())
				}

				req := httptest.NewRequest(http.MethodPut, "/crm/reservation", bytes.NewBuffer(bts))
				req.AddCookie(&http.Cookie{Name: env.CookieParam.CookieName, Value: obj.Token})
				req.Header.Set("Content-Type", "application/json")
				rr := httptest.NewRecorder()

				mux.ServeHTTP(rr, req)

				if rr.Code != http.StatusNoContent {
					t.Errorf("expected status code %d, got %d", http.StatusNoContent, rr.Code)
				}
			})
		})
	})
}

func firstReservation(ctx context.Context) (*model.Reservation, error) {
	var r model.Reservation

	q := "SELECT * FROM reservation LIMIT 1"
	row := db.QueryRowContext(ctx, q)

	err := row.Scan(
		&r.ReservationId, &r.Name, &r.Email, &r.Description, &r.Phone, &r.Price, &r.Status, &r.CreatedAt, &r.ScheduledFor, &r.ExpireAt, &r.StaffId)

	if err != nil {
		return nil, err
	}

	return &r, nil
}

func reservationTimes(t *testing.T, d time.Time, staff1 *staff.StaffEntity, serviceTypes []string, zone string, mux *http.ServeMux, err error) []model.ReservationTimeSlots {
	var sb strings.Builder
	for _, ser := range serviceTypes {
		sb.WriteString("service=" + ser + "&")
	}

	url := fmt.Sprintf(
		"/reservation?day=%v&month=%v&year=%v&staff_id=%s&%stimezone=%s",
		1, int(d.Month()), d.Year(), staff1.UUID.String(), sb.String(), zone,
	)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	var payload []model.ReservationTimeSlots
	if err = json.Unmarshal(rr.Body.Bytes(), &payload); err != nil {
		t.Error(err.Error())
	}

	return payload
}

func createSchedule(t *testing.T, jwtSer auth.IJwtService, staff1 *staff.StaffEntity, d time.Time, mux *http.ServeMux) {
	cred := []models.RolePermissionEnum{
		{Role: models.STAFF, Permissions: []models.PermissionEnum{models.WRITE}},
	}

	obj, err := jwtSer.Encode(context.Background(),
		&models.JwtObj{
			UserId:         staff1.UUID.String(),
			AccessControls: cred,
		},
		utils.TwoDaysInSeconds,
	)

	dto := scheduleModel.SchedulePayload{
		StaffId: staff1.UUID.String(),
		Times: &[]scheduleModel.ScheduleSegmentPayload{
			{
				IsVisible:     true,
				IsReoccurring: false,
				Start:         d.Format(utils.TimeFormat),
				Duration:      8 * 60 * 60, // 8 hrs shift
			},
		},
	}

	dtoBytes, err := json.Marshal(dto)
	if err != nil {
		t.Errorf("failed to marshal SchedulePayload: %s", err)
		return
	}

	// create schedule route to test
	req := httptest.NewRequest(http.MethodPost, "/schedule", bytes.NewBuffer(dtoBytes))
	req.AddCookie(&http.Cookie{Name: env.CookieParam.CookieName, Value: obj.Token})
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	// assert
	if rr.Code != http.StatusCreated {
		t.Errorf("expected status code %d, got %d", http.StatusCreated, rr.Code)
	}
}

func preSaveService(ctx context.Context, a *stores.Adapters) (*service_type.ServiceTypeEntity, error) {
	s := service_type.ServiceTypeEntity{
		Name:        uuid.New().String(),
		Price:       19.56,
		Duration:    3600,
		CleanUpTime: 30 * 60,
		IsVisible:   true,
	}
	err := a.ServiceStore.Save(ctx, &s)
	return &s, err
}

func linkServiceToStaff(a *stores.Adapters, sta *staff.StaffEntity, staSer *service_type.ServiceTypeEntity) error {
	return a.StaffServiceStore.Save(context.Background(), &staff.StaffServiceEntity{
		StaffId:   sta.StaffId,
		ServiceId: staSer.ServiceId,
	})
}

func randomTimezone() (*time.Location, error) {
	var zones []string
	var zoneDirs = []string{
		// update path according to your OS
		"/usr/share/zoneinfo/",
		"/usr/share/lib/zoneinfo/",
		"/usr/lib/locale/TZ/",
	}

	for _, zd := range zoneDirs {
		zones = walkTzDir(zd, zones)

		for idx, zone := range zones {
			zones[idx] = strings.ReplaceAll(zone, zd+"/", "")
		}
	}

	location, err := time.LoadLocation(zones[rand.Intn(len(zones))])
	if err != nil {
		return nil, err
	}

	return location, nil
}

func walkTzDir(path string, zones []string) []string {
	fileInfos, err := ioutil.ReadDir(path)
	if err != nil {
		return zones
	}

	isAlpha := func(s string) bool {
		for _, r := range s {
			if !unicode.IsLetter(r) {
				return false
			}
		}
		return true
	}

	for _, info := range fileInfos {
		if info.Name() != strings.ToUpper(info.Name()[:1])+info.Name()[1:] {
			continue
		}

		if !isAlpha(info.Name()[:1]) {
			continue
		}

		newPath := path + "/" + info.Name()

		if info.IsDir() {
			zones = walkTzDir(newPath, zones)
		} else {
			zones = append(zones, newPath)
		}
	}

	return zones
}
package reservation

import (
	"encoding/json"
	"github.com/iTchTheRightSpot/reservation/pkg"
	"github.com/iTchTheRightSpot/reservation/pkg/middleware"
	"github.com/iTchTheRightSpot/reservation/pkg/models"
	model "github.com/iTchTheRightSpot/reservation/pkg/models/reservation"
	"github.com/iTchTheRightSpot/reservation/pkg/services/reservation"
	log "github.com/iTchTheRightSpot/utility/utils"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ReservationHandler struct {
	mux     *http.ServeMux
	logger  log.ILogger
	ware    *middleware.Middleware
	service reservation.IReservationService
}

func NewReservationHandler(mux *http.ServeMux, l log.ILogger, w *middleware.Middleware, s reservation.IReservationService) *ReservationHandler {
	return &ReservationHandler{mux: mux, logger: l, service: s, ware: w}
}

func (dep *ReservationHandler) Register() {
	dep.mux.HandleFunc("GET /reservation", dep.availableDates)
	dep.mux.HandleFunc("POST /reservation/cancel/{reservation_id}", dep.cancel)
	m1 := middleware.RequestBodyMiddleware[model.ReservationPayload]{Logger: dep.logger}
	dep.mux.Handle("POST /reservation", m1.RequestBody(http.HandlerFunc(dep.create)))

	// protected routes
	rp := models.RolePermissionEnum{Role: models.STAFF, Permissions: []models.PermissionEnum{models.WRITE}}
	// read
	dep.mux.Handle("GET /crm/reservation", dep.ware.Authentication(dep.ware.HasRole(http.HandlerFunc(dep.bookings), &rp.Role)))

	// write
	dep.mux.Handle("POST /crm/reservation", dep.ware.Authentication(dep.ware.HasRoleAndPermissions(m1.RequestBody(http.HandlerFunc(dep.manualCreate)), &rp)))
	m2 := middleware.RequestBodyMiddleware[model.UpdateBookingPayload]{Logger: dep.logger}
	dep.mux.Handle("PUT /crm/reservation", dep.ware.Authentication(dep.ware.HasRoleAndPermissions(m2.RequestBody(http.HandlerFunc(dep.updateBookingStatus)), &rp)))
}

func (dep *ReservationHandler) create(w http.ResponseWriter, r *http.Request) {
	dto, err := pkg.ReadBody[model.ReservationPayload](dep.logger, r)
	if err != nil {
		dep.logger.Error(r.Context(), r.Context(), err.Error())
		log.ErrorResponse(w, err)
		return
	}

	err = dep.service.Create(r.Context(), dto)
	if err != nil {
		dep.logger.Error(r.Context(), err.Error())
		log.ErrorResponse(w, err)
		return
	}

	dep.logger.Log(r.Context(), "reservation created")
	w.WriteHeader(http.StatusCreated)
}

func (dep *ReservationHandler) availableDates(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	day, err := strconv.Atoi(query.Get("day"))
	if err != nil {
		dep.logger.Error(r.Context(), "day missing err: ", err.Error())
		log.ErrorResponse(w, &log.BadRequestError{Message: "day missing"})
		return
	}

	month, err := strconv.Atoi(query.Get("month"))
	if err != nil {
		dep.logger.Error(r.Context(), "month missing err: ", err.Error())
		log.ErrorResponse(w, &log.BadRequestError{Message: "month missing"})
		return
	}

	year, err := strconv.Atoi(query.Get("year"))
	if err != nil {
		dep.logger.Error(r.Context(), "year missing err: ", err.Error())
		log.ErrorResponse(w, &log.BadRequestError{Message: "year missing"})
		return
	}

	p := model.AvailableTimesPayload{
		StaffId:  strings.TrimSpace(query.Get("staff_id")),
		Services: query["service"],
		Day:      day,
		Month:    month,
		Year:     year,
	}

	if err = middleware.ValidatorInstance.Struct(p); err != nil {
		dep.logger.Error(r.Context(), err.Error())
		log.ErrorResponse(w, &log.BadRequestError{Message: err.Error()})
		return
	}

	location := dep.logger.Timezone()
	zone := query.Get("timezone")
	if zone != "" {
		location, err = time.LoadLocation(zone)
		if err != nil {
			dep.logger.Error(r.Context(), err.Error())
			log.ErrorResponse(w, &log.BadRequestError{Message: err.Error()})
			return
		}
	}

	p.StartDateTime = time.Date(p.Year, time.Month(p.Month), day, 0, 0, 0, 0, location)
	p.EndDateTime = time.Date(p.StartDateTime.Year(), p.StartDateTime.Month()+1, 0, 23, 59, 59, 999999999, p.StartDateTime.Location())

	dates, err := dep.service.AvailableDates(r.Context(), &p)
	if err != nil {
		dep.logger.Error(r.Context(), err.Error())
		log.ErrorResponse(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if dates == nil || len(dates) < 1 {
		if err = json.NewEncoder(w).Encode(make([]model.ReservationTimeSlots, 0)); err != nil {
			dep.logger.Error(r.Context(), err.Error())
			http.Error(w, "server error", 500)
		}
		return
	}

	if err = json.NewEncoder(w).Encode(dates); err != nil {
		dep.logger.Error(r.Context(), err.Error())
		http.Error(w, "server error", 500)
	}
}

func (dep *ReservationHandler) cancel(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("reservation_id"), 10, 64)
	if err != nil {
		dep.logger.Error(r.Context(), err.Error())
		log.ErrorResponse(w, &log.BadRequestError{Message: "invalid reservation id"})
		return
	}

	if err = dep.service.Cancel(r.Context(), uint64(id)); err != nil {
		log.ErrorResponse(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	if _, err = w.Write([]byte("reservation cancelled")); err != nil {
		dep.logger.Error(r.Context(), err.Error())
		http.Error(w, "server error", 500)
	}
}

func (dep *ReservationHandler) bookings(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	month, err := strconv.Atoi(query.Get("month"))
	if err != nil {
		dep.logger.Error(r.Context(), "month missing err: ", err.Error())
		log.ErrorResponse(w, &log.BadRequestError{Message: "month missing"})
		return
	}

	year, err := strconv.Atoi(query.Get("year"))
	if err != nil {
		dep.logger.Error(r.Context(), "year missing err: ", err.Error())
		log.ErrorResponse(w, &log.BadRequestError{Message: "year missing"})
		return
	}

	var location *time.Location
	zone := query.Get("timezone")
	if zone == "" {
		dep.logger.Error(r.Context(), "timezone missing")
		log.ErrorResponse(w, &log.BadRequestError{Message: "timezone missing"})
		return
	}

	location, err = time.LoadLocation(zone)
	if err != nil {
		dep.logger.Error(r.Context(), err.Error())
		log.ErrorResponse(w, &log.BadRequestError{Message: err.Error()})
		return
	}

	p := model.CRMBookingsPayload{
		StaffId:  strings.TrimSpace(query.Get("user_id")),
		Month:    month,
		Year:     year,
		Timezone: location,
	}

	if err = middleware.ValidatorInstance.Struct(p); err != nil {
		dep.logger.Error(r.Context(), err.Error())
		log.ErrorResponse(w, &log.BadRequestError{Message: err.Error()})
		return
	}

	dates, err := dep.service.Bookings(r.Context(), &p)
	if err != nil {
		dep.logger.Error(r.Context(), err.Error())
		log.ErrorResponse(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(dates); err != nil {
		dep.logger.Error(r.Context(), err.Error())
		http.Error(w, "server error", 500)
	}
}

func (dep *ReservationHandler) updateBookingStatus(w http.ResponseWriter, r *http.Request) {
	dto, err := pkg.ReadBody[model.UpdateBookingPayload](dep.logger, r)
	if err != nil {
		dep.logger.Error(r.Context(), err.Error())
		log.ErrorResponse(w, err)
		return
	}

	if err = dep.service.UpdateBookingStatus(r.Context(), dto); err != nil {
		dep.logger.Error(r.Context(), err.Error())
		log.ErrorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (dep *ReservationHandler) manualCreate(w http.ResponseWriter, r *http.Request) {
	dto, err := pkg.ReadBody[model.ReservationPayload](dep.logger, r)
	if err != nil {
		dep.logger.Error(r.Context(), err.Error())
		log.ErrorResponse(w, err)
		return
	}

	err = dep.service.ManualCreate(r.Context(), dto)
	if err != nil {
		dep.logger.Error(r.Context(), err.Error())
		log.ErrorResponse(w, err)
		return
	}

	dep.logger.Log(r.Context(), "manually created reservation")
	w.WriteHeader(http.StatusCreated)
}
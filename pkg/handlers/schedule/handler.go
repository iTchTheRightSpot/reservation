package schedule

import (
	"encoding/json"
	"github.com/iTchTheRightSpot/reservation/pkg"
	"github.com/iTchTheRightSpot/reservation/pkg/middleware"
	"github.com/iTchTheRightSpot/reservation/pkg/models"
	model "github.com/iTchTheRightSpot/reservation/pkg/models/schedule"
	"github.com/iTchTheRightSpot/reservation/pkg/services/schedule"
	"github.com/iTchTheRightSpot/reservation/utils"
	mid "github.com/iTchTheRightSpot/utility/middleware"
	log "github.com/iTchTheRightSpot/utility/utils"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ScheduleHandler struct {
	mux     *http.ServeMux
	ware    *middleware.Middleware
	logger  log.ILogger
	service schedule.IScheduleService
}

func NewScheduleHandler(mux *http.ServeMux, w *middleware.Middleware, l log.ILogger, s schedule.IScheduleService) *ScheduleHandler {
	return &ScheduleHandler{mux: mux, ware: w, logger: l, service: s}
}

func (dep *ScheduleHandler) Register() {
	rp := &models.RolePermissionEnum{Role: models.STAFF, Permissions: []models.PermissionEnum{models.WRITE}}

	// read
	dep.mux.Handle("GET /schedules", dep.ware.Authentication(dep.ware.HasRole(http.HandlerFunc(dep.schedules), &rp.Role)))
	dep.mux.Handle("GET /schedules/staff", dep.ware.Authentication(dep.ware.HasRole(http.HandlerFunc(dep.schedulesByStaff), &rp.Role)))

	// write
	m1 := mid.RequestBodyMiddleware[model.SchedulePayload]{Logger: dep.logger, Validator: dep.ware.Validator}
	dep.mux.Handle("POST /schedule", dep.ware.Authentication(dep.ware.HasRoleAndPermissions(m1.RequestBody(http.HandlerFunc(dep.create)), rp)))

	m2 := mid.RequestBodyMiddleware[model.UpdateSchedulePayload]{Logger: dep.logger, Validator: dep.ware.Validator}
	dep.mux.Handle("PUT /schedule", dep.ware.Authentication(dep.ware.HasRoleAndPermissions(m2.RequestBody(http.HandlerFunc(dep.update)), rp)))
	dep.mux.Handle("DELETE /schedule/{schedule_id}", dep.ware.Authentication(
		dep.ware.HasRoleAndPermissions(
			http.HandlerFunc(dep.delete),
			&models.RolePermissionEnum{Role: models.STAFF, Permissions: []models.PermissionEnum{models.DELETE}},
		),
	))
}

func (dep *ScheduleHandler) create(w http.ResponseWriter, r *http.Request) {
	dto, err := pkg.ReadBody[model.SchedulePayload](dep.logger, r)
	if err != nil {
		dep.logger.Error(r.Context(), err.Error())
		log.ErrorResponse(w, &log.BadRequestError{Message: err.Error()})
		return
	}

	if err = dep.service.Create(r.Context(), dto); err != nil {
		dep.logger.Error(r.Context(), err.Error())
		log.ErrorResponse(w, err)
		return
	}

	dep.logger.Log(r.Context(), "new schedule created")
	w.WriteHeader(http.StatusCreated)
}

func (dep *ScheduleHandler) update(w http.ResponseWriter, r *http.Request) {
	dto, err := pkg.ReadBody[model.UpdateSchedulePayload](dep.logger, r)
	if err != nil {
		dep.logger.Error(r.Context(), err.Error())
		log.ErrorResponse(w, &log.BadRequestError{Message: err.Error()})
		return
	}

	if err = dep.service.Update(r.Context(), dto); err != nil {
		dep.logger.Error(r.Context(), err.Error())
		log.ErrorResponse(w, err)
		return
	}

	dep.logger.Log(r.Context(), "schedule updated")
	w.WriteHeader(http.StatusNoContent)
}

func (dep *ScheduleHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("schedule_id"), 10, 64)
	if err != nil {
		dep.logger.Error(r.Context(), err.Error())
		log.ErrorResponse(w, &log.BadRequestError{Message: "invalid schedule_id"})
		return
	}

	if err = dep.service.Delete(r.Context(), id); err != nil {
		dep.logger.Error(r.Context(), err.Error())
		log.ErrorResponse(w, err)
		return
	}

	dep.logger.Log(r.Context(), "schedule deleted")
	w.WriteHeader(http.StatusNoContent)
}

func (dep *ScheduleHandler) schedulesByStaff(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	month, err := strconv.Atoi(query.Get("month"))
	if err != nil {
		log.ErrorResponse(w, &log.BadRequestError{Message: err.Error()})
		return
	}

	year, err := strconv.Atoi(query.Get("year"))
	if err != nil {
		log.ErrorResponse(w, &log.BadRequestError{Message: err.Error()})
		return
	}

	location := dep.logger.Timezone()

	if query.Get("timezone") != "" {
		location, err = time.LoadLocation(query.Get("timezone"))
		if err != nil {
			log.ErrorResponse(w, &log.BadRequestError{Message: err.Error()})
			return
		}
	}

	staffId := query.Get("staff_id")
	if staffId == "" {
		obj, ok := r.Context().Value(utils.UserContextKey).(*models.JwtObj)
		if !ok || obj == nil {
			dep.logger.Error(r.Context(), "jwt object not present in request")
			log.ErrorResponse(w, &log.AuthenticationError{})
			return
		}
		staffId = obj.UserId
	}

	payload := model.AllSchedulesPayload{
		StaffUUID: strings.TrimSpace(staffId),
		Month:     month,
		Year:      year,
		Timezone:  location,
	}

	if err = dep.ware.Validator.Struct(payload); err != nil {
		dep.logger.Error(r.Context(), err.Error())
		log.ErrorResponse(w, &log.BadRequestError{Message: err.Error()})
		return
	}

	arr, err := dep.service.Schedules(r.Context(), &payload)
	if err != nil {
		dep.logger.Error(r.Context(), err.Error())
		log.ErrorResponse(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(arr); err != nil {
		dep.logger.Error(r.Context(), err.Error())
		log.ErrorResponse(w, &log.ServerError{})
	}
}

func (dep *ScheduleHandler) schedules(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	month, err := strconv.Atoi(query.Get("month"))
	if err != nil {
		log.ErrorResponse(w, &log.BadRequestError{Message: err.Error()})
		return
	}

	year, err := strconv.Atoi(query.Get("year"))
	if err != nil {
		log.ErrorResponse(w, &log.BadRequestError{Message: err.Error()})
		return
	}

	location := dep.logger.Timezone()

	if query.Get("timezone") != "" {
		location, err = time.LoadLocation(query.Get("timezone"))
		if err != nil {
			log.ErrorResponse(w, &log.BadRequestError{Message: err.Error()})
			return
		}
	}

	obj, ok := r.Context().Value(utils.UserContextKey).(*models.JwtObj)
	if !ok || obj == nil {
		dep.logger.Error(r.Context(), "jwt object not present in request")
		log.ErrorResponse(w, &log.AuthenticationError{})
		return
	}

	payload := model.AllSchedulesPayload{
		StaffUUID: strings.TrimSpace(obj.UserId),
		Month:     month,
		Year:      year,
		Timezone:  location,
	}

	if err = dep.ware.Validator.Struct(payload); err != nil {
		dep.logger.Error(r.Context(), err.Error())
		log.ErrorResponse(w, &log.BadRequestError{Message: err.Error()})
		return
	}

	arr, err := dep.service.Schedules(r.Context(), &payload)
	if err != nil {
		dep.logger.Error(r.Context(), err.Error())
		log.ErrorResponse(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(arr); err != nil {
		dep.logger.Error(r.Context(), err.Error())
		log.ErrorResponse(w, &log.ServerError{})
	}
}
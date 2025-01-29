package schedule

import (
	"encoding/json"
	"github.com/iTchTheRightSpot/erp-golang/pkg"
	"github.com/iTchTheRightSpot/erp-golang/pkg/middleware"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models"
	model "github.com/iTchTheRightSpot/erp-golang/pkg/models/schedule"
	"github.com/iTchTheRightSpot/erp-golang/pkg/services/schedule"
	"github.com/iTchTheRightSpot/erp-golang/utils"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ScheduleHandler struct {
	mux     *http.ServeMux
	ware    *middleware.Middleware
	logger  utils.ILogger
	service schedule.IScheduleService
}

func NewScheduleHandler(mux *http.ServeMux, w *middleware.Middleware, l utils.ILogger, s schedule.IScheduleService) *ScheduleHandler {
	return &ScheduleHandler{mux: mux, ware: w, logger: l, service: s}
}

func (dep *ScheduleHandler) Register() {
	rp := &models.RolePermissionEnum{Role: models.STAFF, Permissions: []models.PermissionEnum{models.WRITE}}

	// read
	dep.mux.Handle("GET /schedules", dep.ware.Authentication(dep.ware.HasRole(http.HandlerFunc(dep.schedules), &rp.Role)))
	dep.mux.Handle("GET /schedules/staff", dep.ware.Authentication(dep.ware.HasRoleAndPermissions(http.HandlerFunc(dep.schedulesByStaff), rp)))

	// write
	m1 := middleware.RequestBodyMiddleware[model.SchedulePayload]{Logger: dep.logger}
	dep.mux.Handle("POST /schedule", dep.ware.Authentication(dep.ware.HasRoleAndPermissions(m1.RequestBody(http.HandlerFunc(dep.create)), rp)))

	m2 := middleware.RequestBodyMiddleware[model.UpdateSchedulePayload]{Logger: dep.logger}
	dep.mux.Handle("PUT /schedule", dep.ware.Authentication(dep.ware.HasRoleAndPermissions(m2.RequestBody(http.HandlerFunc(dep.update)), rp)))
	dep.mux.Handle("DELETE /schedule/{schedule_id}", dep.ware.Authentication(
		dep.ware.HasRoleAndPermissions(
			http.HandlerFunc(dep.delete),
			&models.RolePermissionEnum{Role: models.STAFF, Permissions: []models.PermissionEnum{models.DELETE}},
		),
	))
}

func (dep *ScheduleHandler) create(w http.ResponseWriter, r *http.Request) {
	dto, err := pkg.ReadBody[model.SchedulePayload](r)
	if err != nil {
		dep.logger.Error(err.Error())
		utils.ErrorResponse(w, &utils.BadRequestError{Message: err.Error()})
		return
	}

	if err = dep.service.Create(r.Context(), dto); err != nil {
		dep.logger.Error(err.Error())
		utils.ErrorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (dep *ScheduleHandler) update(w http.ResponseWriter, r *http.Request) {
	dto, err := pkg.ReadBody[model.UpdateSchedulePayload](r)
	if err != nil {
		dep.logger.Error(err.Error())
		utils.ErrorResponse(w, &utils.BadRequestError{Message: err.Error()})
		return
	}

	if err = dep.service.Update(r.Context(), dto); err != nil {
		dep.logger.Error(err.Error())
		utils.ErrorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (dep *ScheduleHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("schedule_id"), 10, 64)
	if err != nil {
		dep.logger.Error(err.Error())
		utils.ErrorResponse(w, &utils.BadRequestError{Message: "Bad request, invalid schedule_id"})
		return
	}

	if err = dep.service.Delete(r.Context(), id); err != nil {
		dep.logger.Error(err.Error())
		utils.ErrorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (dep *ScheduleHandler) schedulesByStaff(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	month, err := strconv.Atoi(query.Get("month"))
	if err != nil {
		utils.ErrorResponse(w, &utils.BadRequestError{Message: err.Error()})
		return
	}

	year, err := strconv.Atoi(query.Get("year"))
	if err != nil {
		utils.ErrorResponse(w, &utils.BadRequestError{Message: err.Error()})
		return
	}

	location := dep.logger.Timezone()

	if query.Get("timezone") != "" {
		location, err = time.LoadLocation(query.Get("timezone"))
		if err != nil {
			utils.ErrorResponse(w, &utils.BadRequestError{Message: err.Error()})
			return
		}
	}

	staffId := query.Get("staff_id")
	if staffId == "" {
		obj, ok := r.Context().Value(utils.UserContextKey).(*models.JwtObj)
		if !ok || obj == nil {
			dep.logger.Error("jwt object not present in request")
			utils.ErrorResponse(w, &utils.AuthenticationError{})
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

	if err = middleware.ValidatorInstance.Struct(payload); err != nil {
		dep.logger.Error(err.Error())
		utils.ErrorResponse(w, &utils.BadRequestError{Message: err.Error()})
		return
	}

	arr, err := dep.service.Schedules(r.Context(), &payload)
	if err != nil {
		dep.logger.Error(err.Error())
		utils.ErrorResponse(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(arr); err != nil {
		dep.logger.Error(err.Error())
	}
}

func (dep *ScheduleHandler) schedules(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	month, err := strconv.Atoi(query.Get("month"))
	if err != nil {
		utils.ErrorResponse(w, &utils.BadRequestError{Message: err.Error()})
		return
	}

	year, err := strconv.Atoi(query.Get("year"))
	if err != nil {
		utils.ErrorResponse(w, &utils.BadRequestError{Message: err.Error()})
		return
	}

	location := dep.logger.Timezone()

	if query.Get("timezone") != "" {
		location, err = time.LoadLocation(query.Get("timezone"))
		if err != nil {
			utils.ErrorResponse(w, &utils.BadRequestError{Message: err.Error()})
			return
		}
	}

	obj, ok := r.Context().Value(utils.UserContextKey).(*models.JwtObj)
	if !ok || obj == nil {
		dep.logger.Error("jwt object not present in request")
		utils.ErrorResponse(w, &utils.AuthenticationError{})
		return
	}

	payload := model.AllSchedulesPayload{
		StaffUUID: strings.TrimSpace(obj.UserId),
		Month:     month,
		Year:      year,
		Timezone:  location,
	}

	if err = middleware.ValidatorInstance.Struct(payload); err != nil {
		dep.logger.Error(err.Error())
		utils.ErrorResponse(w, &utils.BadRequestError{Message: err.Error()})
		return
	}

	arr, err := dep.service.Schedules(r.Context(), &payload)
	if err != nil {
		dep.logger.Error(err.Error())
		utils.ErrorResponse(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(arr); err != nil {
		dep.logger.Error(err.Error())
	}
}

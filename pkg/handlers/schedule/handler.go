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
	mux := http.NewServeMux()
	role := models.STAFF
	permission := models.WRITE
	m := middleware.RequestBodyMiddleware[model.SchedulePayload]{Logger: dep.logger}

	mux.Handle("POST /", dep.ware.HasPermission(m.RequestBody(http.HandlerFunc(dep.create)), &permission))
	mux.HandleFunc("GET /", dep.schedules)
	dep.mux.Handle("/schedule", dep.ware.Authentication(dep.ware.HasRole(mux, &role)))
}

func (dep *ScheduleHandler) create(w http.ResponseWriter, r *http.Request) {
	dto, err := pkg.ReadBody[model.SchedulePayload](r)
	if err != nil {
		dep.logger.Error(err.Error())
		utils.ConstructErrorResponse(
			w,
			utils.ErrorResponse{
				Status:  http.StatusInternalServerError,
				Message: err.Error(),
			},
		)
		return
	}

	if err = dep.service.Create(r.Context(), dto); err != nil {
		dep.logger.Error(err.Error())
		utils.ConstructErrorResponse(
			w,
			utils.ErrorResponse{
				Status:  http.StatusBadRequest,
				Message: err.Error(),
			},
		)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (dep *ScheduleHandler) schedules(w http.ResponseWriter, r *http.Request) {
	staffUUID := r.URL.Query().Get("staff_id")
	if len(staffUUID) < 1 {
		obj, ok := r.Context().Value(utils.UserContextKey).(*models.JwtObj)
		if !ok || obj == nil {
			dep.logger.Error("linkServiceToStaff invalid staff_id")
			utils.ConstructErrorResponse(
				w,
				utils.ErrorResponse{
					Status:  http.StatusUnauthorized,
					Message: "full authentication is required to access this resource",
				},
			)
			return
		}
		staffUUID = obj.UserId
	}

	month, err := strconv.Atoi(r.URL.Query().Get("month"))
	if err != nil {
		utils.ConstructErrorResponse(
			w,
			utils.ErrorResponse{
				Status:  http.StatusBadRequest,
				Message: err.Error(),
			},
		)
		return
	}

	year, err := strconv.Atoi(r.URL.Query().Get("year"))
	if err != nil {
		utils.ConstructErrorResponse(
			w,
			utils.ErrorResponse{
				Status:  http.StatusBadRequest,
				Message: err.Error(),
			},
		)
		return
	}

	var tz time.Location
	timezone := r.URL.Query().Get("timezone")

	if len(timezone) < 1 {
		tz = *dep.logger.Timezone()
	} else {
		location, err := time.LoadLocation(timezone)
		if err != nil {
			utils.ConstructErrorResponse(
				w,
				utils.ErrorResponse{
					Status:  http.StatusBadRequest,
					Message: err.Error(),
				},
			)
			return
		}
		tz = *location
	}

	payload := model.AllSchedulesPayload{
		StaffUUID: strings.TrimSpace(staffUUID),
		Month:     month,
		Year:      year,
		Timezone:  &tz,
	}

	if err = middleware.ValidatorInstance.Struct(payload); err != nil {
		dep.logger.Error(err.Error())
		utils.ConstructErrorResponse(
			w,
			utils.ErrorResponse{
				Status:  http.StatusBadRequest,
				Message: err.Error(),
			},
		)
		return
	}

	schedules, err := dep.service.Schedules(r.Context(), &payload)
	if err != nil {
		dep.logger.Error(err.Error())
		utils.ConstructErrorResponse(
			w,
			utils.ErrorResponse{
				Status:  http.StatusBadRequest,
				Message: err.Error(),
			},
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response, _ := json.Marshal(schedules)
	if _, err = w.Write(response); err != nil {
		dep.logger.Error(err.Error())
	}
}

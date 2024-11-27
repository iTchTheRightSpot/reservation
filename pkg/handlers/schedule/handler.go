package schedule

import (
	"github.com/iTchTheRightSpot/erp-golang/pkg"
	"github.com/iTchTheRightSpot/erp-golang/pkg/middleware"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models"
	model "github.com/iTchTheRightSpot/erp-golang/pkg/models/schedule"
	"github.com/iTchTheRightSpot/erp-golang/pkg/services/schedule"
	"github.com/iTchTheRightSpot/erp-golang/utils"
	"net/http"
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
	mux.Handle("GET /", http.HandlerFunc(dep.schedules))
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
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("schedule handler GET route hit"))
}

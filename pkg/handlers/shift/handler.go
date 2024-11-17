package shift

import (
	"github.com/iTchTheRightSpot/erp-golang/pkg"
	"github.com/iTchTheRightSpot/erp-golang/pkg/middleware"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models"
	shiftModel "github.com/iTchTheRightSpot/erp-golang/pkg/models/shift"
	"github.com/iTchTheRightSpot/erp-golang/pkg/services/shift"
	"github.com/iTchTheRightSpot/erp-golang/utils"
	"net/http"
)

type ShiftHandler struct {
	mux        *http.ServeMux
	middleware *middleware.Middleware
	logger     utils.ILogger
	service    shift.IShiftService
}

func NewShiftHandler(mux *http.ServeMux, w *middleware.Middleware, l utils.ILogger, s shift.IShiftService) *ShiftHandler {
	return &ShiftHandler{mux: mux, middleware: w, logger: l, service: s}
}

func (dep *ShiftHandler) Register() {
	shiftMux := http.NewServeMux()
	staff := models.STAFF

	m := middleware.RequestBodyMiddleware[shiftModel.ShiftDto]{Logger: dep.logger}
	shiftMux.Handle(
		"POST /",
		dep.middleware.HasRoleAndPermissions(
			m.RequestBody(http.HandlerFunc(dep.create)),
			&models.RolePermission{
				Role:        staff,
				Permissions: []models.PermissionEnum{models.WRITE},
			},
		),
	)
	shiftMux.Handle("GET /", http.HandlerFunc(dep.shifts))

	dep.mux.Handle("/shift", dep.middleware.Authentication(dep.middleware.HasRole(shiftMux, &staff)))
}

func (dep *ShiftHandler) create(w http.ResponseWriter, r *http.Request) {
	dto, err := pkg.ReadBody[shiftModel.ShiftDto](r)
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

func (dep *ShiftHandler) shifts(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("shift handler GET route hit"))
}

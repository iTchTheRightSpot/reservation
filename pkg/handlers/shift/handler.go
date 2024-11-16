package shift

import (
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

	m := middleware.RequestBodyMiddleware[shiftModel.ShiftDto]{Logger: dep.logger}
	shiftMux.Handle(
		"POST /",
		dep.middleware.HasRoleAndPermissions(
			m.RequestBody(http.HandlerFunc(dep.create)),
			&models.RolePermission{
				Role:        models.STAFF,
				Permissions: []models.PermissionEnum{models.WRITE},
			},
		),
	)
	shiftMux.Handle("GET /", http.HandlerFunc(dep.shifts))

	dep.mux.Handle("/shift", dep.middleware.Authentication(shiftMux))
}

func (dep *ShiftHandler) create(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte("shift handler post route hit"))
}

func (dep *ShiftHandler) shifts(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("shift handler GET route hit"))
}

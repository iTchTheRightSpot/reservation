package staff

import (
	"encoding/json"
	"github.com/iTchTheRightSpot/erp-golang/pkg/middleware"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models"
	"github.com/iTchTheRightSpot/erp-golang/pkg/services/staff"
	"github.com/iTchTheRightSpot/erp-golang/utils"
	"net/http"
)

type StaffHandler struct {
	mux     *http.ServeMux
	ware    *middleware.Middleware
	logger  utils.ILogger
	service staff.IStaffService
}

func NewStaffHandler(mux *http.ServeMux, w *middleware.Middleware, l utils.ILogger, s staff.IStaffService) *StaffHandler {
	return &StaffHandler{mux: mux, ware: w, logger: l, service: s}
}

func (dep *StaffHandler) Register() {
	rp := &models.RolePermissionEnum{Role: models.STAFF, Permissions: []models.PermissionEnum{models.WRITE}}

	dep.mux.Handle("GET /staffs", dep.ware.Authentication(dep.ware.HasRoleAndPermissions(http.HandlerFunc(dep.allStaffs), rp)))
}

func (dep *StaffHandler) allStaffs(w http.ResponseWriter, r *http.Request) {
	arr, err := dep.service.AllUsers(r.Context())
	if err != nil {
		utils.ErrorResponse(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(arr); err != nil {
		dep.logger.Error(err.Error())
	}
}

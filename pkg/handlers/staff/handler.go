package staff

import (
	"encoding/json"
	"github.com/iTchTheRightSpot/erp-golang/pkg/middleware"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models"
	"github.com/iTchTheRightSpot/erp-golang/pkg/services/staff"
	"github.com/iTchTheRightSpot/erp-golang/utils"
	"net/http"
	"strings"
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
	rp := &models.RolePermissionEnum{
		Role:        models.STAFF,
		Permissions: []models.PermissionEnum{models.WRITE},
	}

	dep.mux.Handle("GET /staffs", dep.ware.Authentication(dep.ware.HasRoleAndPermissions(http.HandlerFunc(dep.allStaffs), rp)))
	dep.mux.Handle("POST /staff/service", dep.ware.Authentication(dep.ware.HasRoleAndPermissions(http.HandlerFunc(dep.linkServiceToStaff), rp)))
}

func (dep *StaffHandler) allStaffs(w http.ResponseWriter, r *http.Request) {
	arr, err := dep.service.AllStaffs(r.Context())
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

func (dep *StaffHandler) linkServiceToStaff(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("service_name")
	if len(name) < 1 {
		utils.ErrorResponse(w, &utils.BadRequestError{Message: "service_name is missing"})
		return
	}

	uuid := r.URL.Query().Get("staff_id")

	if len(uuid) < 1 {
		obj, ok := r.Context().Value(utils.UserContextKey).(*models.JwtObj)
		if !ok || obj == nil {
			dep.logger.Error("linkServiceToStaff invalid staff_id")
			utils.ErrorResponse(w, &utils.AuthenticationError{})
			return
		}
		uuid = obj.UserId
	}

	err := dep.service.LinkServiceToStaff(r.Context(), strings.TrimSpace(uuid), strings.TrimSpace(name))
	if err != nil {
		dep.logger.Error(err.Error())
		utils.ErrorResponse(w, err)
		return
	}

	dep.logger.Log("successfully linked service to staff")
	w.WriteHeader(http.StatusCreated)
}

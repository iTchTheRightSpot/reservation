package staff

import (
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
	mux := http.NewServeMux()
	role := models.STAFF
	rp := &models.RolePermission{
		Role:        role,
		Permissions: []models.PermissionEnum{models.WRITE},
	}

	mux.Handle("POST /service", dep.ware.HasRoleAndPermissions(http.HandlerFunc(dep.linkServiceToStaff), rp))

	dep.mux.Handle("/staff/", http.StripPrefix("/staff", dep.ware.Authentication(dep.ware.HasRole(mux, &role))))
}

func (dep *StaffHandler) linkServiceToStaff(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("service_name")
	if len(name) < 1 {
		utils.ConstructErrorResponse(w, &utils.BadRequestError{Message: "service_name is missing"})
		return
	}

	staffUUID := r.URL.Query().Get("staff_id")

	if len(staffUUID) < 1 {
		obj, ok := r.Context().Value(utils.UserContextKey).(*models.JwtObj)
		if !ok || obj == nil {
			dep.logger.Error("linkServiceToStaff invalid staff_id")
			utils.ConstructErrorResponse(w, &utils.AuthenticationError{})
			return
		}
		staffUUID = obj.UserId
	}

	err := dep.service.LinkServiceToStaff(r.Context(), strings.TrimSpace(staffUUID), strings.TrimSpace(name))
	if err != nil {
		dep.logger.Error(err.Error())
		utils.ConstructErrorResponse(w, err)
		return
	}

	dep.logger.Log("successfully linked service to staff")
	w.WriteHeader(http.StatusCreated)
}

package service_type

import (
	"encoding/json"
	"github.com/iTchTheRightSpot/erp-golang/pkg"
	"github.com/iTchTheRightSpot/erp-golang/pkg/middleware"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models"
	model "github.com/iTchTheRightSpot/erp-golang/pkg/models/service_type"
	"github.com/iTchTheRightSpot/erp-golang/pkg/services/service_type"
	"github.com/iTchTheRightSpot/erp-golang/utils"
	"net/http"
	"strings"
)

type ServiceTypeHandler struct {
	mux     *http.ServeMux
	logger  utils.ILogger
	service service_type.IServiceType
	ware    *middleware.Middleware
}

func NewServiceTypeHandler(mux *http.ServeMux, l utils.ILogger, s service_type.IServiceType, w *middleware.Middleware) *ServiceTypeHandler {
	return &ServiceTypeHandler{
		mux:     mux,
		logger:  l,
		service: s,
		ware:    w,
	}
}

func (dep *ServiceTypeHandler) Register() {
	rp := &models.RolePermissionEnum{
		Role:        models.STAFF,
		Permissions: []models.PermissionEnum{models.WRITE},
	}

	m := middleware.RequestBodyMiddleware[model.ServiceTypePayload]{Logger: dep.logger}

	dep.mux.HandleFunc("GET /service", dep.services)
	dep.mux.HandleFunc("GET /service/staffs", dep.staffsByServices)

	dep.mux.Handle("POST /service/staff", dep.ware.Authentication(dep.ware.HasRoleAndPermissions(http.HandlerFunc(dep.linkServiceToStaff), rp)))
	dep.mux.Handle("GET /crm/services", dep.ware.Authentication(dep.ware.HasRole(http.HandlerFunc(dep.crmServices), &rp.Role)))
	dep.mux.Handle("POST /crm/service", dep.ware.Authentication(dep.ware.HasRoleAndPermissions(m.RequestBody(http.HandlerFunc(dep.create)), rp)))
	dep.mux.Handle("PUT /crm/service", dep.ware.Authentication(dep.ware.HasRoleAndPermissions(m.RequestBody(http.HandlerFunc(dep.update)), rp)))
}

func (dep *ServiceTypeHandler) create(w http.ResponseWriter, r *http.Request) {
	dto, err := pkg.ReadBody[model.ServiceTypePayload](r)
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

func (dep *ServiceTypeHandler) services(w http.ResponseWriter, r *http.Request) {
	arr, err := dep.service.ServiceTypes(r.Context())
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

func (dep *ServiceTypeHandler) update(w http.ResponseWriter, r *http.Request) {
	dto, err := pkg.ReadBody[model.ServiceTypePayload](r)
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

func (dep *ServiceTypeHandler) staffsByServices(w http.ResponseWriter, r *http.Request) {
	services := r.URL.Query()["name"]
	if services == nil || len(services) < 1 {
		utils.ErrorResponse(w, &utils.BadRequestError{Message: "bad request, missing services type(s)"})
		return
	}

	arr, err := dep.service.StaffsByServiceTypes(r.Context(), &services)
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

func (dep *ServiceTypeHandler) crmServices(w http.ResponseWriter, r *http.Request) {
	arr, err := dep.service.CRMServiceTypes(r.Context())
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

func (dep *ServiceTypeHandler) linkServiceToStaff(w http.ResponseWriter, r *http.Request) {
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

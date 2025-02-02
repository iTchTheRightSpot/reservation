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
	rp := &models.RolePermissionEnum{Role: models.STAFF, Permissions: []models.PermissionEnum{models.WRITE}}

	dep.mux.HandleFunc("GET /service", dep.services)
	dep.mux.HandleFunc("GET /service/staffs", dep.staffsByServices)

	// protected routes
	// reads
	dep.mux.Handle("GET /crm/services", dep.ware.Authentication(
		dep.ware.HasRole(http.HandlerFunc(dep.crmServices), &rp.Role),
	))
	dep.mux.Handle("GET /crm/services/staff", dep.ware.Authentication(
		dep.ware.HasRoleAndPermissions(
			http.HandlerFunc(dep.serviceTypesByStaffUUID),
			&models.RolePermissionEnum{Role: models.STAFF, Permissions: []models.PermissionEnum{models.READ}},
		),
	))

	// writes
	m1 := middleware.RequestBodyMiddleware[model.LinkServiceTypeToStaffPayload]{Logger: dep.logger}
	dep.mux.Handle("POST /crm/service/staff", dep.ware.Authentication(
		dep.ware.HasRoleAndPermissions(
			m1.RequestBody(http.HandlerFunc(dep.linkServiceToStaff)),
			rp,
		),
	))
	dep.mux.Handle("DELETE /crm/service/staff", dep.ware.Authentication(
		dep.ware.HasRoleAndPermissions(
			http.HandlerFunc(dep.deLinkServiceFromStaff),
			&models.RolePermissionEnum{Role: models.STAFF, Permissions: []models.PermissionEnum{models.DELETE}},
		)),
	)
	m2 := middleware.RequestBodyMiddleware[model.ServiceTypePayload]{Logger: dep.logger}
	dep.mux.Handle("POST /crm/service", dep.ware.Authentication(
		dep.ware.HasRoleAndPermissions(m2.RequestBody(http.HandlerFunc(dep.create)), rp)))
	dep.mux.Handle("PUT /crm/service", dep.ware.Authentication(
		dep.ware.HasRoleAndPermissions(m2.RequestBody(http.HandlerFunc(dep.update)), rp)))
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
	dto, err := pkg.ReadBody[model.LinkServiceTypeToStaffPayload](r)
	if err != nil {
		dep.logger.Error(err.Error())
		utils.ErrorResponse(w, &utils.BadRequestError{Message: err.Error()})
		return
	}

	if err = dep.service.LinkServiceToStaff(r.Context(), dto); err != nil {
		dep.logger.Error(err.Error())
		utils.ErrorResponse(w, err)
		return
	}

	dep.logger.Log("successfully linked service to staff")
	w.WriteHeader(http.StatusCreated)
}

func (dep *ServiceTypeHandler) serviceTypesByStaffUUID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("staff_id")
	if len(id) < 1 {
		utils.ErrorResponse(w, &utils.BadRequestError{Message: "staff_id missing. bad request"})
		return
	}

	arr, err := dep.service.ServicesByStaffUUID(r.Context(), id)
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

func (dep *ServiceTypeHandler) deLinkServiceFromStaff(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)
	if err := json.NewEncoder(w).Encode(utils.InsertionError{Message: "route not implemented"}); err != nil {
		dep.logger.Error(err.Error())
	}
}

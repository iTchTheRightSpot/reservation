package service_type

import (
	"encoding/json"
	"github.com/iTchTheRightSpot/erp-golang/pkg"
	"github.com/iTchTheRightSpot/erp-golang/pkg/middleware"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models"
	model "github.com/iTchTheRightSpot/erp-golang/pkg/models/service"
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
	rp := &models.RolePermissionEnum{
		Role:        models.STAFF,
		Permissions: []models.PermissionEnum{models.WRITE},
	}

	m := middleware.RequestBodyMiddleware[model.ServiceTypePayload]{Logger: dep.logger}

	dep.mux.HandleFunc("GET /service", dep.services)
	dep.mux.HandleFunc("GET /service/staffs", dep.staffsByServices)
	dep.mux.Handle("POST /service", dep.ware.Authentication(dep.ware.HasRoleAndPermissions(m.RequestBody(http.HandlerFunc(dep.create)), rp)))
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

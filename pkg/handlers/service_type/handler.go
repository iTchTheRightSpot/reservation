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

type ServiceHandler struct {
	mux     *http.ServeMux
	logger  utils.ILogger
	service service_type.IService
	ware    *middleware.Middleware
}

func NewServiceHandler(mux *http.ServeMux, l utils.ILogger, s service_type.IService, w *middleware.Middleware) *ServiceHandler {
	return &ServiceHandler{
		mux:     mux,
		logger:  l,
		service: s,
		ware:    w,
	}
}

func (dep *ServiceHandler) Register() {
	rp := &models.RolePermission{
		Role:        models.STAFF,
		Permissions: []models.PermissionEnum{models.WRITE},
	}

	m := middleware.RequestBodyMiddleware[model.ServiceTypePayload]{Logger: dep.logger}

	dep.mux.HandleFunc("GET /service", dep.services)
	dep.mux.Handle("POST /service", dep.ware.Authentication(dep.ware.HasRoleAndPermissions(m.RequestBody(http.HandlerFunc(dep.create)), rp)))
}

func (dep *ServiceHandler) create(w http.ResponseWriter, r *http.Request) {
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

func (dep *ServiceHandler) services(w http.ResponseWriter, r *http.Request) {
	arr, err := dep.service.Services(r.Context())
	if err != nil {
		dep.logger.Error(err.Error())
		utils.ErrorResponse(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response, _ := json.Marshal(arr)
	if _, err = w.Write(response); err != nil {
		dep.logger.Error(err.Error())
	}
}

package service

import (
	"github.com/iTchTheRightSpot/erp-golang/pkg"
	"github.com/iTchTheRightSpot/erp-golang/pkg/middleware"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models"
	model "github.com/iTchTheRightSpot/erp-golang/pkg/models/service"
	"github.com/iTchTheRightSpot/erp-golang/pkg/services/service"
	"github.com/iTchTheRightSpot/erp-golang/utils"
	"net/http"
)

type ServiceHandler struct {
	mux     *http.ServeMux
	logger  utils.ILogger
	service service.IService
	ware    *middleware.Middleware
}

func NewServiceHandler(mux *http.ServeMux, l utils.ILogger, s service.IService, w *middleware.Middleware) *ServiceHandler {
	return &ServiceHandler{
		mux:     mux,
		logger:  l,
		service: s,
		ware:    w,
	}
}

func (dep *ServiceHandler) Register() {
	mux := http.NewServeMux()
	staff := models.STAFF
	rp := &models.RolePermission{
		Role:        staff,
		Permissions: []models.PermissionEnum{models.WRITE},
	}
	m := middleware.RequestBodyMiddleware[model.ServicePayload]{Logger: dep.logger}

	mux.Handle("POST /", dep.ware.HasRoleAndPermissions(m.RequestBody(http.HandlerFunc(dep.create)), rp))
	dep.mux.Handle("/service", dep.ware.Authentication(dep.ware.HasRole(mux, &staff)))
}

func (dep *ServiceHandler) create(w http.ResponseWriter, r *http.Request) {
	dto, err := pkg.ReadBody[model.ServicePayload](r)
	if err != nil {
		dep.logger.Error(err.Error())
		utils.ConstructErrorResponse(w, &utils.BadRequestError{Message: err.Error()})
		return
	}

	if err = dep.service.Create(r.Context(), dto); err != nil {
		dep.logger.Error(err.Error())
		utils.ConstructErrorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

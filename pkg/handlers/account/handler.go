package account

import (
	"github.com/iTchTheRightSpot/erp-golang/pkg"
	"github.com/iTchTheRightSpot/erp-golang/pkg/middleware"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/profile"
	"github.com/iTchTheRightSpot/erp-golang/pkg/services/account"
	"github.com/iTchTheRightSpot/erp-golang/pkg/services/auth"
	"github.com/iTchTheRightSpot/erp-golang/utils"
	"net/http"
)

type AccountHandler struct {
	mux     *http.ServeMux
	ware    *middleware.Middleware
	logger  utils.ILogger
	ps      auth.IPasswordService
	service account.IAccountService
}

func NewAccountHandler(mux *http.ServeMux, w *middleware.Middleware, l utils.ILogger, ps auth.IPasswordService, s account.IAccountService) *AccountHandler {
	return &AccountHandler{mux: mux, ware: w, logger: l, ps: ps, service: s}
}

func (dep *AccountHandler) Register() {
	ware := middleware.RequestBodyMiddleware[profile.ProfilePayload]{Logger: dep.logger}
	rp := &models.RolePermission{
		Role:        models.STAFF,
		Permissions: []models.PermissionEnum{models.WRITE},
	}

	dep.mux.Handle("POST /account/register", dep.ware.Authentication(dep.ware.HasRoleAndPermissions(ware.RequestBody(http.HandlerFunc(dep.register)), rp)))
}

func (dep *AccountHandler) register(w http.ResponseWriter, r *http.Request) {
	p, err := pkg.ReadBody[profile.ProfilePayload](r)
	if err != nil {
		dep.logger.Error(err.Error())
		utils.ErrorResponse(w, &utils.BadRequestError{Message: err.Error()})
		return
	}

	if err = dep.ps.PasswordRegex(p.Password); err != nil {
		utils.ErrorResponse(w, &utils.BadRequestError{Message: err.Error()})
		return
	}

	if err = dep.service.Register(r.Context(), p); err != nil {
		dep.logger.Error(err.Error())
		utils.ErrorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

package account

import (
	"github.com/iTchTheRightSpot/erp-golang/config"
	"github.com/iTchTheRightSpot/erp-golang/pkg"
	"github.com/iTchTheRightSpot/erp-golang/pkg/middleware"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models"
	"github.com/iTchTheRightSpot/erp-golang/pkg/services/account"
	"github.com/iTchTheRightSpot/erp-golang/pkg/services/auth"
	"github.com/iTchTheRightSpot/erp-golang/utils"
	"net/http"
)

type AccountHandler struct {
	mux     *http.ServeMux
	ware    *middleware.Middleware
	logger  utils.ILogger
	env     *config.SecretVariables
	ps      auth.IPasswordService
	service account.IAccountService
}

func NewAccountHandler(mux *http.ServeMux, w *middleware.Middleware, l utils.ILogger, env *config.SecretVariables, ps auth.IPasswordService, s account.IAccountService) *AccountHandler {
	return &AccountHandler{mux: mux, ware: w, logger: l, env: env, ps: ps, service: s}
}

func (dep *AccountHandler) Register() {
	rp := &models.RolePermissionEnum{
		Role:        models.STAFF,
		Permissions: []models.PermissionEnum{models.WRITE},
	}

	ware1 := middleware.RequestBodyMiddleware[models.ProfilePayload]{Logger: dep.logger}
	dep.mux.Handle("POST /account/register", dep.ware.Authentication(dep.ware.HasRoleAndPermissions(ware1.RequestBody(http.HandlerFunc(dep.register)), rp)))

	ware2 := middleware.RequestBodyMiddleware[models.Login]{Logger: dep.logger}
	dep.mux.Handle("POST /account/login", ware2.RequestBody(http.HandlerFunc(dep.login)))
}

func (dep *AccountHandler) register(w http.ResponseWriter, r *http.Request) {
	p, err := pkg.ReadBody[models.ProfilePayload](r)
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

func (dep *AccountHandler) login(w http.ResponseWriter, r *http.Request) {
	p, err := pkg.ReadBody[models.Login](r)
	if err != nil {
		dep.logger.Error(err.Error())
		utils.ErrorResponse(w, &utils.BadRequestError{Message: err.Error()})
		return
	}

	o, err := dep.service.Login(r.Context(), p)
	if err != nil {
		utils.ErrorResponse(w, err)
		return
	}

	pkg.WriteCookie(w, dep.env.CookieParam, o.Token, o.ExpireAt)
	w.WriteHeader(http.StatusCreated)
}

package account

import (
	"encoding/json"
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
	// public
	m1 := middleware.RequestBodyMiddleware[models.Login]{Logger: dep.logger}
	dep.mux.Handle("POST /account/login", m1.RequestBody(http.HandlerFunc(dep.login)))

	// protected
	// read
	dep.mux.Handle("GET /active", dep.ware.Authentication(http.HandlerFunc(dep.activeUser)))

	// write
	dep.mux.Handle("POST /logout", dep.ware.Authentication(http.HandlerFunc(dep.logout)))
	m2 := middleware.RequestBodyMiddleware[models.ProfilePayload]{Logger: dep.logger}
	dep.mux.Handle("POST /account/register", dep.ware.Authentication(
		dep.ware.HasRoleAndPermissions(
			m2.RequestBody(http.HandlerFunc(dep.register)),
			&models.RolePermissionEnum{Role: models.STAFF, Permissions: []models.PermissionEnum{models.WRITE}},
		),
	))

	m3 := middleware.RequestBodyMiddleware[models.AddRoleAndPermissionPayload]{Logger: dep.logger}
	dep.mux.Handle("POST /account/role-permission", dep.ware.Authentication(
		dep.ware.HasRoleAndPermissions(
			m3.RequestBody(http.HandlerFunc(dep.addRoleAndPermission)),
			&models.RolePermissionEnum{Role: models.STAFF, Permissions: []models.PermissionEnum{models.WRITE}},
		),
	))

	dep.mux.Handle("DELETE /account/role", dep.ware.Authentication(
		dep.ware.HasRoleAndPermissions(
			http.HandlerFunc(dep.deleteRole),
			&models.RolePermissionEnum{Role: models.STAFF, Permissions: []models.PermissionEnum{models.DELETE}},
		),
	))
	dep.mux.Handle("DELETE /account/permission/{permission_id}", dep.ware.Authentication(
		dep.ware.HasRoleAndPermissions(
			http.HandlerFunc(dep.deletePermission),
			&models.RolePermissionEnum{Role: models.STAFF, Permissions: []models.PermissionEnum{models.DELETE}},
		),
	))
}

func (dep *AccountHandler) activeUser(w http.ResponseWriter, r *http.Request) {
	obj, ok := r.Context().Value(utils.UserContextKey).(*models.JwtObj)
	if !ok || obj == nil {
		dep.logger.Error("activeUser: invalid user context")
		utils.ErrorResponse(w, &utils.AuthenticationError{})
		return
	}

	a, err := dep.service.ActiveUser(r.Context(), obj)
	if err != nil {
		utils.ErrorResponse(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(a); err != nil {
		dep.logger.Error(err.Error())
	}
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

func (dep *AccountHandler) logout(w http.ResponseWriter, r *http.Request) {
	if r.Cookies() == nil || len(r.Cookies()) == 0 {
		dep.logger.Error("cookie present")
		utils.ErrorResponse(w, &utils.AuthenticationError{})
		return
	}

	cookie, err := r.Cookie(dep.env.CookieParam.CookieName)
	if err != nil || cookie == nil {
		dep.logger.Error(err)
		utils.ErrorResponse(w, &utils.AuthenticationError{Message: "invalid cookie"})
		return
	}

	// max-age takes higher precedence to expires https://stackoverflow.com/questions/70715904/golang-cookie-max-age-vs-expire
	cookie.MaxAge = -1 // expire now
	cookie.Path = "/"
	cookie.HttpOnly = true
	cookie.SameSite = dep.env.CookieParam.SameSite
	cookie.Secure = dep.env.CookieParam.CookieSecure

	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusNoContent)
}

func (dep *AccountHandler) addRoleAndPermission(w http.ResponseWriter, r *http.Request) {
	p, err := pkg.ReadBody[models.AddRoleAndPermissionPayload](r)
	if err != nil {
		dep.logger.Error(err.Error())
		utils.ErrorResponse(w, &utils.BadRequestError{Message: err.Error()})
		return
	}

	if err = dep.service.AddRoleAndPermission(r.Context(), p); err != nil {
		utils.ErrorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (dep *AccountHandler) deleteRole(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (dep *AccountHandler) deletePermission(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

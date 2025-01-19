package account

import (
	"github.com/iTchTheRightSpot/erp-golang/pkg"
	"github.com/iTchTheRightSpot/erp-golang/pkg/middleware"
	"github.com/iTchTheRightSpot/erp-golang/pkg/models/profile"
	"github.com/iTchTheRightSpot/erp-golang/pkg/services/account"
	"github.com/iTchTheRightSpot/erp-golang/utils"
	"net/http"
)

type AccountHandler struct {
	mux     *http.ServeMux
	logger  utils.ILogger
	service account.IAccountService
}

func NewAccountHandler(mux *http.ServeMux, l utils.ILogger, s account.IAccountService) *AccountHandler {
	return &AccountHandler{mux: mux, logger: l, service: s}
}

func (dep *AccountHandler) Register() {
	ware := middleware.RequestBodyMiddleware[profile.ProfilePayload]{Logger: dep.logger}

	dep.mux.Handle("POST /account/register", ware.RequestBody(http.HandlerFunc(dep.register)))
}

func (dep *AccountHandler) register(w http.ResponseWriter, r *http.Request) {
	p, err := pkg.ReadBody[profile.ProfilePayload](r)
	if err != nil {
		dep.logger.Error(err.Error())
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

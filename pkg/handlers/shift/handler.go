package shift

import (
	"github.com/iTchTheRightSpot/erp-golang/pkg/middleware"
	shiftModel "github.com/iTchTheRightSpot/erp-golang/pkg/models/shift"
	"github.com/iTchTheRightSpot/erp-golang/pkg/services/shift"
	"github.com/iTchTheRightSpot/erp-golang/utils"
	"net/http"
)

type ShiftHandler struct {
	mux        *http.ServeMux
	middleware middleware.Middleware[any]
	logger     utils.ILogger
	service    shift.IShiftService
}

func NewShiftHandler(mux *http.ServeMux, w middleware.Middleware[any], l utils.ILogger, s shift.IShiftService) *ShiftHandler {
	return &ShiftHandler{mux: mux, middleware: w, logger: l, service: s}
}

func (dep *ShiftHandler) RegisterRoutes() {
	m := middleware.RequestBodyMiddleware[shiftModel.ShiftDto]{}
	dep.mux.Handle("POST /shift", m.RequestBody(dep.middleware.Permission(dep.middleware.Authentication(http.HandlerFunc(dep.create)))))
}

func (dep *ShiftHandler) create(w http.ResponseWriter, r *http.Request) {

}

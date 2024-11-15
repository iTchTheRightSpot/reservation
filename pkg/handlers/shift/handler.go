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
	middleware *middleware.Middleware
	logger     utils.ILogger
	service    shift.IShiftService
}

func NewShiftHandler(mux *http.ServeMux, w *middleware.Middleware, l utils.ILogger, s shift.IShiftService) *ShiftHandler {
	return &ShiftHandler{mux: mux, middleware: w, logger: l, service: s}
}

func (dep *ShiftHandler) RegisterRoutes() {
	m := middleware.RequestBodyMiddleware[shiftModel.ShiftDto]{Logger: dep.logger}
	next := http.HandlerFunc(dep.create)
	dep.mux.Handle("POST /shift", dep.middleware.ChainAuth(m.RequestBody(next)))
}

func (dep *ShiftHandler) create(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("shift handler post route hit"))
}

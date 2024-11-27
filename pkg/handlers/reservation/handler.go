package reservation

import (
	"github.com/iTchTheRightSpot/erp-golang/pkg"
	"github.com/iTchTheRightSpot/erp-golang/pkg/middleware"
	model "github.com/iTchTheRightSpot/erp-golang/pkg/models/reservation"
	"github.com/iTchTheRightSpot/erp-golang/pkg/services/reservation"
	"github.com/iTchTheRightSpot/erp-golang/utils"
	"net/http"
)

type ReservationHandler struct {
	mux     *http.ServeMux
	logger  utils.ILogger
	service reservation.IReservationService
}

func NewReservationHandler(mux *http.ServeMux, l utils.ILogger, s reservation.IReservationService) *ReservationHandler {
	return &ReservationHandler{mux: mux, logger: l, service: s}
}

func (dep *ReservationHandler) Register() {
	mux := http.NewServeMux()
	m := middleware.RequestBodyMiddleware[model.ReservationPayload]{Logger: dep.logger}

	mux.Handle("POST /", m.RequestBody(http.HandlerFunc(dep.create)))
	mux.HandleFunc("GET /", dep.availableDates)

	dep.mux.Handle("/reservation", mux)
}

func (dep *ReservationHandler) create(w http.ResponseWriter, r *http.Request) {
	dto, err := pkg.ReadBody[model.ReservationPayload](r)
	if err != nil {
		dep.logger.Error(err.Error())
		utils.ConstructErrorResponse(
			w,
			utils.ErrorResponse{
				Status:  http.StatusInternalServerError,
				Message: err.Error(),
			},
		)
		return
	}

	if err = dep.service.Create(r.Context(), dto); err != nil {
		dep.logger.Error(err.Error())
		utils.ConstructErrorResponse(
			w,
			utils.ErrorResponse{
				Status:  http.StatusBadRequest,
				Message: err.Error(),
			},
		)
		return
	}

	dep.logger.Log("new reservation")
	w.WriteHeader(http.StatusCreated)
}

func (dep *ReservationHandler) availableDates(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

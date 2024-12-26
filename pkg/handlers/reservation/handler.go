package reservation

import (
	"encoding/json"
	"github.com/iTchTheRightSpot/erp-golang/pkg"
	"github.com/iTchTheRightSpot/erp-golang/pkg/middleware"
	model "github.com/iTchTheRightSpot/erp-golang/pkg/models/reservation"
	"github.com/iTchTheRightSpot/erp-golang/pkg/services/reservation"
	"github.com/iTchTheRightSpot/erp-golang/utils"
	"net/http"
	"strconv"
	"strings"
	"time"
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
	ware := middleware.RequestBodyMiddleware[model.ReservationPayload]{Logger: dep.logger}
	dep.mux.HandleFunc("GET /reservation", dep.availableDates)
	dep.mux.Handle("POST /reservation", ware.RequestBody(http.HandlerFunc(dep.create)))
	dep.mux.HandleFunc("POST /reservation/cancel/{reservation_id}", dep.cancel)
}

func (dep *ReservationHandler) create(w http.ResponseWriter, r *http.Request) {
	dto, err := pkg.ReadBody[model.ReservationPayload](r)
	if err != nil {
		dep.logger.Error(err.Error())
		utils.ErrorResponse(w, err)
		return
	}

	err = dep.service.Create(r.Context(), dto)
	if err != nil {
		dep.logger.Error(err.Error())
		utils.ErrorResponse(w, err)
		return
	}

	dep.logger.Log("reservation created")
	w.WriteHeader(http.StatusCreated)
}

func (dep *ReservationHandler) availableDates(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	day, err := strconv.Atoi(query.Get("day"))
	if err != nil {
		dep.logger.Error("day missing err: ", err.Error())
		utils.ErrorResponse(w, &utils.BadRequestError{Message: "day missing"})
		return
	}

	month, err := strconv.Atoi(query.Get("month"))
	if err != nil {
		dep.logger.Error("month missing err: ", err.Error())
		utils.ErrorResponse(w, &utils.BadRequestError{Message: "month missing"})
		return
	}

	year, err := strconv.Atoi(query.Get("year"))
	if err != nil {
		dep.logger.Error("year missing err: ", err.Error())
		utils.ErrorResponse(w, &utils.BadRequestError{Message: "year missing"})
		return
	}

	p := model.AvailableTimesPayload{
		StaffId:  strings.TrimSpace(query.Get("staff_id")),
		Services: query["service"],
		Day:      day,
		Month:    month,
		Year:     year,
	}

	if err = middleware.ValidatorInstance.Struct(p); err != nil {
		dep.logger.Error(err.Error())
		utils.ErrorResponse(w, &utils.BadRequestError{Message: err.Error()})
		return
	}

	location := dep.logger.Timezone()
	zone := query.Get("timezone")
	if zone != "" {
		location, err = time.LoadLocation(zone)
		if err != nil {
			dep.logger.Error(err.Error())
			utils.ErrorResponse(w, &utils.BadRequestError{Message: err.Error()})
			return
		}
	}

	p.StartDateTime = time.Date(p.Year, time.Month(p.Month), day, 0, 0, 0, 0, location)
	p.EndDateTime = p.StartDateTime.AddDate(0, 1, -1)

	arr, err := dep.service.AvailableDates(r.Context(), &p)
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

func (dep *ReservationHandler) cancel(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("reservation_id"), 10, 64)
	if err != nil {
		dep.logger.Error(err.Error())

		utils.ErrorResponse(w, &utils.BadRequestError{Message: "invalid reservation id"})
		return
	}

	if err = dep.service.Cancel(r.Context(), uint64(id)); err != nil {
		utils.ErrorResponse(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	if _, err = w.Write([]byte("reservation cancelled")); err != nil {
		dep.logger.Error(err.Error())
	}
}
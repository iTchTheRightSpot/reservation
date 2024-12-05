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
		utils.ErrorResponse(w, err)
		return
	}

	err = dep.service.Create(r.Context(), dto)
	if err != nil {
		dep.logger.Error(err.Error())
		utils.ErrorResponse(w, err)
		return
	}

	dep.logger.Log("created a new reservation")
	w.WriteHeader(http.StatusCreated)
}

func (dep *ReservationHandler) availableDates(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	day, err := strconv.Atoi(query.Get("day"))
	if err != nil {
		dep.logger.Error(err.Error())
		utils.ErrorResponse(w, &utils.BadRequestError{Message: err.Error()})
		return
	}

	month, err := strconv.Atoi(query.Get("month"))
	if err != nil {
		dep.logger.Error(err.Error())
		utils.ErrorResponse(w, &utils.BadRequestError{Message: err.Error()})
		return
	}

	year, err := strconv.Atoi(query.Get("year"))
	if err != nil {
		dep.logger.Error(err.Error())
		utils.ErrorResponse(w, &utils.BadRequestError{Message: err.Error()})
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
	if query.Get("timezone") != "" {
		location, err = time.LoadLocation(query.Get("timezone"))
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

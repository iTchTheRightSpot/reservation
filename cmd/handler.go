package cmd

import (
	"database/sql"
	"github.com/iTchTheRightSpot/erp-golang/config"
	"github.com/iTchTheRightSpot/erp-golang/pkg/handlers/reservation"
	"github.com/iTchTheRightSpot/erp-golang/pkg/handlers/schedule"
	"github.com/iTchTheRightSpot/erp-golang/pkg/handlers/service"
	"github.com/iTchTheRightSpot/erp-golang/pkg/handlers/staff"
	"github.com/iTchTheRightSpot/erp-golang/pkg/middleware"
	"github.com/iTchTheRightSpot/erp-golang/utils"
	"net/http"
	"strings"
)

type HandlerRegistry struct {
	logger   utils.ILogger
	env      *config.SecretVariables
	mux      *http.ServeMux
	services *serviceRegistry
	ware     *middleware.Middleware
}

func NewHandlerRegistry(mux *http.ServeMux, db *sql.DB, l utils.ILogger, e *config.SecretVariables) *HandlerRegistry {
	s := newServiceRegistry(db, l, e)
	return &HandlerRegistry{
		logger:   l,
		env:      e,
		mux:      mux,
		services: s,
		ware:     &middleware.Middleware{Logger: l, Auth: s.JwtService, Param: e.CookieParam},
	}
}

func (dep *HandlerRegistry) prefix() string {
	r := dep.env.RoutePrefix
	var prefix strings.Builder
	if string(r[len(r)-1]) == "/" {
		prefix.WriteString(r[0 : len(r)-1])
	} else {
		prefix.WriteString(r)
	}
	return prefix.String()
}

func (dep *HandlerRegistry) Initialize() http.Handler {
	v1 := http.NewServeMux()

	// register handlers
	schedule.NewScheduleHandler(v1, dep.ware, dep.logger, dep.services.ScheduleService).Register()
	service.NewServiceHandler(v1, dep.logger, dep.services.ServiceImpl, dep.ware).Register()
	staff.NewStaffHandler(v1, dep.ware, dep.logger, dep.services.StaffService).Register()
	reservation.NewReservationHandler(v1, dep.logger, dep.services.ReservationService).Register()

	// register v1 with mux
	dep.mux.Handle(dep.env.RoutePrefix, http.StripPrefix(dep.prefix(), v1))

	return dep.ware.Initialize(dep.mux)
}

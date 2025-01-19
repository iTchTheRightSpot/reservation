package cmd

import (
	"database/sql"
	"github.com/iTchTheRightSpot/erp-golang/config"
	"github.com/iTchTheRightSpot/erp-golang/pkg/handlers/account"
	"github.com/iTchTheRightSpot/erp-golang/pkg/handlers/reservation"
	"github.com/iTchTheRightSpot/erp-golang/pkg/handlers/schedule"
	"github.com/iTchTheRightSpot/erp-golang/pkg/handlers/service_type"
	"github.com/iTchTheRightSpot/erp-golang/pkg/handlers/staff"
	"github.com/iTchTheRightSpot/erp-golang/pkg/middleware"
	"github.com/iTchTheRightSpot/erp-golang/utils"
	"net/http"
)

type HandlerRegistry struct {
	lg   utils.ILogger
	env  *config.SecretVariables
	mux  *http.ServeMux
	ss   *serviceRegistry
	ware *middleware.Middleware
}

func NewHandlerRegistry(db *sql.DB, l utils.ILogger, e *config.SecretVariables) *HandlerRegistry {
	s := newServiceRegistry(db, l, e)
	return &HandlerRegistry{
		lg:   l,
		env:  e,
		mux:  http.NewServeMux(),
		ss:   s,
		ware: &middleware.Middleware{Logger: l, Auth: s.JwtService, Param: e.CookieParam},
	}
}

func (dep *HandlerRegistry) Initialize() http.Handler {
	v1 := http.NewServeMux()

	// register handlers
	schedule.NewScheduleHandler(v1, dep.ware, dep.lg, dep.ss.ScheduleService).Register()
	service_type.NewServiceTypeHandler(v1, dep.lg, dep.ss.ServiceImpl, dep.ware).Register()
	staff.NewStaffHandler(v1, dep.ware, dep.lg, dep.ss.StaffService).Register()
	reservation.NewReservationHandler(v1, dep.lg, dep.ss.ReservationService).Register()
	account.NewAccountHandler(v1, dep.ware, dep.lg, dep.ss.PasswordService, dep.ss.AccountService)

	dep.mux.Handle("/api/v1/", http.StripPrefix("/api/v1", v1))

	return dep.ware.Initialize(dep.mux)
}

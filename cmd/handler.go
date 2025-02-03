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
	f    http.FileSystem
}

func NewHandlerRegistry(db *sql.DB, l utils.ILogger, e *config.SecretVariables, f http.FileSystem) *HandlerRegistry {
	s := newServiceRegistry(db, l, e)
	return &HandlerRegistry{
		f:    f,
		lg:   l,
		env:  e,
		mux:  http.NewServeMux(),
		ss:   s,
		ware: &middleware.Middleware{Logger: l, Auth: s.JwtService, Param: e.CookieParam, ApiPrefix: e.ApiPrefix, FileSystem: f},
	}
}

func (dep *HandlerRegistry) Initialize() http.Handler {
	v1 := http.NewServeMux()

	schedule.NewScheduleHandler(v1, dep.ware, dep.lg, dep.ss.ScheduleService).Register()
	service_type.NewServiceTypeHandler(v1, dep.lg, dep.ss.ServiceImpl, dep.ware).Register()
	staff.NewStaffHandler(v1, dep.ware, dep.lg, dep.ss.StaffService).Register()
	reservation.NewReservationHandler(v1, dep.lg, dep.ware, dep.ss.ReservationService).Register()
	account.NewAccountHandler(v1, dep.ware, dep.lg, dep.env, dep.ss.PasswordService, dep.ss.AccountService).Register()

	dep.mux.Handle(dep.env.ApiPrefix, http.StripPrefix(dep.env.ApiPrefix[:len(dep.env.ApiPrefix)-1], v1))

	dep.mux.Handle("/", http.StripPrefix("/", http.FileServer(dep.f)))

	return dep.ware.Initialize(dep.mux)
}

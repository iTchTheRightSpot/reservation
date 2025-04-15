package cmd

import (
	"database/sql"
	"github.com/iTchTheRightSpot/reservation/config"
	"github.com/iTchTheRightSpot/reservation/pkg/handlers/account"
	"github.com/iTchTheRightSpot/reservation/pkg/handlers/reservation"
	"github.com/iTchTheRightSpot/reservation/pkg/handlers/schedule"
	"github.com/iTchTheRightSpot/reservation/pkg/handlers/service_type"
	"github.com/iTchTheRightSpot/reservation/pkg/handlers/staff"
	"github.com/iTchTheRightSpot/reservation/pkg/middleware"
	"github.com/iTchTheRightSpot/utility/utils"
	"net/http"
)

type HandlerRegistry struct {
	lg   utils.ILogger
	env  *config.SecretVariables
	mux  *http.ServeMux
	sr   *serviceRegistry
	ware *middleware.Middleware
	f    http.FileSystem
}

func NewHandlerRegistry(db *sql.DB, l utils.ILogger, e *config.SecretVariables, f http.FileSystem) *HandlerRegistry {
	s := newServiceRegistry(db, l, e)
	ware := middleware.Middleware{Logger: l, Auth: s.JwtService, Param: e.CookieParam, ApiPrefix: e.ApiPrefix, FS: f}
	return &HandlerRegistry{f: f, lg: l, env: e, mux: http.NewServeMux(), sr: s, ware: &ware}
}

func (dep *HandlerRegistry) Initialize() http.Handler {
	v1 := http.NewServeMux()

	schedule.NewScheduleHandler(v1, dep.ware, dep.lg, dep.sr.ScheduleService).Register()
	service_type.NewServiceTypeHandler(v1, dep.lg, dep.sr.ServiceImpl, dep.ware).Register()
	staff.NewStaffHandler(v1, dep.ware, dep.lg, dep.sr.StaffService).Register()
	reservation.NewReservationHandler(v1, dep.lg, dep.ware, dep.sr.ReservationService).Register()
	account.NewAccountHandler(v1, dep.ware, dep.lg, dep.env, dep.sr.PasswordService, dep.sr.AccountService).Register()

	dep.mux.Handle(dep.env.ApiPrefix, http.StripPrefix(dep.env.ApiPrefix[:len(dep.env.ApiPrefix)-1], v1))

	dep.mux.Handle("/", http.StripPrefix("/", http.FileServer(dep.f)))

	return dep.ware.Initialize(dep.mux)
}
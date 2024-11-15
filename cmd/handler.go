package cmd

import (
	"database/sql"
	"github.com/iTchTheRightSpot/erp-golang/config"
	"github.com/iTchTheRightSpot/erp-golang/pkg/handlers/shift"
	"github.com/iTchTheRightSpot/erp-golang/pkg/middleware"
	"github.com/iTchTheRightSpot/erp-golang/utils"
	"net/http"
)

type HandlerRegistry struct {
	log        utils.ILogger
	env        *config.SecretVariables
	mux        *http.ServeMux
	services   *serviceRegistry
	middleware *middleware.Middleware
}

func NewHandlerRegistry(mux *http.ServeMux, db *sql.DB, l utils.ILogger, e *config.SecretVariables) *HandlerRegistry {
	s := newServiceRegistry(db, l, e)
	m := &middleware.Middleware{Logger: l, Auth: s.JwtService, Param: e.CookieParam}
	return &HandlerRegistry{
		log:        l,
		env:        e,
		mux:        mux,
		services:   s,
		middleware: m,
	}
}

func (dep *HandlerRegistry) Initialize() http.Handler {
	// v1 subroutine
	v1 := http.NewServeMux()

	// register handlers
	shift.NewShiftHandler(v1, dep.middleware, dep.log, dep.services.ShiftService).RegisterRoutes()

	// register v1 with mux
	dep.mux.Handle("/api/v1/", http.StripPrefix("/api/v1", v1))

	return dep.middleware.Initialize(dep.mux)
}

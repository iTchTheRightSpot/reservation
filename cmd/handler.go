package cmd

import (
	"database/sql"
	"github.com/iTchTheRightSpot/erp-golang/config"
	"github.com/iTchTheRightSpot/erp-golang/pkg/handlers/shift"
	"github.com/iTchTheRightSpot/erp-golang/pkg/middleware"
	"github.com/iTchTheRightSpot/erp-golang/utils"
	"net/http"
	"strings"
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
	return &HandlerRegistry{
		log:        l,
		env:        e,
		mux:        mux,
		services:   s,
		middleware: &middleware.Middleware{Logger: l, Auth: s.JwtService, Param: e.CookieParam},
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
	shift.NewShiftHandler(v1, dep.middleware, dep.log, dep.services.ShiftService).Register()

	// register v1 with mux
	dep.mux.Handle(dep.env.RoutePrefix, http.StripPrefix(dep.prefix(), v1))

	return dep.middleware.Initialize(dep.mux)
}

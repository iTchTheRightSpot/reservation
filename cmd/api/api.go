package api

import (
	"database/sql"
	"github.com/iTchTheRightSpot/erp-golang/cmd"
	"github.com/iTchTheRightSpot/erp-golang/config"
	"github.com/iTchTheRightSpot/erp-golang/utils"
	"net/http"
	"time"
)

type ErpServer struct {
	Db     *sql.DB
	Logger utils.ILogger
	Env    *config.SecretVariables
}

func (s *ErpServer) Serve() {
	s.Logger.Log("starting server on PORT", s.Env.Address)

	server := http.Server{
		Addr:              s.Env.Address,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       5 * time.Second,
		Handler:           cmd.NewHandlerRegistry(s.Db, s.Logger, s.Env).Initialize(),
	}

	s.Logger.Fatal("server stopped ", server.ListenAndServe())
}

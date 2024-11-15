package api

import (
	"database/sql"
	"fmt"
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
	a := cmd.NewHandlerRegistry(http.NewServeMux(), s.Db, s.Logger, s.Env)
	server := http.Server{
		Addr:              s.Env.Address,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       5 * time.Second,
		Handler:           a.Initialize(),
	}

	s.Logger.Log(fmt.Sprintf("starting server on PORT %s", server.Addr))

	s.Logger.Fatal("server failed to start: %v", server.ListenAndServe())
}

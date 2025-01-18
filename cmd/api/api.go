package api

import (
	"database/sql"
	"github.com/iTchTheRightSpot/erp-golang/cmd"
	"github.com/iTchTheRightSpot/erp-golang/config"
	"github.com/iTchTheRightSpot/erp-golang/utils"
	"github.com/rs/cors"
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

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{s.Env.FrontEnd},
		AllowedMethods:   []string{http.MethodPost, http.MethodGet, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowedHeaders:   []string{"Content-Type", "Accept", "Cookie"},
		AllowCredentials: true,
	})

	registry := cmd.NewHandlerRegistry(s.Db, s.Logger, s.Env)

	server := http.Server{
		Addr:              s.Env.Address,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       5 * time.Second,
		Handler:           c.Handler(registry.Initialize()),
	}

	s.Logger.Fatal("server stopped ", server.ListenAndServe())
}

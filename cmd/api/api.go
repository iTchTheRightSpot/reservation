package api

import (
	"context"
	"database/sql"
	"github.com/iTchTheRightSpot/erp-golang/cmd"
	"github.com/iTchTheRightSpot/erp-golang/config"
	"github.com/iTchTheRightSpot/erp-golang/frontend"
	"github.com/iTchTheRightSpot/utility/utils"
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
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{s.Env.FrontEnd},
		AllowedMethods:   []string{http.MethodPost, http.MethodGet, http.MethodPut, http.MethodDelete},
		AllowedHeaders:   []string{"Origin", "Content-Type", "Accept", "Cookie"},
		ExposedHeaders:   []string{"Content-Length"},
		AllowCredentials: true,
	})

	f := frontend.FrontendStruct{Logger: s.Logger}
	open, err := f.FileSystem()
	if err != nil {
		s.Logger.Fatal(err.Error())
	}

	server := http.Server{
		Addr:              s.Env.Address,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       5 * time.Second,
		Handler:           c.Handler(cmd.NewHandlerRegistry(s.Db, s.Logger, s.Env, open).Initialize()),
	}

	s.Logger.Log(context.Background(), "starting server on PORT", s.Env.Address)
	s.Logger.Fatal("server stopped ", server.ListenAndServe())
}
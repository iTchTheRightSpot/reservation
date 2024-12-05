package main

import (
	"github.com/iTchTheRightSpot/erp-golang/cmd/api"
	"github.com/iTchTheRightSpot/erp-golang/config"
	"github.com/iTchTheRightSpot/erp-golang/database"
	"github.com/iTchTheRightSpot/erp-golang/utils"
	"log"
)

func main() {
	// logger
	logger, err := utils.NewLogger("UTC")
	if err != nil {
		return
	}

	// config
	obj := config.SecretVariables{}
	env, err := obj.Config()
	if err != nil {
		return
	}

	// connect to database
	db, err := database.ConnectToPostgres(env.DbConnectionString)
	if err != nil {
		log.Fatal(err)
		return
	}

	// start server
	server := api.ErpServer{Db: db, Logger: logger, Env: env}
	server.Serve()
}

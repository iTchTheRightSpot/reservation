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
		log.Fatal("failed to instantiate logger ", err.Error())
	}

	// config
	obj := config.SecretVariables{}
	env := obj.Config()

	// connect to database
	db, err := database.ConnectToPostgres(env.DbConnectionString)
	if err != nil {
		logger.Fatal(err.Error())
	}

	if env.CookieParam.CookieSecure {
		if err = database.Migrate(db); err != nil {
			logger.Fatal(err.Error())
		}
	}

	// start server
	server := api.ErpServer{Db: db, Logger: logger, Env: env}
	server.Serve()
}

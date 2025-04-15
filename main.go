package main

import (
	"context"
	"database/sql"
	"github.com/iTchTheRightSpot/erp-golang/cmd/api"
	"github.com/iTchTheRightSpot/erp-golang/config"
	"github.com/iTchTheRightSpot/erp-golang/database"
	"github.com/iTchTheRightSpot/utility/utils"
)

func main() {
	// initialize env
	obj := config.SecretVariables{}
	env := obj.Config()

	lg := utils.DevLogger("UTC")
	if env.Profile == "production" {
		lg = utils.ProdLogger("UTC", env.Discord)
	}

	db, err := database.ConnectToPostgres(env.DbConnectionString)
	if err != nil {
		lg.Fatal(err.Error())
	}

	defer func(db *sql.DB) {
		if err = db.Close(); err != nil {
			lg.Error(context.Background(), "error defer closing db connection "+err.Error())
		}
	}(db)

	if env.CookieParam.CookieSecure {
		if err = database.Migrate(db); err != nil {
			lg.Fatal(err.Error())
		}
	}

	// start server
	server := api.ErpServer{Db: db, Logger: lg, Env: env}
	server.Serve()
}
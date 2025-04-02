package main

import (
	"database/sql"
	"github.com/iTchTheRightSpot/erp-golang/cmd/api"
	"github.com/iTchTheRightSpot/erp-golang/config"
	"github.com/iTchTheRightSpot/erp-golang/database"
	"github.com/iTchTheRightSpot/erp-golang/utils"
)

func main() {
	// initialize env
	obj := config.SecretVariables{}
	env := obj.Config()

	// logger
	var l = utils.NewDevLogger()
	if env.Profile == "production" {
		lg, err := utils.NewLogger("UTC", env.Discord)
		if err != nil {
			l.Fatal("failed to instantiate logger ", err.Error())
		}
		l = lg
	}

	// connect to database
	db, err := database.ConnectToPostgres(env.DbConnectionString)
	if err != nil {
		l.Fatal(err.Error())
	}

	 defer func(db *sql.DB) {
                if err = db.Close(); err != nil {
                        l.Error("error defer closing db connection " + err.Error())
                }
        }(db)

	if env.CookieParam.CookieSecure {
		if err = database.Migrate(db); err != nil {
			l.Fatal(err.Error())
		}
	}

	// start server
	server := api.ErpServer{Db: db, Logger: l, Env: env}
	server.Serve()
}

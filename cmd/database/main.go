package main

import (
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/iTchTheRightSpot/reservation/config"
	"github.com/iTchTheRightSpot/reservation/database"
	"log"
)

func main() {
	secret := config.SecretVariables{}
	env := secret.Config()

	db, err := database.ConnectToPostgres(env.DbConnectionString)
	if err != nil {
		log.Fatal(err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://./database/migrations/", "postgres", driver)
	if err != nil {
		log.Fatal(err)
	}

	if err = m.Up(); err != nil {
		log.Fatal(err)
	}
}
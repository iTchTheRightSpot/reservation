package database

import (
	"database/sql"
	"embed"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/lib/pq"
	"log"
	"time"
)

//go:embed migrations/*.sql
var files embed.FS

// ConnectToPostgres https://go.dev/doc/tutorial/database-access
func ConnectToPostgres(conn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", conn)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	// https://www.alexedwards.net/blog/configuring-sqldb
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err = db.Ping(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	log.Println("database connection established")
	return db, nil
}

func Migrate(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	d, err := iofs.New(files, "migrations")
	if err != nil {
		log.Fatal(err)
		return err
	}

	m, err := migrate.NewWithInstance("iofs", d, "postgres", driver)
	if err != nil {
		return err
	}

	if err = m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("no new migrations to apply")
		} else {
			return err
		}
	}

	return nil
}

package database

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"time"
)

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

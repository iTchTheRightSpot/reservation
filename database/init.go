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

//func ConnectToPostgres(conn string) (*pgxpool.Pool, error) {
//	config, err := pgxpool.ParseConfig(conn)
//	if err != nil {
//		log.Fatalf("unable to parse connection string: %v", err)
//	}
//
//	config.MaxConns = 10                     // equivalent to SetMaxOpenConns
//	config.MinConns = 5                      // equivalent to SetMaxIdleConns
//	config.MaxConnLifetime = 5 * time.Minute // equivalent to SetConnMaxLifetime
//	//config.MaxConnIdleTime = 30 * time.Minute // set idle timeout
//
//	pool, err := pgxpool.NewWithConfig(context.Background(), config)
//	if err != nil {
//		log.Fatalf("unable to create connection pool: %v", err)
//	}
//	defer pool.Close()
//
//	log.Println("database connection established")
//	return pool, nil
//}

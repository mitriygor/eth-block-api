package db

import (
	"database/sql"
	"eth-listener/config"
	"log"
	"os"
	"time"
)

var counts int64

var db *sql.DB

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func ConnectToDB(DSN string) *sql.DB {
	dsn := os.Getenv(DSN)

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres not yet ready ...")
			counts++
		} else {
			log.Println("Connected to Postgres!")
			return connection
		}

		if counts > config.MaxAttempts {
			log.Println(err)
			return nil
		}

		log.Println("Waiting for two seconds....")
		time.Sleep(2 * time.Second)
		continue
	}
}

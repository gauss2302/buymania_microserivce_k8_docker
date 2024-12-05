package db

import (
	"database/sql"
	"fmt"
	"log"
)

func NewPostgresDB(host, port, user, password, dbname string) (*sql.DB, error) {
	// Connect to the db
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("error opening db: %w", err)
	}

	// Check if the db is alive
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("error pinging db: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)

	log.Println("Successfully connected to the db")
	return db, nil
}

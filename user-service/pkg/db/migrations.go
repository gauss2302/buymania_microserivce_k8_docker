package db

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"log"
)

func RunMigrations(dbHost, dbPort, dbUser, dbPassword, dbName string) error {
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	m, err := migrate.New(
		"file://pkg/db/migrations", dbURL)
	if err != nil {
		return fmt.Errorf("error creating migration instance: %w", err)
	}

	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("error running migrations: %w", err)
	}

	log.Println("Successfully ran migrations")
	return nil
}

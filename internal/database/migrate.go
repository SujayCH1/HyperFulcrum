package database

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations() error {

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	migrationPath, err := getMigrationPath()
	if err != nil {
		return err
	}

	m, err := migrate.New(
		"file://"+filepath.ToSlash(migrationPath),
		dsn,
	)

	if err != nil {
		return err
	}

	migrationErr := m.Up()

	sourceErr, databaseErr := m.Close()

	if migrationErr != nil && migrationErr != migrate.ErrNoChange {
		return migrationErr
	}

	if sourceErr != nil {
		return sourceErr
	}

	if databaseErr != nil {
		return databaseErr
	}

	return nil
}

func getMigrationPath() (string, error) {
	migrationPath := os.Getenv("MIGRATIONS_PATH")
	if migrationPath != "" {
		return filepath.Abs(migrationPath)
	}

	executablePath, err := os.Executable()
	if err != nil {
		return "", err
	}

	migrationPath = filepath.Join(
		filepath.Dir(executablePath),
		"migrations",
	)

	if _, err := os.Stat(migrationPath); err == nil {
		return migrationPath, nil
	}

	return filepath.Abs("migrations")
}

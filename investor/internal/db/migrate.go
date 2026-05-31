package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"

	// Register the pgx stdlib driver for database/sql.
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/alekparkhomenko/investor/platform/pkg/logger"
)

// RunMigrations runs all pending database migrations using goose.
// It opens a connection to the database using the pgx driver,
// then applies any pending migration files found in the
// investor/migrations directory (relative to the module root).
func RunMigrations(ctx context.Context, dsn string, log *logger.Logger) error {
	log.Info(ctx, "starting database migrations", logger.Fields{
		"component": "db",
	})

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("opening database connection: %w", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return fmt.Errorf("pinging database: %w", err)
	}

	if err := goose.Up(db, "migrations"); err != nil {
		return fmt.Errorf("running migrations: %w", err)
	}

	log.Info(ctx, "database migrations completed", logger.Fields{
		"component": "db",
	})

	return nil
}

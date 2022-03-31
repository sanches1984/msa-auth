package resources

import (
	"database/sql"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/sanches1984/auth/config"
)

const driverName = "postgres"

func InitDatabase(logger zerolog.Logger) (*sql.DB, error) {
	db, err := sql.Open(driverName, config.SQLDSN())
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	logger.Info().Str("sqldsn", config.SQLDSN()).Msg("db connected")

	if err := migrateDatabase(db, logger); err != nil {
		return nil, err
	}

	return db, nil
}

func migrateDatabase(db *sql.DB, logger zerolog.Logger) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	migration, err := migrate.NewWithDatabaseInstance("file://app/migrations", driverName, driver)
	if err != nil {
		return err
	}

	beforeVersion, dirty, err := migration.Version()
	if err != nil && beforeVersion != 0 {
		return err
	}

	logger.Info().Uint("version", beforeVersion).Msg("migration started")

	if dirty {
		logger.Warn().Msg("previous migration failed")
	}

	err = migration.Up()

	if err != nil && err != migrate.ErrNoChange {
		return err
	} else if err == migrate.ErrNoChange {
		logger.Info().Msg("no new database changes")
	}

	afterVersion, dirty, err := migration.Version()
	if err != nil && beforeVersion != 0 {
		return err
	}

	logger.Info().Uint("version", afterVersion).Msg("migration done")

	if dirty {
		logger.Warn().Msg("previous migration failed")
	}

	return nil
}

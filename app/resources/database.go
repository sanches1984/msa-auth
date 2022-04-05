package resources

import (
	"github.com/go-pg/pg/v9"
	"github.com/rs/zerolog"
	"github.com/sanches1984/auth/config"
	database "github.com/sanches1984/gopkg-pg-orm"
	"github.com/sanches1984/gopkg-pg-orm/migrate"
)

func InitDatabase(logger zerolog.Logger) (database.IClient, error) {
	dsn := config.SQLDSN()
	opts, err := pg.ParseURL(dsn)
	if err != nil {
		return nil, err
	}

	db := database.Connect(config.App, opts)
	if _, err := db.Exec("SELECT 1"); err != nil {
		return nil, err
	}

	logger.Info().Str("dsn", dsn).Msg("db connected")

	if err := migrate.NewMigrator("app/migrations", dsn, migrate.WithLogger(logger)).Run(); err != nil {
		return nil, err
	}

	return db, nil
}

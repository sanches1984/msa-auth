package resources

import (
	"github.com/go-pg/pg/v9"
	"github.com/rs/zerolog"
	"github.com/sanches1984/auth/config"
	database "github.com/sanches1984/gopkg-pg-orm"
	"github.com/sanches1984/gopkg-pg-orm/migrate"
)

func InitDatabase(migrationsPath string, logger zerolog.Logger) (database.IClient, error) {
	dsn := config.Env().SQLDSN
	opts, err := pg.ParseURL(dsn)
	if err != nil {
		return nil, err
	}

	opts.DialTimeout = config.Env().ConnectTimeout
	opts.ReadTimeout = config.Env().ReadTimeout
	opts.WriteTimeout = config.Env().ReadTimeout
	db := database.Connect(config.Env().AppName, opts)
	if _, err := db.Exec("SELECT 1"); err != nil {
		return nil, err
	}

	logger.Info().Str("dsn", dsn).Msg("db connected")

	if err := migrate.NewMigrator(migrationsPath, dsn, migrate.WithLogger(logger)).Run(); err != nil {
		return nil, err
	}

	return db, nil
}

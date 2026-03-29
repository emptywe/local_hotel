package pgdriver

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	migratepgx "github.com/golang-migrate/migrate/v4/database/pgx"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"

	_ "github.com/emptywe/local_hotel/internal/repository/postgres/migrations"
)

type Config struct {
	Username string
	Password string
	Host     string
	Port     string
	DbName   string
}

const (
	maxOpenDbConn = 10
	maxDbLifetime = 5 * time.Minute
)

// NewDBPool creates new postgres connection pool using pgx pgdriver
func NewDBPool(cfg Config, ctx context.Context) (*pgxpool.Pool, error) {
	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DbName)
	db, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, err
	}
	db.Config().MaxConns = maxOpenDbConn
	db.Config().MaxConnLifetime = maxDbLifetime
	if err = db.Ping(ctx); err != nil {
		return nil, err
	}

	if err = migrateDB(cfg, "schema_migration", false); err != nil {
		return nil, err
	}
	return db, nil
}

// migrateDB process db migrations
func migrateDB(cfg Config, table string, down bool) error {
	db, err := sql.Open("pgx", fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s",
		cfg.Host, cfg.Port, cfg.DbName, cfg.Username, cfg.Password))
	if err != nil {
		return err
	}
	defer db.Close()
	driver, err := migratepgx.WithInstance(db, &migratepgx.Config{MigrationsTable: table})
	if err != nil {
		return err
	}
	migrator, err := migrate.NewWithDatabaseInstance("embed://", table, driver)
	if err != nil {
		return err
	}
	if down {
		err = migrator.Down()
		if err != nil && err.Error() == "no change" { // "no change" is not an error
			err = nil
		}
	} else {
		err = migrator.Up()
		if err != nil && err.Error() == "no change" { // "no change" is not an error
			err = nil
		}
	}

	return err
}

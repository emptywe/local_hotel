package dbrepo

import (
	"github.com/emptywe/local_hotel/config"
	"github.com/emptywe/local_hotel/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

// postgresDBRepo postgres repo struct
type postgresDBRepo struct {
	App *config.AppConfig
	DB  *pgxpool.Pool
}

// NewPostgresRepo creating new postgres repo
func NewPostgresRepo(conn *pgxpool.Pool, app *config.AppConfig) repository.DataBaseRepo {
	return &postgresDBRepo{
		App: app,
		DB:  conn,
	}
}

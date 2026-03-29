package handler

import (
	"github.com/emptywe/local_hotel/config"
	"github.com/emptywe/local_hotel/internal/render"
	"github.com/emptywe/local_hotel/internal/repository"
	"github.com/emptywe/local_hotel/internal/repository/postgres/dbrepo"
	"github.com/jackc/pgx/v5/pgxpool"
)

var Handle *Handler // Handle - the repository for handlers

// Handler is repository type for handlers
type Handler struct {
	App      *config.AppConfig
	renderer *render.Renderer
	DB       repository.DataBaseRepo
}

// NewHandler create and set new handler repository
func NewHandler(a *config.AppConfig, renderer *render.Renderer, db *pgxpool.Pool) {
	Handle = &Handler{
		App:      a,
		renderer: renderer,
		DB:       dbrepo.NewPostgresRepo(db, a),
	}
}

// NewTestHandler create and set new test handler repository
func NewTestHandler(a *config.AppConfig, renderer *render.Renderer) {
	Handle = &Handler{
		App:      a,
		renderer: renderer,
		DB:       dbrepo.NewTestRepo(a),
	}
}

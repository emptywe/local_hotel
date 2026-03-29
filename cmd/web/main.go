package main

import (
	"context"
	"encoding/gob"
	"fmt"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/emptywe/local_hotel/config"
	"github.com/emptywe/local_hotel/internal/handler"
	"github.com/emptywe/local_hotel/internal/render"
	"github.com/emptywe/local_hotel/internal/repository/postgres/pgdriver"
	"github.com/emptywe/local_hotel/internal/router"
	"github.com/emptywe/local_hotel/internal/server"
	"github.com/emptywe/local_hotel/model"
	"github.com/emptywe/local_hotel/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

const (
	defAddr         = "127.0.0.1:8080" // Address for local development
	dockerAddr      = "0.0.0.0:8080"   // Address for docker containers
	sessionLifetime = time.Hour * 24
)

var pageStatus = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_request_get_page_status_count", // metric name
		Help: "Count of status returned by page.",
	},
	[]string{"page", "status"}, // labels
)

func init() {
	// we need to register the counter so prometheus can collect this metric
	prometheus.MustRegister(pageStatus)
}

func main() {
	pass, _ := bcrypt.GenerateFromPassword([]byte("password"), 10)
	fmt.Println(string(pass))
	var app config.AppConfig

	defer close(app.MailChan)

	if err := setupLogger(); err != nil {
		panic(err)
	}

	db, err := connectToPostgres()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	zap.S().Debug("Connected to postgres database")

	handlers, err := setup(&app, db)
	if err != nil {
		panic(err)
	}
	zap.S().Debug("Logger initialised")

	go ListenForMail(&app)

	srv := &server.Server{
		Addr:    defAddr,
		Handler: handlers,
	}
	pageStatus.WithLabelValues("hello", "world").Inc()
	app.Stats = pageStatus
	zap.S().Debug("Server start")
	if err = srv.Run(); err != nil && err != http.ErrServerClosed {
		zap.S().Fatalf("problem runing server: %v", err)
	}
}

func setup(app *config.AppConfig, db *pgxpool.Pool) (http.Handler, error) {

	app.MailChan = make(chan model.MailData, 1000)
	// telling application that we will store more complex types
	gob.Register(model.Reservation{})
	gob.Register(model.User{})
	gob.Register(model.Room{})
	gob.Register(model.Restriction{})
	gob.Register(map[string]int{})
	// for production change to true
	app.InProduction = false

	renderer, err := setupRenderer(app)
	if err != nil {
		return nil, err
	}
	app.UseCache = false

	setupSession(app)

	handler.NewHandler(app, renderer, db)
	handlers := router.InitRoutes(app)

	return handlers, nil
}

func setupLogger() error {
	logCfg := logger.NewConfig(logger.Settings{
		DisableCaller:     false,
		DisableStacktrace: true,
		Colour:            false,
		Level:             0,
	})
	if err := logCfg.InitLogger(); err != nil {
		return err
	}
	return nil
}

func connectToPostgres() (*pgxpool.Pool, error) {
	db, err := pgdriver.NewDBPool(pgdriver.Config{
		Username: "postgres",
		Password: "1234",
		Host:     "localhost",
		Port:     "5432",
		DbName:   "bookings",
	}, context.Background())
	return db, err
}

func setupRenderer(app *config.AppConfig) (*render.Renderer, error) {
	renderer := render.NewRenderer(app)
	tc, err := renderer.CreateTemplateCache()
	if err != nil {
		return nil, err
	}
	app.TemplateCache = tc
	return renderer, nil
}

func setupSession(app *config.AppConfig) {
	session := scs.New()
	session.Lifetime = sessionLifetime
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction
	app.Session = session
}

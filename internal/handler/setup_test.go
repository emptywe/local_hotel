package handler

import (
	"encoding/gob"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/emptywe/local_hotel/config"
	"github.com/emptywe/local_hotel/internal/render"
	"github.com/emptywe/local_hotel/model"
	"github.com/emptywe/local_hotel/pkg/logger"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type testHandler struct{}

func (mh *testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}

var (
	mockApp         config.AppConfig
	pathToTemplates = "./../../html-templates"
)

func TestMain(m *testing.M) {
	// telling application that we will store more complex types
	gob.Register(model.Reservation{})
	gob.Register(model.User{})
	gob.Register(model.Room{})
	gob.Register(model.Restriction{})
	// for production change to true
	mockApp.InProduction = false
	err := setupLogger()
	if err != nil {
		os.Exit(2)
	}

	renderer, err := setupRenderer(&mockApp)
	if err != nil {
		os.Exit(2)
	}

	mockApp.MailChan = make(chan model.MailData, 10)
	defer close(mockApp.MailChan)
	listenForMail(&mockApp)

	mockApp.UseCache = true
	setupSession(&mockApp)
	NewTestHandler(&mockApp, renderer)

	os.Exit(m.Run())
}

func setupLogger() error {
	logCfg := logger.NewConfig(logger.Settings{
		DisableCaller:     false,
		DisableStacktrace: true,
		Colour:            false,
		Level:             4,
	})
	if err := logCfg.InitLogger(); err != nil {
		return err
	}
	return nil
}

func setupRenderer(app *config.AppConfig) (*render.Renderer, error) {
	renderer := render.NewRenderer(app)
	tc, err := CreateTestTemplateCache()
	if err != nil {
		return nil, err
	}
	app.TemplateCache = tc
	return renderer, nil
}

func setupSession(app *config.AppConfig) {
	session := scs.New()
	session.Lifetime = time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction
	app.Session = session
}

// CreateTemplateCache creates a template cache as a map
func CreateTestTemplateCache() (map[string]*template.Template, error) {
	tplCache := make(map[string]*template.Template)
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.gohtml", pathToTemplates))
	if err != nil {
		return tplCache, err
	}

	for _, page := range pages {

		name := filepath.Base(page)
		zap.S().Debugf("Page is: %s", page)
		ts, err := template.New(name).Funcs(nil).ParseFiles(page)
		if err != nil {
			return tplCache, err
		}

		match, err := filepath.Glob(fmt.Sprintf("%s/*.layout.gohtml", pathToTemplates))
		if err != nil {
			return tplCache, err
		}

		if len(match) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.gohtml", pathToTemplates))
			if err != nil {
				return tplCache, err
			}
		}

		tplCache[name] = ts
	}

	return tplCache, nil
}

func listenForMail(app *config.AppConfig) {
	go func() {
		for _ = range app.MailChan {

		}
	}()
}

func getRoutes(app *config.AppConfig) http.Handler {
	router := mux.NewRouter()

	router.Use(handlers.RecoveryHandler())
	//router.Use(Handle.NoSurf)
	router.Use(Handle.SessionLoad)

	router.HandleFunc("/", Handle.Home).Methods("GET")
	router.HandleFunc("/about", Handle.About).Methods("GET")
	router.HandleFunc("/alter-doctor", Handle.Doctor).Methods("GET")
	router.HandleFunc("/hayat", Handle.Hayat).Methods("GET")

	router.HandleFunc("/search-availability", Handle.Availability).Methods("GET")
	router.HandleFunc("/search-availability", Handle.CheckAvailability).Methods("POST")
	router.HandleFunc("/search-availability-json", Handle.CheckAvailabilityJSON).Methods("POST")
	router.HandleFunc("/choose-room/{id}", Handle.ChooseRoom).Methods("GET")
	router.HandleFunc("/book-room", Handle.BookRoom).Methods("GET")

	router.HandleFunc("/make-reservation", Handle.Reservation).Methods("GET")
	router.HandleFunc("/make-reservation", Handle.PostReservation).Methods("POST")
	router.HandleFunc("/reservation-summary", Handle.ReservationSummary).Methods("GET")

	router.HandleFunc("/contact", Handle.Contact).Methods("GET")

	fileServer := http.FileServer(http.Dir("./static/"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fileServer))

	return router
}

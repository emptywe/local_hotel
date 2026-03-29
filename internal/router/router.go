package router

import (
	"net/http"

	"github.com/emptywe/local_hotel/config"
	"github.com/emptywe/local_hotel/internal/handler"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// InitRoutes initiate routes with gorilla mux and return Handler
func InitRoutes(app *config.AppConfig) http.Handler {
	router := mux.NewRouter()

	router.Use(handlers.RecoveryHandler())
	router.Use(handler.Handle.NoSurf)
	router.Use(handler.Handle.SessionLoad)

	router.Handle("/metrics", promhttp.Handler())

	router.HandleFunc("/", handler.Handle.Home).Methods("GET")
	router.HandleFunc("/about", handler.Handle.About).Methods("GET")
	router.HandleFunc("/alter-doctor", handler.Handle.Doctor).Methods("GET")
	router.HandleFunc("/hayat", handler.Handle.Hayat).Methods("GET")

	router.HandleFunc("/search-availability", handler.Handle.Availability).Methods("GET")
	router.HandleFunc("/search-availability", handler.Handle.CheckAvailability).Methods("POST")
	router.HandleFunc("/search-availability-json", handler.Handle.CheckAvailabilityJSON).Methods("POST")
	router.HandleFunc("/choose-room/{id}", handler.Handle.ChooseRoom).Methods("GET")
	router.HandleFunc("/book-room", handler.Handle.BookRoom).Methods("GET")

	router.HandleFunc("/make-reservation", handler.Handle.Reservation).Methods("GET")
	router.HandleFunc("/make-reservation", handler.Handle.PostReservation).Methods("POST")
	router.HandleFunc("/reservation-summary", handler.Handle.ReservationSummary).Methods("GET")

	router.HandleFunc("/contact", handler.Handle.Contact).Methods("GET")

	router.HandleFunc("/user/login", handler.Handle.ShowLogin).Methods("GET")
	router.HandleFunc("/user/login", handler.Handle.PostShowLogin).Methods("POST")
	router.HandleFunc("/user/logout", handler.Handle.Logout).Methods("GET")

	adminRouter := router.PathPrefix("/admin").Methods("GET", "POST").Subrouter()
	adminRouter.Use(handler.Handle.Auth)
	adminRouter.HandleFunc("/dashboard", handler.Handle.AdminDashboard).Methods("GET")

	adminRouter.HandleFunc("/reservations-new", handler.Handle.AdminNewReservations).Methods("GET")
	adminRouter.HandleFunc("/reservations-all", handler.Handle.AdminAllReservations).Methods("GET")

	adminRouter.HandleFunc("/reservations-calendar", handler.Handle.AdminCalendarReservations).Methods("GET")
	adminRouter.HandleFunc("/reservations-calendar", handler.Handle.AdminCalendarPostReservations).Methods("POST")

	adminRouter.HandleFunc("/reservations/{src}/{id}/show", handler.Handle.AdminShowReservation).Methods("GET")
	adminRouter.HandleFunc("/reservations/{src}/{id}", handler.Handle.AdminPostShowReservation).Methods("POST")

	adminRouter.HandleFunc("/process-reservation/{src}/{id}/do", handler.Handle.AdminProcessReservation).Methods("GET")
	adminRouter.HandleFunc("/delete-reservation/{src}/{id}/delete", handler.Handle.AdminDeleteReservation).Methods("GET")

	fileServer := http.FileServer(http.Dir("./static/"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fileServer))

	return router
}

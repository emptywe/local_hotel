package handler

import (
	"net/http"

	"github.com/emptywe/local_hotel/internal/helpers"
	"github.com/justinas/nosurf"
)

// NoSurf adds CSRF protection to all POST requests
func (h *Handler) NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   h.App.InProduction,
		SameSite: http.SameSiteLaxMode,
	})

	return csrfHandler
}

// SessionLoad loads and saves session on every request
func (h *Handler) SessionLoad(next http.Handler) http.Handler {
	return h.App.Session.LoadAndSave(next)
}

// Auth - checks if user are logged in
func (h *Handler) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if !helpers.IsAuthenticated(r, h.App.Session) {
			h.App.Session.Put(r.Context(), "error", "Log in first")
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}

package handler

import (
	"errors"
	"net/http"

	"github.com/emptywe/local_hotel/internal/forms"
	"github.com/emptywe/local_hotel/model"
	"go.uber.org/zap"
)

// Home is handler for home page
func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	h.renderer.Template(w, r, "home.page.gohtml", &model.TemplateData{})
	h.App.Stats.WithLabelValues("Home", "200").Inc()
}

// About is handler for about page
func (h *Handler) About(w http.ResponseWriter, r *http.Request) {
	h.renderer.Template(w, r, "about.page.gohtml", &model.TemplateData{})
	h.App.Stats.WithLabelValues("About", "200").Inc()
}

// Reservation renders the make reservation page and displays form
func (h *Handler) Reservation(w http.ResponseWriter, r *http.Request) {
	res, ok := h.App.Session.Get(r.Context(), "reservation").(model.Reservation)
	if !ok {
		err := errors.New("cannot get reservation from session")
		zap.S().Error(err)
		h.App.Session.Put(r.Context(), "error", err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	room, err := h.DB.GetRoomById(res.RoomId)
	if err != nil {
		zap.S().Errorf("can't get room by id: %v", err)
		h.App.Session.Put(r.Context(), "error", "can't find room")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	res.Room.RoomName = room.RoomName

	h.App.Session.Put(r.Context(), "reservation", res)

	startDateString := res.StartDate.Format("2006-01-02")
	endDateString := res.EndDate.Format("2006-01-02")

	stringMap := make(map[string]string)
	stringMap["start_date"] = startDateString
	stringMap["end_date"] = endDateString

	data := make(map[string]interface{})
	data["reservation"] = res
	h.renderer.Template(w, r, "make-reservation.page.gohtml", &model.TemplateData{
		Form:      forms.New(nil),
		Data:      data,
		StringMap: stringMap,
	})
	h.App.Stats.WithLabelValues("Reservation", "200").Inc()
}

// Doctor renders the make alter doctor page and displays form
func (h *Handler) Doctor(w http.ResponseWriter, r *http.Request) {
	h.renderer.Template(w, r, "doctor.page.gohtml", &model.TemplateData{})
	h.App.Stats.WithLabelValues("Doctor", "200").Inc()
}

// Hayat renders the make hayat page and displays form
func (h *Handler) Hayat(w http.ResponseWriter, r *http.Request) {
	h.renderer.Template(w, r, "hayat.page.gohtml", &model.TemplateData{})
	h.App.Stats.WithLabelValues("Hayat", "200").Inc()
}

// Availability renders the make search availability page and displays form
func (h *Handler) Availability(w http.ResponseWriter, r *http.Request) {
	h.renderer.Template(w, r, "search-availability.page.gohtml", &model.TemplateData{})
	h.App.Stats.WithLabelValues("Availability", "200").Inc()
}

// Contact renders the make contact page and displays form
func (h *Handler) Contact(w http.ResponseWriter, r *http.Request) {
	h.renderer.Template(w, r, "contact.page.gohtml", &model.TemplateData{})
	h.App.Stats.WithLabelValues("Contact", "200").Inc()
}

// ReservationSummary renders the reservation summary page and displays form
func (h *Handler) ReservationSummary(w http.ResponseWriter, r *http.Request) {

	reservation, ok := h.App.Session.Get(r.Context(), "reservation").(model.Reservation)
	if !ok {
		zap.S().Error("can't get reservation from request")
		h.App.Session.Put(r.Context(), "error", "can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	h.App.Session.Remove(r.Context(), "reservation")

	data := make(map[string]interface{})
	data["reservation"] = reservation

	startDateString := reservation.StartDate.Format("2006-01-02")
	endDateString := reservation.EndDate.Format("2006-01-02")
	stringMap := make(map[string]string)
	stringMap["start_date"] = startDateString
	stringMap["end_date"] = endDateString

	h.renderer.Template(w, r, "reservation-summary.page.gohtml", &model.TemplateData{
		Data:      data,
		StringMap: stringMap,
	})
	h.App.Stats.WithLabelValues("Reservation summary", "200").Inc()
}

// ShowLogin - shows the login page
func (h *Handler) ShowLogin(w http.ResponseWriter, r *http.Request) {
	h.renderer.Template(w, r, "login.page.gohtml", &model.TemplateData{
		Form: forms.New(nil),
	})

}

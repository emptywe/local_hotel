package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/emptywe/local_hotel/internal/forms"
	"github.com/emptywe/local_hotel/internal/helpers"
	"github.com/emptywe/local_hotel/model"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// PostShowLogin - handles logging the user in
func (h *Handler) PostShowLogin(w http.ResponseWriter, r *http.Request) {
	_ = h.App.Session.RenewToken(r.Context())

	if err := r.ParseForm(); err != nil {
		zap.S().Error("can't empty login form")
		h.App.Session.Put(r.Context(), "error", "empty login form")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("email", "password")
	if !form.Valid() {
		_ = h.renderer.Template(w, r, "login.page.template", &model.TemplateData{Form: form})
		return
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	id, err := h.DB.Authenticate(email, password)
	if err != nil {
		zap.S().Errorf("can't authenticate user: %s", err.Error())
		h.App.Session.Put(r.Context(), "error", "invalid login credentials")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	h.App.Session.Put(r.Context(), "user_id", id)
	h.App.Session.Put(r.Context(), "flash", "logged in successfully")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Logout - logs the user out
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	_ = h.App.Session.Destroy(r.Context())
	_ = h.App.Session.RenewToken(r.Context())
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

// AdminDashboard - render admin dashboard page
func (h *Handler) AdminDashboard(w http.ResponseWriter, r *http.Request) {
	_ = h.renderer.Template(w, r, "admin-dashboard.page.gohtml", &model.TemplateData{})
}

// AdminNewReservations - shows all new unprocessed reservations
func (h *Handler) AdminNewReservations(w http.ResponseWriter, r *http.Request) {

	reservations, err := h.DB.NewReservations()
	if err != nil {
		zap.S().Errorf("can't get all reservations: %s", err.Error())
		helpers.JsonErrorResponse(w, "can't get reservations: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data := make(map[string]interface{})
	data["reservations"] = reservations

	_ = h.renderer.Template(w, r, "admin-new-reservations.page.gohtml", &model.TemplateData{
		Data: data,
	})
}

// AdminAllReservations - shows all reservations
func (h *Handler) AdminAllReservations(w http.ResponseWriter, r *http.Request) {
	reservations, err := h.DB.AllReservations()
	if err != nil {
		zap.S().Errorf("can't get all reservations: %s", err.Error())
		helpers.JsonErrorResponse(w, "can't get reservations: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data := make(map[string]interface{})
	data["reservations"] = reservations

	_ = h.renderer.Template(w, r, "admin-all-reservations.page.gohtml", &model.TemplateData{
		Data: data,
	})
}

// AdminShowReservation - shows the reservation in the admin tool
func (h *Handler) AdminShowReservation(w http.ResponseWriter, r *http.Request) {
	params := strings.Split(r.RequestURI, "/")
	id, err := strconv.Atoi(params[len(params)-2])
	if err != nil {
		zap.S().Errorf("invalid reservation id: %s", err.Error())
		helpers.JsonErrorResponse(w, "invalid reservation id", http.StatusBadRequest)
		return
	}
	src := params[len(params)-3]
	stringMap := make(map[string]string)
	stringMap["src"] = src

	year := r.URL.Query().Get("y")
	month := r.URL.Query().Get("m")

	stringMap["month"] = month
	stringMap["year"] = year

	res, err := h.DB.GetReservationById(id)
	if err != nil {
		zap.S().Errorf("can't get reservation id: %s", err.Error())
		helpers.JsonErrorResponse(w, "can't get reservation id", http.StatusInternalServerError)
		return
	}

	data := make(map[string]interface{})
	data["reservation"] = res

	_ = h.renderer.Template(w, r, "admin-reservations-show.page.gohtml", &model.TemplateData{
		StringMap: stringMap,
		Data:      data,
		Form:      forms.New(nil),
	})
}

// AdminPostShowReservation - inserts the reservation changes from admin tool
func (h *Handler) AdminPostShowReservation(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		zap.S().Errorf("invalid reservation form: %s", err.Error())
		helpers.JsonErrorResponse(w, "invalid reservation from", http.StatusBadRequest)
		return
	}

	params := strings.Split(r.RequestURI, "/")
	id, err := strconv.Atoi(params[4])
	if err != nil {
		zap.S().Errorf("invalid reservation id: %s", err.Error())
		helpers.JsonErrorResponse(w, "invalid reservation id", http.StatusBadRequest)
		return
	}
	src := params[3]
	stringMap := make(map[string]string)
	stringMap["src"] = src

	res, err := h.DB.GetReservationById(id)

	res.FirstName = r.Form.Get("first_name")
	res.LastName = r.Form.Get("last_name")
	res.Email = r.Form.Get("email")
	res.Phone = r.Form.Get("phone")

	if err = h.DB.UpdateReservation(res); err != nil {
		zap.S().Errorf("can't update reservation: %s", err.Error())
		helpers.JsonErrorResponse(w, "can't update reservation", http.StatusInternalServerError)
		return
	}

	month := r.Form.Get("month")
	year := r.Form.Get("year")

	h.App.Session.Put(r.Context(), "flash", "Changes saved")
	if year == "" {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)
	} else {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-calendar?y=%s&m-%s", year, month), http.StatusSeeOther)
	}

}

// AdminCalendarReservations - shows the calendar with reservations
func (h *Handler) AdminCalendarReservations(w http.ResponseWriter, r *http.Request) {
	now := time.Now()

	if r.URL.Query().Get("y") != "" {
		year, _ := strconv.Atoi(r.URL.Query().Get("y"))
		month, _ := strconv.Atoi(r.URL.Query().Get("m"))
		now = time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	}

	data := make(map[string]interface{})
	data["now"] = now

	next := now.AddDate(0, 1, 0)
	last := now.AddDate(0, -1, 0)

	nextMonth := next.Format("01")
	nextMonthYear := next.Format("2006")
	lastMonth := last.Format("01")
	lastMonthYear := last.Format("2006")

	stringMap := make(map[string]string)
	stringMap["next_month"] = nextMonth
	stringMap["next_month_year"] = nextMonthYear
	stringMap["last_month"] = lastMonth
	stringMap["last_month_year"] = lastMonthYear
	stringMap["this_month"] = now.Format("01")
	stringMap["this_month_year"] = now.Format("2006")

	currYear, currMonth, _ := now.Date()
	currLoc := now.Location()
	firstOfMonth := time.Date(currYear, currMonth, 1, 0, 0, 0, 0, currLoc)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

	intMap := make(map[string]int)
	intMap["days_in_month"] = lastOfMonth.Day()

	rooms, err := h.DB.AllRooms()
	if err != nil {
		zap.S().Errorf("can't get all rooms: %s", err.Error())
		helpers.JsonErrorResponse(w, "can't get all rooms", http.StatusInternalServerError)
		return
	}
	data["rooms"] = rooms

	for _, room := range rooms {
		reservationMap := make(map[string]int)
		blockMap := make(map[string]int)
		for d := firstOfMonth; !d.After(lastOfMonth); d = d.AddDate(0, 0, 1) {
			reservationMap[d.Format("2006-01-2")] = 0
			blockMap[d.Format("2006-01-2")] = 0
		}

		restrictions, err := h.DB.GetRestrictionsForRoomByDate(room.Id, firstOfMonth, lastOfMonth)
		if err != nil {
			zap.S().Errorf("can't get room restrictions: %s", err.Error())
			helpers.JsonErrorResponse(w, "can't get room restrictions", http.StatusInternalServerError)
			return
		}

		for _, restriction := range restrictions {
			if restriction.ReservationId > 0 {
				for d := restriction.StartDate; !d.After(restriction.EndDate); d = d.AddDate(0, 0, 1) {
					reservationMap[d.Format("2006-01-2")] = restriction.ReservationId
				}
			} else {
				blockMap[restriction.StartDate.Format("2006-01-2")] = restriction.Id
			}
		}

		data[fmt.Sprintf("reservation_map_%d", room.Id)] = reservationMap
		data[fmt.Sprintf("block_map_%d", room.Id)] = blockMap

		h.App.Session.Put(r.Context(), fmt.Sprintf("block_map_%d", room.Id), blockMap)
	}

	_ = h.renderer.Template(w, r, "admin-calendar-reservations.page.gohtml", &model.TemplateData{
		StringMap: stringMap,
		Data:      data,
		IntMap:    intMap,
	})
}

// AdminProcessReservation - marks as processed reservation by id
func (h *Handler) AdminProcessReservation(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		zap.S().Errorf("invalid reservation id: %s", err.Error())
		helpers.JsonErrorResponse(w, "invalid reservation id", http.StatusBadRequest)
		return
	}
	src := mux.Vars(r)["src"]

	if err = h.DB.UpdateProcessedReservation(id, model.Processed); err != nil {
		zap.S().Errorf("Can't mark reservation as processed: %s", err.Error())
		helpers.JsonErrorResponse(w, "Can't mark reservation as processed", http.StatusInternalServerError)
		return
	}

	year := r.URL.Query().Get("y")
	month := r.URL.Query().Get("m")

	h.App.Session.Put(r.Context(), "flash", "Reservation marked as processed")
	if year == "" {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)
	} else {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-calendar?y=%s&m=%s", year, month), http.StatusSeeOther)
	}
}

// AdminDeleteReservation - deletes reservation by id
func (h *Handler) AdminDeleteReservation(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		zap.S().Errorf("invalid reservation id: %s", err.Error())
		helpers.JsonErrorResponse(w, "invalid reservation id", http.StatusBadRequest)
		return
	}
	src := mux.Vars(r)["src"]

	if err = h.DB.DeleteReservation(id); err != nil {
		zap.S().Errorf("Can't delete reservation: %s", err.Error())
		helpers.JsonErrorResponse(w, "Can't delete reservation ", http.StatusInternalServerError)
		return
	}
	h.App.Session.Put(r.Context(), "flash", "Reservation deleted")

	year := r.URL.Query().Get("y")
	month := r.URL.Query().Get("m")

	http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)
	if year == "" {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)
	} else {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-calendar?y=%s&m=%s", year, month), http.StatusSeeOther)
	}
}

// AdminCalendarPostReservations - handles post of reservation calendar
func (h *Handler) AdminCalendarPostReservations(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		zap.S().Errorf("can't parse form for AdminCalendarPostReservations: %s", err.Error())
		helpers.JsonErrorResponse(w, "can't parse form", http.StatusBadRequest)
		return
	}

	year, _ := strconv.Atoi(r.Form.Get("y"))
	month, _ := strconv.Atoi(r.Form.Get("m"))

	rooms, err := h.DB.AllRooms()
	if err != nil {
		zap.S().Errorf("can't get all rooms: %s", err.Error())
		helpers.JsonErrorResponse(w, "can't get all rooms", http.StatusInternalServerError)
		return
	}

	form := forms.New(r.PostForm)

	for _, room := range rooms {
		curMap := h.App.Session.Get(r.Context(), fmt.Sprintf("block_map_%d", room.Id)).(map[string]int)
		for name, value := range curMap {

			if value > 0 && !form.Has(fmt.Sprintf("remove_block_%d_%s", room.Id, name)) {
				if err = h.DB.DeleteBlockById(value); err != nil {
					zap.S().Errorf("can't delete block: %s", err.Error())
				}
			}
		}
	}

	for name, _ := range r.PostForm {
		if strings.HasPrefix(name, "add_block") {
			exploded := strings.Split(name, "_")
			roomId, _ := strconv.Atoi(exploded[2])
			t, _ := time.Parse("2006-01-2", exploded[3])
			if err = h.DB.InsertBlockForRoom(roomId, t); err != nil {
				zap.S().Errorf("can't insert block: %s", err.Error())
			}
		}
	}

	h.App.Session.Put(r.Context(), "flash", "Changes Saved")
	http.Redirect(w, r, fmt.Sprintf("/admin/reservations-calendar?y=%d&m=%d", year, month), http.StatusSeeOther)
}

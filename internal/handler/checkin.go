package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/emptywe/local_hotel/internal/forms"
	"github.com/emptywe/local_hotel/internal/helpers"
	"github.com/emptywe/local_hotel/model"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// CheckAvailability renders the make search availability page and displays form
func (h *Handler) CheckAvailability(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		zap.S().Errorf("can't parse form %v", err)
		helpers.ServerError(w, err)
		return
	}

	startDateString := r.Form.Get("start")
	endDateString := r.Form.Get("end")
	// TODO: add start != end
	startDate, err := time.Parse(model.DateLayout, startDateString)
	if err != nil {
		zap.S().Errorf("can't parse start date: %v", err)
		helpers.ServerError(w, err)
		return
	}
	endDate, err := time.Parse(model.DateLayout, endDateString)
	if err != nil {
		zap.S().Errorf("can't parse end date: %v", err)
		helpers.ServerError(w, err)
		return
	}

	rooms, err := h.DB.SearchAvailabilityForAllRooms(startDate, endDate)
	if err != nil && err != sql.ErrNoRows {
		zap.S().Errorf("problem search for availability %v", err)
		helpers.ServerError(w, err)
		return
	}

	if len(rooms) == 0 {
		h.App.Session.Put(r.Context(), "error", "This dates are not available")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})
	data["rooms"] = rooms

	res := model.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}

	h.App.Session.Put(r.Context(), "reservation", res)

	h.renderer.Template(w, r, "choose_room.page.gohtml", &model.TemplateData{
		Data: data,
	})

}

// CheckAvailabilityJSON handles request for availability and send JSON response
func (h *Handler) CheckAvailabilityJSON(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		zap.S().Errorf("can't parse form: %v", err)
		helpers.JsonErrorResponse(w, "can't parse form", http.StatusInternalServerError)
		return
	}

	startDateString := r.Form.Get("start")
	endDateString := r.Form.Get("end")
	roomIdString := r.Form.Get("room_id")

	startDate, err := time.Parse(model.DateLayout, startDateString)
	if err != nil {
		zap.S().Errorf("can't parse start date: %v", err)
		helpers.JsonErrorResponse(w, "can't parse start date", http.StatusInternalServerError)
		return
	}

	endDate, err := time.Parse(model.DateLayout, endDateString)
	if err != nil {
		zap.S().Errorf("can't parse end date: %v", err)
		helpers.JsonErrorResponse(w, "can't parse end date", http.StatusInternalServerError)
		return
	}

	roomId, err := strconv.Atoi(roomIdString)
	if err != nil {
		zap.S().Errorf("can't parse room id: %v", err)
		helpers.JsonErrorResponse(w, "can't parse room id", http.StatusInternalServerError)
		return
	}

	available, err := h.DB.SearchAvailabilityByDatesByRoomId(startDate, endDate, roomId)
	if err != nil {
		zap.S().Errorf("error connecting to database: %v", err)
		helpers.JsonErrorResponse(w, "error connecting to database", http.StatusInternalServerError)
		return
	}

	resp := model.JsonResponse{
		Ok:        available,
		Message:   "",
		StartDate: startDateString,
		EndDate:   endDateString,
		RoomId:    roomIdString,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)

}

// PostReservation handles the posting of a reservation form
func (h *Handler) PostReservation(w http.ResponseWriter, r *http.Request) {

	reservation, ok := h.App.Session.Get(r.Context(), "reservation").(model.Reservation)
	if !ok {
		err := errors.New("can't get reservation from session")
		h.App.Session.Put(r.Context(), "error", err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	err := r.ParseForm()
	if err != nil {
		zap.S().Errorf("can't parse request form: %v", err)
		h.App.Session.Put(r.Context(), "error", "can't parse request form")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	reservation.FirstName = r.Form.Get("first_name")
	reservation.LastName = r.Form.Get("last_name")
	reservation.Email = r.Form.Get("email")
	reservation.Phone = r.Form.Get("phone")

	form := forms.New(r.PostForm)

	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		//http.Error(w, "invalid from", http.StatusTemporaryRedirect)
		h.renderer.Template(w, r, "make-reservation.page.gohtml", &model.TemplateData{
			Form: form,
			Data: data,
		})
		zap.S().Error("invalid form")
		return
	}

	reservationId, err := h.DB.InsertReservation(reservation)
	if err != nil {
		zap.S().Errorf("can't insert reservation: %v", err)
		h.App.Session.Put(r.Context(), "error", "can't insert reservation")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	restriction := model.RoomRestriction{
		StartDate:     reservation.StartDate,
		EndDate:       reservation.EndDate,
		RoomId:        reservation.RoomId,
		ReservationId: reservationId,
		RestrictionId: 1,
	}

	if err = h.DB.InsertRoomRestriction(restriction); err != nil {
		zap.S().Errorf("can't insert room restriction: %v", err)
		h.App.Session.Put(r.Context(), "error", "can't insert room restriction")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	htmlMessage := fmt.Sprintf(`
		<strong>Reservation confirmation</strong>
		Dear %s:,<br>
		This is confirm your reservation from %s to %s.
	`, reservation.FirstName, reservation.StartDate.Format("2006-01-02"), reservation.EndDate.Format("2006-01-02"))

	msg := model.MailData{
		To:       reservation.Email,
		From:     "out@gmail.com",
		Subject:  "Reservation confirmation",
		Content:  htmlMessage,
		Template: "basic.gohtml",
	}
	h.App.MailChan <- msg

	htmlMessage = fmt.Sprintf(`
		<strong>Reservation confirmation</strong>
		For %s:,<br>
		Reservation has been made from %s to %s.
	`, reservation.Room.RoomName, reservation.StartDate.Format("2006-01-02"), reservation.EndDate.Format("2006-01-02"))

	msg = model.MailData{
		To:       "owner@gmail.com",
		From:     "out@gmail.com",
		Subject:  "Reservation confirmation",
		Content:  htmlMessage,
		Template: "basic.gohtml",
	}
	h.App.MailChan <- msg

	h.App.Session.Put(r.Context(), "reservation", reservation)
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

// ChooseRoom inserts a chosen room id into reservation in session
func (h *Handler) ChooseRoom(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	roomId, err := strconv.Atoi(vars["id"])
	if err != nil {
		zap.S().Errorf("invalid room id: %v", err)
		helpers.ServerError(w, err)
		return
	}

	res, ok := h.App.Session.Get(r.Context(), "reservation").(model.Reservation)
	if !ok {
		err = errors.New("cannot get reservation from session")
		zap.S().Error(err)
		helpers.ServerError(w, err)
		return
	}

	res.RoomId = roomId

	h.App.Session.Put(r.Context(), "reservation", res)

	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}

// BookRoom creates a reservation by URL parameters and redirects to make reservation page
func (h *Handler) BookRoom(w http.ResponseWriter, r *http.Request) {

	var res model.Reservation

	roomIdString := r.URL.Query().Get("id")
	startDateString := r.URL.Query().Get("s")
	endDateString := r.URL.Query().Get("e")

	roomId, err := strconv.Atoi(roomIdString)
	if err != nil {
		zap.S().Errorf("can't parse id: %v", err)
		helpers.ServerError(w, err)
		return
	}

	startDate, err := time.Parse(model.DateLayout, startDateString)
	if err != nil {
		zap.S().Errorf("can't parse start date: %v", err)
		helpers.ServerError(w, err)
		return
	}
	endDate, err := time.Parse(model.DateLayout, endDateString)
	if err != nil {
		zap.S().Errorf("can't parse end date: %v", err)
		helpers.ServerError(w, err)
		return
	}

	room, err := h.DB.GetRoomById(roomId)
	if err != nil {
		zap.S().Errorf("can't get room by id: %v", err)
		helpers.ServerError(w, err)
		return
	}

	res.RoomId = roomId
	res.StartDate = startDate
	res.EndDate = endDate
	res.Room.RoomName = room.RoomName

	h.App.Session.Put(r.Context(), "reservation", res)

	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}

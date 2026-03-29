package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/emptywe/local_hotel/model"
	"github.com/gorilla/mux"
)

type postData struct {
	key   string
	value string
}

var TestsGet = []struct {
	name               string
	url                string
	expectedStatusCode int
}{
	{"home", "/", http.StatusOK},
	{"about", "/about", http.StatusOK},
	{"alter-doctor", "/alter-doctor", http.StatusOK},
	{"hayat", "/hayat", http.StatusOK},
	{"search availability", "/search-availability", http.StatusOK},
	{"contact", "/contact", http.StatusOK},
}

func TestHandlers(t *testing.T) {
	router := getRoutes(&mockApp)
	ts := httptest.NewTLSServer(router)
	defer ts.Close()

	for _, test := range TestsGet {
		resp, err := ts.Client().Get(ts.URL + test.url)
		if err != nil {
			t.Fatal(fmt.Sprintf("test %s failed: %v", test.name, err))
		}
		if resp.StatusCode != test.expectedStatusCode {
			t.Errorf("expected %d, but got %d at test: %s", test.expectedStatusCode, resp.StatusCode, test.name)
		}
	}
}

var TestsReservation = []struct {
	name               string
	reservation        interface{}
	expectedStatusCode int
}{
	{"OK", model.Reservation{
		RoomId: 1,
		Room: model.Room{
			Id:       1,
			RoomName: "Alter-doctor",
		},
	}, http.StatusOK},
	{"Without reservation", nil, http.StatusTemporaryRedirect},
	{"Wrong room id", model.Reservation{
		RoomId: 3,
		Room: model.Room{
			Id:       3,
			RoomName: "Alter-doctor",
		},
	}, http.StatusTemporaryRedirect},
}

func TestHandle_Reservation(t *testing.T) {

	for _, test := range TestsReservation {
		rr, err := MakeRequest("GET", "/make-reservation", nil, test.reservation, Handle.Reservation, nil)
		if err != nil {
			t.Error(err)
		}
		if rr.Code != test.expectedStatusCode {
			t.Errorf("Reservation hanler test: %s failed got: %d, expected: %d", test.name, rr.Code, test.expectedStatusCode)
		}
	}
}

var defaultReservationPOSTUrl = url.Values{
	"first_name": []string{"John"},
	"last_name":  []string{"Smith"},
	"email":      []string{"email.email@email.com"},
	"phone":      []string{"987-654-321-0"},
}

var defaultReservationPOSTReservation = model.Reservation{
	StartDate: time.Now().Add(time.Hour * 24),
	EndDate:   time.Now().Add(time.Hour * 48),
	RoomId:    1,
	Room: model.Room{
		Id:       1,
		RoomName: "Alter-doctor",
	},
}

var TestsReservationPOST = []struct {
	name               string
	reservation        interface{}
	reqBody            url.Values
	expectedStatusCode int
}{
	{"OK", defaultReservationPOSTReservation, defaultReservationPOSTUrl, http.StatusSeeOther},
	{"Missing form", defaultReservationPOSTReservation, nil, http.StatusTemporaryRedirect},
	{"Missing reservation", nil, defaultReservationPOSTUrl, http.StatusTemporaryRedirect},
	{"Invalid first_name form field", defaultReservationPOSTReservation,
		url.Values{
			"first_name": []string{""},
			"last_name":  []string{"Smith"},
			"email":      []string{"email.email@email.com"},
			"phone":      []string{"987-654-321-0"},
		}, http.StatusTemporaryRedirect},
	{"Invalid email form field", defaultReservationPOSTReservation,
		url.Values{
			"first_name": []string{"John"},
			"last_name":  []string{"Smith"},
			"email":      []string{"email"},
			"phone":      []string{"987-654-321-0"},
		}, http.StatusTemporaryRedirect},
	{"Invalid dates in reservation", model.Reservation{
		StartDate: time.Now().Truncate(time.Hour * 48),
		EndDate:   time.Now().Truncate(time.Hour * 24),
		RoomId:    1,
		Room: model.Room{
			Id:       1,
			RoomName: "Alter-doctor",
		},
	}, defaultReservationPOSTUrl, http.StatusTemporaryRedirect},
	{"Invalid room in reservation", model.Reservation{
		StartDate: time.Now().Add(time.Hour * 24),
		EndDate:   time.Now().Add(time.Hour * 48),
		RoomId:    0,
		Room: model.Room{
			Id: 0,
		},
	}, defaultReservationPOSTUrl, http.StatusTemporaryRedirect},
	{"Insert room restriction fail", model.Reservation{
		StartDate: time.Now().Add(time.Hour * 24),
		EndDate:   time.Now().Add(time.Hour * 48),
		RoomId:    1999,
		Room: model.Room{
			Id:       1999,
			RoomName: "Alter-doctor",
		},
	}, defaultReservationPOSTUrl, http.StatusTemporaryRedirect},
}

func TestHandle_ReservationPOST(t *testing.T) {

	for _, test := range TestsReservationPOST {
		rr, err := MakeRequest("POST", "/make-reservation", test.reqBody, test.reservation, Handle.PostReservation, nil)
		if err != nil {
			t.Error(err)
		}
		if rr.Code != test.expectedStatusCode {
			t.Errorf("ReservationPOST hanler test: %s failed got: %d, expected: %d", test.name, rr.Code, test.expectedStatusCode)
		}
	}
}

var defaultCheckAvailabilityUrl = url.Values{
	"start": []string{time.Now().Add(time.Hour * 24).Format(model.DateLayout)},
	"end":   []string{time.Now().Add(time.Hour * 48).Format(model.DateLayout)},
}

var TestsCheckAvailability = []struct {
	name               string
	reqBody            url.Values
	expectedStatusCode int
}{
	{"OK", defaultCheckAvailabilityUrl, http.StatusOK},
	{"Missing form", nil, http.StatusInternalServerError},
	{"Invalid start date", url.Values{
		"start": []string{"start"},
		"end":   []string{time.Now().Add(time.Hour * 48).Format(model.DateLayout)},
	}, http.StatusInternalServerError},
	{"Invalid end date", url.Values{
		"start": []string{time.Now().Add(time.Hour * 24).Format(model.DateLayout)},
		"end":   []string{"end"},
	}, http.StatusInternalServerError},
	{"Unable search availability", url.Values{
		"start": []string{time.Now().Truncate(time.Hour * 24).Format(model.DateLayout)},
		"end":   []string{time.Now().Add(time.Hour * 48).Format(model.DateLayout)},
	}, http.StatusInternalServerError},
	{"Zero available rooms", url.Values{
		"start": []string{time.Now().Add(time.Hour * 24).Format(model.DateLayout)},
		"end":   []string{time.Now().Add(time.Hour * 24).Format(model.DateLayout)},
	}, http.StatusSeeOther},
}

func TestHandler_CheckAvailability(t *testing.T) {

	for _, test := range TestsCheckAvailability {
		rr, err := MakeRequest("POST", "/search-availability", test.reqBody, nil, Handle.CheckAvailability, nil)
		if err != nil {
			t.Error(err)
		}

		if rr.Code != test.expectedStatusCode {
			t.Errorf("CheckAvailability handler test: %s failed got: %d, expected: %d", test.name, rr.Code, test.expectedStatusCode)
		}
	}
}

var defaultCheckAvailabilityJSONUrl = url.Values{
	"start":   []string{time.Now().Add(time.Hour * 24).Format(model.DateLayout)},
	"end":     []string{time.Now().Add(time.Hour * 48).Format(model.DateLayout)},
	"room_id": []string{"1"},
}

var TestsCheckAvailabilityJSON = []struct {
	name               string
	reqBody            url.Values
	ok                 bool
	expectedStatusCode int
}{
	{"OK", defaultCheckAvailabilityJSONUrl, true, http.StatusOK},
	{"Missing form", nil, false, http.StatusInternalServerError},
	{"Invalid start date", url.Values{
		"start":   []string{"start"},
		"end":     []string{time.Now().Add(time.Hour * 48).Format(model.DateLayout)},
		"room_id": []string{"1"},
	}, false, http.StatusInternalServerError},
	{"Invalid end date", url.Values{
		"start":   []string{time.Now().Add(time.Hour * 24).Format(model.DateLayout)},
		"end":     []string{"end"},
		"room_id": []string{"1"},
	}, false, http.StatusInternalServerError},
	{"Invalid room id", url.Values{
		"start":   []string{time.Now().Add(time.Hour * 24).Format(model.DateLayout)},
		"end":     []string{time.Now().Add(time.Hour * 48).Format(model.DateLayout)},
		"room_id": []string{"id"},
	}, false, http.StatusInternalServerError},
	{"Unable search availability", url.Values{
		"start":   []string{time.Now().Add(time.Hour * 24).Format(model.DateLayout)},
		"end":     []string{time.Now().Add(time.Hour * 48).Format(model.DateLayout)},
		"room_id": []string{"1999"},
	}, false, http.StatusInternalServerError},
}

func TestHandler_CheckAvailabilityJSON(t *testing.T) {

	var responseBody model.JsonResponse

	for _, test := range TestsCheckAvailabilityJSON {
		rr, err := MakeRequest("POST", "/search-availability-json", test.reqBody, nil, Handle.CheckAvailabilityJSON, nil)
		if err != nil {
			t.Error(err)
		}
		if err = json.NewDecoder(rr.Body).Decode(&responseBody); err != nil {
			t.Errorf("can't decode response body: %v", err)
		}

		if responseBody.Ok != test.ok || rr.Code != test.expectedStatusCode {
			t.Errorf("CheckAvailabilityJSON handler test: %s failed got: %d, expected: %d", test.name, rr.Code, test.expectedStatusCode)
		}
	}
}

var defailtReservationSummaryReservation = model.Reservation{
	StartDate: time.Now().Add(time.Hour * 24),
	EndDate:   time.Now().Add(time.Hour * 48),
}

var TestsReservationSummary = []struct {
	name               string
	reservation        interface{}
	expectedStatusCode int
}{
	{"OK", defailtReservationSummaryReservation, http.StatusOK},
	{"Missing reservation", nil, http.StatusTemporaryRedirect},
}

func TestHandler_ReservationSummary(t *testing.T) {

	for _, test := range TestsReservationSummary {
		rr, err := MakeRequest("GET", "/reservation-summary", nil, test.reservation, Handle.ReservationSummary, nil)
		if err != nil {
			t.Error(err)
		}

		if rr.Code != test.expectedStatusCode {
			t.Errorf("ReservationSummary handler test: %s failed got: %d, expected: %d", test.name, rr.Code, test.expectedStatusCode)
		}
	}
}

var TestsChooseRoom = []struct {
	name               string
	reservation        interface{}
	id                 interface{}
	expectedStatusCode int
}{
	{"OK", model.Reservation{}, 1, http.StatusSeeOther},
	{"Missing id", model.Reservation{}, nil, http.StatusInternalServerError},
	{"Missing reservation", nil, 1, http.StatusInternalServerError},
}

func TestHandler_ChooseRoom(t *testing.T) {

	for _, test := range TestsChooseRoom {
		rr, err := MakeRequest("GET", "/choose-room/{id}", nil, test.reservation, Handle.ChooseRoom, test.id)
		if err != nil {
			t.Error(err)
		}

		if rr.Code != test.expectedStatusCode {
			t.Errorf("ChooseRoom handler test: %s failed got: %d, expected: %d", test.name, rr.Code, test.expectedStatusCode)
		}
	}
}

var TestsBookRoom = []struct {
	name               string
	query              string
	expectedStatusCode int
}{
	{"OK", "?id=1&s=2023-01-02&e=2023-01-03", http.StatusSeeOther},
	{"Invalid id", "?id=e&s=2023-01-02&e=2023-01-03", http.StatusInternalServerError},
	{"Invalid start date", "?id=1&s=start&e=2023-01-03", http.StatusInternalServerError},
	{"Invalid end date", "?id=1&s=2023-01-02&e=end", http.StatusInternalServerError},
	{"Unable get room by id", "?id=3&s=2023-01-02&e=2023-01-03", http.StatusInternalServerError},
}

func TestHandler_BookRoom(t *testing.T) {

	for _, test := range TestsBookRoom {
		rr, err := MakeRequest("GET", "/book-room", nil, nil, Handle.BookRoom, test.query)
		if err != nil {
			t.Error(err)
		}

		if rr.Code != test.expectedStatusCode {
			t.Errorf("BookRoom handler test: %s failed got: %d, expected: %d", test.name, rr.Code, test.expectedStatusCode)
		}
	}
}

func MakeRequest(method, url string, reqBody url.Values, reservation interface{}, handlerFunc http.HandlerFunc, query interface{}) (*httptest.ResponseRecorder, error) {
	var (
		req *http.Request
	)
	router := mux.NewRouter()
	router.HandleFunc(url, handlerFunc).Methods(method)
	if query != nil {
		switch query.(type) {
		case int:
			url = strings.Replace(url, "{id}", fmt.Sprintf("%d", query.(int)), 1)
		case string:
			url = url + query.(string)
		}
	}
	if reqBody == nil {
		req, _ = http.NewRequest(method, url, nil)
	} else {
		req, _ = http.NewRequest(method, url, strings.NewReader(reqBody.Encode()))
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	ctx, err := getCtx(req)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)

	if reservation != nil {
		mockApp.Session.Put(ctx, "reservation", reservation)
	}

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)
	//handlerFunc.ServeHTTP(rr, req)

	return rr, nil
}

func getCtx(r *http.Request) (context.Context, error) {
	ctx, err := mockApp.Session.Load(r.Context(), "X-Session")
	if err != nil {
		return nil, err
	}
	return ctx, nil

}

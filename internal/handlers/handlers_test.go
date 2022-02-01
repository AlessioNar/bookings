package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/tsawler/bookings-app/internal/models"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name                string
	url                 string
	method              string
	expecetedStatusCode int
}{
	{"home", "/", "GET", http.StatusOK},
	{"about", "/about", "GET", http.StatusOK},
	{"gq", "/generals-quarters", "GET", http.StatusOK},
	{"ms", "/majors-suite", "GET", http.StatusOK},
	{"sa", "/search-availability", "GET", http.StatusOK},
	{"contact", "/contact", "GET", http.StatusOK},

	//{"post-search-availability", "/search-availability", "POST", []postData{
	//	{key: "start", value: "2022-01-29"},
	//	{key: "end", value: "2022-02-01"},
	//}, http.StatusOK},

	//{"post-search-avail-json", "/search-availability-json", "POST", []postData{
	//	{key: "start", value: "2022-01-29"},
	//	{key: "end", value: "2022-02-01"},
	//}, http.StatusOK},
	//{"make-reservation-post", "/make-reservation", "POST", []postData{
	//	{key: "first_name", value: "John"},
	//	{key: "last_name", value: "Smith"},
	//	{key: "email", value: "john.smith@email.org"},
	//	{key: "phone", value: "555-555-5555"},
	//}, http.StatusOK},
}

func TestHandler(t *testing.T) {
	routes := getRoutes()
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, e := range theTests {
		resp, err := ts.Client().Get(ts.URL + e.url)
		if err != nil {
			t.Log(err)
			t.Fatal(err)
		}

		if resp.StatusCode != e.expecetedStatusCode {
			t.Errorf("for %s, expected %d but got %d", e.name, e.expecetedStatusCode, resp.StatusCode)
		}

	}
}

func TestRepository_Reservation(t *testing.T) {

	// Initialize reservation
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
	}

	req, _ := http.NewRequest("GET", "/make-reservation", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.Reservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Reservation handler return wrong response code: got %d wanted %d", rr.Code, http.StatusOK)
	}

	// test case where reservation is not in session (reset everything)
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler return wrong response code: got %d wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test with non-existing room
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	// Set RoomID to 100
	reservation.RoomID = 100
	session.Put(ctx, "reservation", reservation)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler return wrong response code: got %d wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

}

func TestRepository_PostReservation(t *testing.T) {

	// Generate string and pass it to new request
	reqBody := "start_date=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=John")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=john.smith@email.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=123456789")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	req, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler return wrong response code: got %d wanted %d", rr.Code, http.StatusSeeOther)
	}

	// Test for missing post body
	req, _ = http.NewRequest("POST", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler return wrong response code for missing post body: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// Test for invalid start date
	reqBody = "start_date=invalid"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=John")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=john.smith@email.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=123456789")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler return wrong response code for invalid start date: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// Test for invalid end date
	reqBody = "start_date=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=invalid")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=John")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=john.smith@email.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=123456789")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler return wrong response code for invalid end date: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// Test for invalid roomID
	reqBody = "start_date=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=John")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=john.smith@email.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=123456789")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=invalid")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler return wrong response code for invalid room id: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// Test for invalid data
	reqBody = "start_date=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=J")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=john.smith@email.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=123456789")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler return wrong response code for invalid data: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// Test for failure to insert reservation into database
	reqBody = "start_date=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=John")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=john.smith@email.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=123456789")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=2")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler failed when trying to insert reservation in DB: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// Test for failure to insert restriction into database
	reqBody = "start_date=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=John")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=john.smith@email.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=123456789")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1000")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler failed when trying to insert restriction in DB: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

}

func getCtx(req *http.Request) context.Context {
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return ctx
}

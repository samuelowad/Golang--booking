package handlers

import (
	"context"
	"fmt"
	"github.com/samuelowad/bookings/internal/models"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name               string
	url                string
	method             string
	expectedStatusCode int
}{
	{"home", "/", "GET", http.StatusOK},
	{"about", "/about", "GET", http.StatusOK},
	//{"post-search-available", "/search-ava", "POST", []postData{
	//	{key: "start", value: "2013-01-01"},
	//	{key: "end", value: "2013-01-04"},
	//}, http.StatusOK},
	//{"post-search-available-json", "/search-ava-json", "POST", []postData{
	//	{key: "first_name", value: "test"},
	//	{key: "last_name", value: "test"},
	//	{key: "email", value: "test@test.com"},
	//	{key: "phone", value: "+2348886645658"},
	//}, http.StatusOK},
}

func TestHandler(t *testing.T) {
	routes := getRoute()
	testServer := httptest.NewTLSServer(routes)

	defer testServer.Close()

	for _, e := range theTests {
		res, err := testServer.Client().Get(testServer.URL + e.url)
		if err != nil {
			t.Log(err)
			t.Fatal(err)
		}

		if res.StatusCode != e.expectedStatusCode {
			t.Errorf("for  %s, exxpected %d, got %d", e.name, e.expectedStatusCode, res.StatusCode)
		}
	}
}

func TestRepository_Reservation(t *testing.T) {
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General",
		},
	}

	req, _ := http.NewRequest("GET", "/make-reservation", nil)
	ctx := getCtx(req)

	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.MakeReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("reservation handler returned %d, instead of %d", rr.Code, http.StatusOK)
	}

	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)

	req = req.WithContext(ctx)

	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("reservation handler returned %d, instead of %d", rr.Code, http.StatusOK)
	}

	//test with non-existing room
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)

	req = req.WithContext(ctx)

	rr = httptest.NewRecorder()
	reservation.RoomID = 23
	session.Put(ctx, "reservation", reservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("reservation handler returned %d, instead of %d", rr.Code, http.StatusOK)
	}

}

func TestRepository_PostReservation(t *testing.T) {
	cd, _ := convertDateStringToTime("2050-01-01", "2050-01-02")
	reservation := models.Reservation{
		StartDate: cd.startDate,
		EndDate:   cd.endDate,

		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General",
		},
	}
	reqBody := "last_name=sam"

	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=sam")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=sam@sam.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=1234567890")
	req, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx := getCtx(req)

	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-Form-urlencoded")
	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)
	handler := http.HandlerFunc(Repo.PostMakeReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("post reservation handler returned %d, instead of %d", rr.Code, http.StatusOK)
	}

	//	test for missing session
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)

	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-Form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostMakeReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("post reservation handler returned %d, instead of %d for missing session", rr.Code, http.StatusOK)
	}

	//	test for  invaild user data
	reqBody = "last_name=sam"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=sam")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=sam")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=1234567890")
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)

	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-Form-urlencoded")
	rr = httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)
	handler = http.HandlerFunc(Repo.PostMakeReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("post reservation handler returned %d, instead of %d for invalid userdata", rr.Code, http.StatusOK)
	}

	//	test for failing to insert reservation
	reqBody = "last_name=sam"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=sam")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=sam@sam.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=1234567890")
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)

	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-Form-urlencoded")
	rr = httptest.NewRecorder()

	//new room ID
	reservation.RoomID = 2
	session.Put(ctx, "reservation", reservation)
	handler = http.HandlerFunc(Repo.PostMakeReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("post reservation handler returned %d, instead of %d for inserting user reservation", rr.Code, http.StatusSeeOther)
	}

	//	test for failing to insert room restriction
	reqBody = "last_name=sam"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=sam")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=sam@sam.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=1234567890")
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)

	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-Form-urlencoded")
	rr = httptest.NewRecorder()

	reservation.RoomID = 1000
	session.Put(ctx, "reservation", reservation)
	handler = http.HandlerFunc(Repo.PostMakeReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("post reservation handler returned %d, instead of %d for inserting user reservation", rr.Code, http.StatusSeeOther)
	}
}

func getCtx(req *http.Request) context.Context {
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return ctx
}

package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/samuelowad/bookings/src/driver"
	"github.com/samuelowad/bookings/src/helpers"
	"github.com/samuelowad/bookings/src/repository"
	"github.com/samuelowad/bookings/src/repository/dbrepo"
	"net/http"
	"strconv"
	"time"

	"github.com/samuelowad/bookings/src/utils"

	"github.com/samuelowad/bookings/src/config"
	"github.com/samuelowad/bookings/src/models"
	"github.com/samuelowad/bookings/src/render"
)

var Repo *Repository

type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

//NewRepo creates a new repository
func NewRepo(app *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: app,
		DB:  dbrepo.NewPostgresDBRepo(db.SQL, app),
	}
}

func NewHandler(r *Repository) {
	Repo = r
}

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.Template(w, "home.page.tmpl", &models.TemplateData{}, r)
	//fmt.Fprintf(w, "hello")

}

func (m *Repository) About(w http.ResponseWriter, r *http.Request) {

	render.Template(w, "about.page.tmpl", &models.TemplateData{}, r)

}

func (m *Repository) Reserve(w http.ResponseWriter, r *http.Request) {

	render.Template(w, "about.page.tmpl", &models.TemplateData{}, r)

}

func (m *Repository) General(w http.ResponseWriter, r *http.Request) {

	render.Template(w, "general.page.tmpl", &models.TemplateData{}, r)

}

func (m *Repository) Executive(w http.ResponseWriter, r *http.Request) {

	render.Template(w, "executive.page.tmpl", &models.TemplateData{}, r)

}

func (m *Repository) SearchAva(w http.ResponseWriter, r *http.Request) {

	render.Template(w, "search-ava.page.tmpl", &models.TemplateData{}, r)

}

// post search ava
func (m *Repository) PostSearchAva(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")
	fmt.Println()

	w.Write([]byte(fmt.Sprintf("start date %s and end date is %s", start, end)))
}

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

func (m *Repository) AvailabilityJson(w http.ResponseWriter, r *http.Request) {
	resp := jsonResponse{
		OK:      true,
		Message: "Availabile",
	}

	out, err := json.MarshalIndent(resp, "", "   ")
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.Write(out)

}

func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {

	render.Template(w, "contact.page.tmpl", &models.TemplateData{}, r)

}

func (m *Repository) MakeReservation(w http.ResponseWriter, r *http.Request) {
	var empyReservation models.Reservation

	data := make(map[string]interface{})
	data["reservation"] = empyReservation

	render.Template(w, "make-res.page.tmpl", &models.TemplateData{
		Form: utils.New(nil),
		Data: data,
	}, r)

}

func (m *Repository) PostMakeReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	startDate := r.Form.Get("start_date")
	endDate := r.Form.Get("end_date")
	dateLayout := "2006-01-02"
	NewStartDate, err := time.Parse(dateLayout, startDate)
	if err != nil {
		helpers.ServerError(w, err)
		return

	}
	newEndDate, err := time.Parse(dateLayout, endDate)
	if err != nil {
		helpers.ServerError(w, err)
		return

	}
	roomID, err := strconv.Atoi(r.Form.Get("room_id"))

	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Phone:     r.Form.Get("phone"),
		Email:     r.Form.Get("email"),
		StartDate: NewStartDate,
		EndDate:   newEndDate,
		RoomID:    roomID,
	}

	form := utils.New(r.PostForm)

	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3, r)
	form.IsEmail("email")
	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		render.Template(w, "make-res.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		}, r)
		return
	}
	m.App.Session.Put(r.Context(), "reservation", reservation)

	newReservationID, err := m.DB.InsertReservation(reservation)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	restriction := models.RoomRestriction{
		StartDate:     NewStartDate,
		EndDate:       newEndDate,
		RoomID:        roomID,
		ReservationID: newReservationID,
		RestrictionID: 1,
	}
	err = m.DB.InsertRoomRestriction(restriction)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	http.Redirect(w, r, "/res-summary", http.StatusSeeOther)

}

func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, found := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)

	if !found {
		m.App.ErrorLog.Println("cant get error from session")
		m.App.Session.Put(r.Context(), "error", "cannot get Reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	m.App.Session.Remove(r.Context(), "reservation")
	data := make(map[string]interface{})
	data["reservation"] = reservation

	render.Template(w, "res-summary.page.tmpl", &models.TemplateData{
		Data: data,
	}, r)

}

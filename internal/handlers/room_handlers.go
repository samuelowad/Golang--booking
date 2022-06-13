package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/samuelowad/bookings/internal/driver"
	"github.com/samuelowad/bookings/internal/helpers"
	"github.com/samuelowad/bookings/internal/repository"
	"github.com/samuelowad/bookings/internal/repository/dbrepo"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/samuelowad/bookings/internal/utils"

	"github.com/samuelowad/bookings/internal/config"
	"github.com/samuelowad/bookings/internal/models"
	"github.com/samuelowad/bookings/internal/render"
)

var Repo *Repository

type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

type dateTime struct {
	startDate, endDate time.Time
}

//NewRepo creates a new repository
func NewRepo(app *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: app,
		DB:  dbrepo.NewPostgresDBRepo(db.SQL, app),
	}
}

//NewTestRepo creates a new repository
func NewTestRepo(app *config.AppConfig) *Repository {
	return &Repository{
		App: app,
		DB:  dbrepo.NewTestingRepo(app),
	}
}

func NewHandler(r *Repository) {
	Repo = r
}

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.Template(w, "home.page.tmpl", &models.TemplateData{}, r)

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

	date, err := convertDateStringToTime(start, end)

	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	rooms, err := m.DB.SearchAvailabilityForAllRooms(date.startDate, date.endDate)

	if err != nil {

		log.Println(err)
		helpers.ServerError(w, err)
		return
	}

	if len(rooms) == 0 {
		m.App.Session.Put(r.Context(), "error", "No Availability")
		http.Redirect(w, r, "/search-ava", http.StatusSeeOther)
	}
	data := make(map[string]interface{})
	data["rooms"] = rooms

	res := models.Reservation{

		StartDate: date.startDate,
		EndDate:   date.endDate,
	}
	m.App.Session.Put(r.Context(), "reservation", res)

	render.Template(w, "choose-room.page.tmpl", &models.TemplateData{
		Data: data,
	}, r)
}

type jsonResponse struct {
	OK        bool   `json:"ok"`
	Message   string `json:"message"`
	RoomID    string `json:"room_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

func (m *Repository) AvailabilityJson(w http.ResponseWriter, r *http.Request) {

	roomID, _ := strconv.Atoi(r.Form.Get("room_id"))
	sd := r.Form.Get("start")
	ed := r.Form.Get("end")
	date, err := convertDateStringToTime(sd, ed)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	available, err := m.DB.SearchAvailabilityByDateByRoomID(date.startDate, date.endDate, roomID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	resp := jsonResponse{
		OK:        available,
		Message:   "",
		StartDate: sd,
		EndDate:   ed,
		RoomID:    strconv.Itoa(roomID),
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
	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "cannot get reservation from Session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	room, err := m.DB.GetRoomByID(res.RoomID)

	if err != nil {
		m.App.Session.Put(r.Context(), "error", "cannot find room")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	res.Room.RoomName = room.RoomName

	m.App.Session.Put(r.Context(), "reservation", res)

	startDate := res.StartDate.Format("2006-01-02")
	endDate := res.EndDate.Format("2006-01-02")
	stringMap := make(map[string]string)
	stringMap["start_date"] = startDate
	stringMap["end_date"] = endDate

	data := make(map[string]interface{})
	data["reservation"] = res

	render.Template(w, "make-res.page.tmpl", &models.TemplateData{
		Form:      utils.New(nil),
		Data:      data,
		StringMap: stringMap,
	}, r)

}

func (m *Repository) PostMakeReservation(w http.ResponseWriter, r *http.Request) {

	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)

	if !ok {
		m.App.Session.Put(r.Context(), "error", "cannot get reservation from Session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "cannot get reservation from Session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	reservation.FirstName = r.Form.Get("first_name")
	reservation.LastName = r.Form.Get("last_name")
	reservation.Phone = r.Form.Get("phone")
	reservation.Email = r.Form.Get("email")

	form := utils.New(r.PostForm)

	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3, r)
	form.IsEmail("email")
	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation
		http.Error(w, "error", http.StatusSeeOther)
		render.Template(w, "make-res.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		}, r)
		return
	}

	newReservationID, err := m.DB.InsertReservation(reservation)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "cannot create reservation")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	m.App.Session.Put(r.Context(), "reservation", reservation)

	restriction := models.RoomRestriction{
		StartDate:     reservation.StartDate,
		EndDate:       reservation.EndDate,
		RoomID:        reservation.RoomID,
		ReservationID: newReservationID,
		RestrictionID: 1,
	}
	err = m.DB.InsertRoomRestriction(restriction)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "cannot create reservation ")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	//send notification to client
	mailMessage := fmt.Sprintf(`<strong>Reservation Confirmation</strong><br>
			dear %s,<br>
				this is to comfirm your reservation from %s to %s.`, reservation.FirstName, reservation.StartDate.Format("2006-01-02"), reservation.EndDate.Format("2006-01-02"))
	msg := models.MailData{reservation.Email, "test@test.com", "Reservation comfirm", mailMessage, "email.html"}

	m.App.MailChan <- msg

	//send notification to home owner
	mailMessage = fmt.Sprintf(`<strong>Reservation Notification</strong><br>
			dear %s,<br>
				this is to inform you of a new reservation from %s to %s.`, reservation.FirstName, reservation.StartDate.Format("2006-01-02"), reservation.EndDate.Format("2006-01-02"))
	msg = models.MailData{"homeowner@me.com", "test@test.com", "Reservation notification", mailMessage, "email.html"}

	m.App.MailChan <- msg

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
	sd := reservation.StartDate.Format("2006-01-02")
	ed := reservation.EndDate.Format("2006-01-02")

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	render.Template(w, "res-summary.page.tmpl", &models.TemplateData{
		Data:      data,
		StringMap: stringMap,
	}, r)

}

func (m *Repository) ChooseRoom(w http.ResponseWriter, r *http.Request) {
	roomId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, err)
		return
	}

	res.RoomID = roomId
	m.App.Session.Put(r.Context(), "reservation", res)

	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}

//BookRoom books room, builds session and redirect to reservation
func (m *Repository) BookRoom(w http.ResponseWriter, r *http.Request) {

	roomID, _ := strconv.Atoi(r.URL.Query().Get("id"))
	start := r.URL.Query().Get("s")
	end := r.URL.Query().Get("e")

	date, err := convertDateStringToTime(start, end)
	room, err := m.DB.GetRoomByID(roomID)

	if err != nil {
		helpers.ServerError(w, errors.New("cannot get reservation"))
		return
	}

	var res models.Reservation

	res.Room.RoomName = room.RoomName
	res.RoomID = roomID
	res.StartDate = date.startDate
	res.EndDate = date.endDate

	m.App.Session.Put(r.Context(), "reservation", res)

	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}

func convertDateStringToTime(startDate, endDate string) (dateTime, error) {
	dateLayout := "2006-01-02"
	NewStartDate, err := time.Parse(dateLayout, startDate)
	if err != nil {
		return dateTime{}, err

	}
	newEndDate, err := time.Parse(dateLayout, endDate)
	if err != nil {
		return dateTime{}, err
	}

	return dateTime{NewStartDate, newEndDate}, nil
}

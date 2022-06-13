package handlers

import (
	"github.com/samuelowad/bookings/internal/helpers"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/samuelowad/bookings/internal/models"
	"github.com/samuelowad/bookings/internal/render"
	"github.com/samuelowad/bookings/internal/utils"
)

func (m *Repository) Login(w http.ResponseWriter, r *http.Request) {
	err := render.Template(w, "login.page.tmpl", &models.TemplateData{Form: utils.New(nil)}, r)
	if err != nil {
		return
	}
}

//PostLogin handles the login
func (m *Repository) PostLogin(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.RenewToken(r.Context())
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	form := utils.New(r.PostForm)
	form.Required("email", "password")
	form.IsEmail("email")
	if !form.Valid() {
		render.Template(w, "login.page.tmpl", &models.TemplateData{
			Form: form,
		}, r)
		return
	}
	id, _, err := m.DB.Authenticate(email, password)

	if err != nil {
		log.Println(err)
		m.App.Session.Put(r.Context(), "error", "invalid login")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return

	}
	m.App.Session.Put(r.Context(), "user_id", id)
	m.App.Session.Put(r.Context(), "flash", "login successful")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (m *Repository) Logout(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.Destroy(r.Context())
	_ = m.App.Session.RenewToken(r.Context())
	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func (m *Repository) AdminDashboard(w http.ResponseWriter, r *http.Request) {
	render.Template(w, "admin-dashboard.page.tmpl", &models.TemplateData{}, r)
}

func (m *Repository) AdminNewReservations(w http.ResponseWriter, r *http.Request) {

	reservations, err := m.DB.AllNewReservations()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	data := make(map[string]interface{})
	data["reservations"] = reservations

	render.Template(w, "admin-new-reservations.page.tmpl", &models.TemplateData{
		Data: data,
	}, r)
}

func (m *Repository) AdminAllReservations(w http.ResponseWriter, r *http.Request) {
	reservations, err := m.DB.AllReservations()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	data := make(map[string]interface{})
	data["reservations"] = reservations
	render.Template(w, "admin-all-reservations.page.tmpl", &models.TemplateData{Data: data}, r)
}

//AdminShowReservations shows the reservation in the admin tool
func (m *Repository) AdminShowReservations(w http.ResponseWriter, r *http.Request) {
	exploded := strings.Split(r.RequestURI, "/")
	id, err := strconv.Atoi(exploded[4])

	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	src := exploded[3]

	stringMap := make(map[string]string)
	stringMap["src"] = src

	res, err := m.DB.GetReservationByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["reservation"] = res

	render.Template(w, "admin-reservations-show.page.tmpl", &models.TemplateData{
		StringMap: stringMap, Data: data, Form: utils.New(nil),
	}, r)

}

func (m *Repository) AdminReservationsCalendar(w http.ResponseWriter, r *http.Request) {
	render.Template(w, "admin-reservations-calender.page.tmpl", &models.TemplateData{}, r)
}

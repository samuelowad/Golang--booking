package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/samuelowad/bookings/src/config"
	"github.com/samuelowad/bookings/src/models"
	"github.com/samuelowad/bookings/src/render"
)

var Repo *Repository

type Repository struct {
	App *config.AppConfig
}

//NewRepo creates a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

func NewHandler(r *Repository) {
	Repo = r
}

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIp := r.RemoteAddr
	m.App.Session.Put(r.Context(), "user_ip", remoteIp)
	render.RenderTemplate(w, "home.page.tmpl", &models.TemplateData{}, r)
	//fmt.Fprintf(w, "hello")

}

func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	stringMap :=
		make(map[string]string)

	stringMap["test"] = "hello again"

	remoteIp := m.App.Session.GetString(r.Context(), "user_ip")

	stringMap["user_ip"] = remoteIp
	render.RenderTemplate(w, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	}, r)

}

func (m *Repository) Reserve(w http.ResponseWriter, r *http.Request) {

	render.RenderTemplate(w, "about.page.tmpl", &models.TemplateData{}, r)

}

func (m *Repository) General(w http.ResponseWriter, r *http.Request) {

	render.RenderTemplate(w, "general.page.tmpl", &models.TemplateData{}, r)

}

func (m *Repository) Executive(w http.ResponseWriter, r *http.Request) {

	render.RenderTemplate(w, "executive.page.tmpl", &models.TemplateData{}, r)

}

func (m *Repository) SearchAva(w http.ResponseWriter, r *http.Request) {

	render.RenderTemplate(w, "search-ava.page.tmpl", &models.TemplateData{}, r)

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
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")

	w.Write(out)

}

func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {

	render.RenderTemplate(w, "contact.page.tmpl", &models.TemplateData{}, r)

}

func (m *Repository) MakeReservation(w http.ResponseWriter, r *http.Request) {

	render.RenderTemplate(w, "make-res.page.tmpl", &models.TemplateData{}, r)

}

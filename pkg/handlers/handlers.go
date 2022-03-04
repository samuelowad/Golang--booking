package handlers

import (
	"github.com/samuelowad/bookings/pkg/config"
	"github.com/samuelowad/bookings/pkg/models"
	"github.com/samuelowad/bookings/pkg/render"
	"net/http"
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
	render.RenderTemplate(w, "home.page.tmpl", &models.TemplateData{})
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
	})

}

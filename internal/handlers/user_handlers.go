package handlers

import (
	"github.com/samuelowad/bookings/internal/models"
	"github.com/samuelowad/bookings/internal/render"
	"github.com/samuelowad/bookings/internal/utils"
	"net/http"
)

func (m *Repository) Login(w http.ResponseWriter, r *http.Request) {
	err := render.Template(w, "login.page.tmpl", &models.TemplateData{Form: utils.New(nil)}, r)
	if err != nil {
		return
	}
}

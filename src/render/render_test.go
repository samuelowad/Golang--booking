package render

import (
	"github.com/samuelowad/bookings/src/models"
	"net/http"
	"testing"
)

func TestAddDefaultData(t *testing.T) {
	var td models.TemplateData

	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	session.Put(r.Context(), "flash", "123")
	result := AddDefaultData(&td, r)

	if result.Flash != "123" {
		t.Error("Flash value of 123")
	}

}

func getSession() (*http.Request, error) {
	r, err := http.NewRequest("GET", "/some-url", nil)

	if err != nil {
		return nil, err
	}
	ctx := r.Context()

	ctx, _ = session.Load(ctx, r.Header.Get("X-Session"))
	r = r.WithContext(ctx)

	return r, nil

}

func TestRenderTemplate(t *testing.T) {
	pathToTemplate = "../../template"

	tc, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}
	app.TemplateCache = tc

	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	var ww myWriter

	//err = RenderTemplate(&ww, "home.page.tmpl", &models.TemplateData{}, r)
	//if err != nil {
	//	fmt.Println(err)
	//	t.Error("error writing template")
	//}

	err = RenderTemplate(&ww, "not.page.tmpl", &models.TemplateData{}, r)
	if err == nil {
		t.Error("rendered non existing page")
	}
}

func TestNewTemplates(t *testing.T) {
	NewTemplates(app)
}

func TestCreateTemplateCache(t *testing.T) {
	pathToTemplate = "../../template"
	_, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}
}

package render

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/justinas/nosurf"
	"github.com/samuelowad/bookings/src/config"
	"github.com/samuelowad/bookings/src/models"
)

var functions = template.FuncMap{}

var app *config.AppConfig

//set config
func NewTemplates(a *config.AppConfig) {
	app = a

}

func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.CSRFToken = nosurf.Token(r)
	return td
}

func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData, r *http.Request) {

	var tc map[string]*template.Template
	if app.UseCache {
		tc = app.TemplateCache

	} else {
		tc, _ = CreateTemplateCache()
	}

	t, ok := tc[tmpl]

	if !ok {
		log.Fatal("could not get template")
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td, r)
	_ = t.Execute(buf, td)

	_, err := buf.WriteTo(w)
	if err != nil {
		fmt.Println("error writing")
	}

}

//create template cache
func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob("./templates/*.page.tmpl")

	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err

		}
		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err

		}
		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache, err

			}
		}
		myCache[name] = ts
	}

	return myCache, nil

}

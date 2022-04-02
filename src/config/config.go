package config

import (
	"github.com/alexedwards/scs/v2"
	"html/template"
	"log"
)

//AppConfig holds application config
type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template

	InfoLog  *log.Logger
	ErrorLog *log.Logger
	InProd   bool
	Session  *scs.SessionManager
}

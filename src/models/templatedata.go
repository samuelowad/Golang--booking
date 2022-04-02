package models

import "github.com/samuelowad/bookings/src/utils"

//TemplateData holds data sent from handler to template
type TemplateData struct {
	StringMap map[string]string
	IntMap    map[string]int
	FloatMap  map[string]float32
	Data      map[string]interface{}
	CSRFToken string
	Flash     string
	Warning   string
	Error     string
	Form      *utils.Form
}

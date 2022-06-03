package utils

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"net/http"
	"net/url"
	"strings"
)

type Form struct {
	url.Values
	Error errors
}

//Valid returns true if no errors
func (f *Form) Valid() bool {
	return len(f.Error) == 0
}

//New init a form struct
func New(data url.Values) *Form {
	return &Form{data, errors(map[string][]string{})}
}

func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Error.Add(field, "required")
		}
	}
}

//Has checks if form field is in Post and not empty
func (f *Form) Has(field string, r *http.Request) bool {
	x := r.Form.Get(field)
	if x == "" {
		return false
	}
	return true
}

//MinLength checks if string has MinLength
func (f *Form) MinLength(field string, length int, r *http.Request) bool {
	x := r.Form.Get(field)
	if len(x) < length {
		f.Error.Add(field, fmt.Sprintf("THIS field MUST BE AT LEAST %d characters long", length))
		return false
	}
	return true
}

//IsEmail checks for vaild email address
func (f *Form) IsEmail(field string) {
	if !govalidator.IsEmail(f.Get(field)) {
		f.Error.Add(field, "Invalid email address")
	}
}

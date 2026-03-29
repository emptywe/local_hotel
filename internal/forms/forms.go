package forms

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"net/url"
	"strings"
)

// Form creates a custom form struct and embeds an url.Values object
type Form struct {
	url.Values
	Errors errors
}

// New initializes a form struct
func New(data url.Values) *Form {
	return &Form{
		Values: data,
		Errors: make(map[string][]string),
	}
}

// Has checks if form field is in post and not empty
func (f *Form) Has(field string) bool {
	if f.Get(field) == "" {
		return false
	}
	return true
}

// Valid returns true if there is no errors
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

// MinLength checks for string minimum length
func (f *Form) MinLength(field string, length int) bool {
	if len(f.Get(field)) < length {
		f.Errors.Add(field, fmt.Sprintf("this field must be at least %d characters long", length))
		return false
	}
	return true
}

// Required checks for required fields
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		if strings.TrimSpace(f.Get(field)) == "" {
			f.Errors.Add(field, "this field cannot be blank")
		}
	}
}

func (f *Form) IsEmail(field string) {
	if !govalidator.IsEmail(f.Get(field)) {
		f.Errors.Add(field, "invalid email address")
	}
}

package forms

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
)

// Form creates a custom form struct and embeds a url.Values object
type Form struct {
	url.Values
	Errors errors
}

// Valid returns true if there are no errors otherwise false
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

// New initializes a form struct
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// checks for required fields
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank.")
		}
	}
}

// Has checks if form field is in post and not empty
func (f *Form) Has(field string) bool {
	x := f.Get(field)
	if x == "" {
		return false
	}
	return true
}

// IsEmail checks for valid email address
func (f *Form) IsEmail(field string) {
	if !govalidator.IsEmail(f.Get(field)) {
		f.Errors.Add(field, "Invalid email address.")
	}
}

// LengthTest tests a field for minimum and maximum length
func (f *Form) LengthTest(field string, minLength, maxLength int) bool {
	x := f.Get(field)
	if len(x) < minLength {
		f.Errors.Add(field, fmt.Sprintf("Input must be at least %d characters long", minLength))
		return false
	} else {
		if len(x) > maxLength {
			f.Errors.Add(field, fmt.Sprintf("Input must less than %d characters long", maxLength))
			return false
		}
	}
	return true
}

func (f *Form) PasswordsMatch(p1, p2 string) bool {
	fmt.Println("PasswordMAtch", p1, p2)
	if p1 != p2 {
		f.Errors.Add("password_2", "Passwords do not match")
		return false
	}
	return true
}

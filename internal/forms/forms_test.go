package forms

import (
	"net/url"
	"testing"
)

// func TestForm_MinLength(t *testing.T) {
// 	r := httptest.NewRequest("POST", "/whatever", nil)
// 	form := New(r.PostForm)

// 	form.MinLength("x", 10)
// 	if form.Valid() {
// 		t.Error("form shows min length for non-existent field")
// 	}

// 	isError := form.Errors.Get("x")
// 	if isError == "" {
// 		t.Error("should have an error, but did not get one")
// 	}

// 	postedValues := url.Values{}
// 	postedValues.Add("some_field", "some value")
// 	form = New(postedValues)

// 	form.MinLength("some_field", 100)
// 	if form.Valid() {
// 		t.Error("shows minlength of 100 met when data is shorter")
// 	}

// 	postedValues = url.Values{}
// 	postedValues.Add("another_field", "abc123")
// 	form = New(postedValues)

// 	form.MinLength("another_field", 1)
// 	if !form.Valid() {
// 		t.Error("shows minlength of 1 is not met when it is")
// 	}

// 	isError = form.Errors.Get("another_field")
// 	if isError != "" {
// 		t.Error("should not have an error, but got one")
// 	}

// }
func TestForm_Valid(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)

	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("form shows valid when required fields are missing")
	}

	postedData = url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "a")
	postedData.Add("c", "a")

	form = New(postedData)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("form shows invalid when required fields are present")
	}
}

func TestFormHas(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)

	has := form.Has("whatever")
	if has {
		t.Error("form shows has field when it does not")
	}

	postedData = url.Values{}
	postedData.Add("a", "a")
	form = New(postedData)

	has = form.Has("a")
	if !has {
		t.Error("shows form does not have field when it should")
	}
}

func TestIsEmail(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)

	form.IsEmail("x")
	if form.Valid() {
		t.Error("form shows valid email for non-existent field")
	}

	postedData = url.Values{}
	postedData.Add("email", "abc")
	form = New(postedData)

	form.IsEmail("email")
	if form.Valid() {
		t.Error("form shows valid email for invalid input")
	}

	postedData = url.Values{}
	postedData.Add("email", "abc@abc.com")
	form = New(postedData)

	form.IsEmail("email")
	if !form.Valid() {
		t.Error("form shows invalid email for valid input")
	}
}

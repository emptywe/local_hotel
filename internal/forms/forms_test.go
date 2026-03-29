package forms

import (
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestNew(t *testing.T) {
	req := httptest.NewRequest("POST", "/", nil)
	if form := New(req.PostForm); form == nil {
		t.Error("nil form value")
	}
}

func TestHas(t *testing.T) {
	form := setupForm()
	if !form.Has("testKey") {
		t.Error("can't find added key")
	}
	if form.Has("non-exisingKey") {
		t.Error("found non-existing key")
	}
}

func TestValid(t *testing.T) {
	form := setupForm()
	if !form.Valid() {
		t.Error("invalidate valid form")
	}
	form.Errors.Add("testKey", "testError")
	if form.Valid() {
		t.Error("validate invalid form")
	}
}

func TestMinLength(t *testing.T) {
	form := setupForm()
	if !form.MinLength("testKey", 3) {
		t.Error("invalidate min length with valid length")
	}
	if form.MinLength("testKey", 10) {
		t.Error("validate min length with invalid length")
	}
}

func TestIsEmail(t *testing.T) {
	form := setupForm()
	form.Add("email", "test@test.com")
	form.IsEmail("email")
	if len(form.Errors) > 0 {
		t.Error("invalidate valid email address")
	}
	form.Set("email", "test-invalid")
	form.IsEmail("email")
	if len(form.Errors) <= 0 {
		t.Error("validate invalid email address")
	}
}

func TestRequired(t *testing.T) {
	form := setupForm()
	form.Required("testKey")
	if len(form.Errors) > 0 {
		t.Errorf("invalidate valid form with requiered fields")
	}
	form.Add("testNewField", "testField")
	form.Required("testKey", "testNewField")
	if len(form.Errors) > 0 {
		t.Errorf("invalidate valid form with requiered fields")
	}
	form.Required("testKey", "testNewField", "non-existTest")
	if len(form.Errors) <= 0 {
		t.Errorf("validate invalid form with requiered fields")
	}
}

func setupForm() *Form {
	req := httptest.NewRequest("POST", "/", nil)
	req.PostForm = url.Values{"testKey": []string{"testValue"}}
	return New(req.PostForm)
}

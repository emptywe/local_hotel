package forms

import "testing"

func TestErrorsGet(t *testing.T) {
	form := setupForm()
	if form.Errors.Get("testKey") != "" {
		t.Error("get error when should get empty string")
	}
	form.Errors["testKey"] = []string{"testError"}
	if form.Errors.Get("testKey") != "testError" {
		t.Error("get invalid error")
	}
}

func TestErrorsAdd(t *testing.T) {
	form := setupForm()
	form.Errors.Add("testKey", "testError")
	if len(form.Errors) <= 0 {
		t.Error("failed adding new error")
	}
	if form.Errors.Get("testKey") != "testError" {
		t.Error("get invalid error")
	}
}

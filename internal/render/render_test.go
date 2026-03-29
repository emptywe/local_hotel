package render

import (
	"net/http"
	"testing"

	"github.com/emptywe/local_hotel/model"
)

func TestNewRenderer(t *testing.T) {
	if renderer := NewRenderer(&mockApp); renderer == nil {
		t.Error("returned invalid renderer")
	}
}

func TestCreateTemplateCache(t *testing.T) {
	pathToTemplates = "./../../html-templates"
	r := NewRenderer(&mockApp)
	tc, err := r.CreateTemplateCache()
	if err != nil {
		t.Errorf("failed to create template cache %v", err)
	}
	if tc == nil {
		t.Errorf("nil template cache %v", err)
	}
}

func TestAddDefaultData(t *testing.T) {
	var td model.TemplateData
	r := NewRenderer(&mockApp)
	req, err := getSession()
	if err != nil {
		t.Errorf("%v", err)
	}
	mockSession.Put(req.Context(), "flash", "test")
	r.AddDefaultData(&td, req)
	if td.Flash != "test" {
		t.Error("flash value didn't pass through function")
	}
}

func TestRenderTemplates(t *testing.T) {
	pathToTemplates = "./../../html-templates"
	r := NewRenderer(&mockApp)
	tc, err := r.CreateTemplateCache()
	if err != nil {
		t.Errorf("can't create templates cache: %v", err)
	}
	mockApp.TemplateCache = tc
	req, err := getSession()
	if err != nil {
		t.Errorf("can't create mok request: %v", err)
	}

	var ww mockWriter

	if err = r.Template(&ww, req, "home.page.gohtml", &model.TemplateData{}); err != nil {
		t.Errorf("can't render template to response: %v", err)
	}
	mockApp.UseCache = true
	if err = r.Template(&ww, req, "home.page.gohtml", &model.TemplateData{}); err != nil {
		t.Errorf("can't render template to response when using cahce: %v", err)
	}

	if err = r.Template(&ww, req, "non-existing.page.gohtml", &model.TemplateData{}); err == nil {
		t.Error("try to render non-existing template to response")
	}
}

func getSession() (*http.Request, error) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		return nil, err
	}
	ctx := req.Context()
	ctx, _ = mockSession.Load(ctx, req.Header.Get("X-Session"))
	req = req.WithContext(ctx)
	return req, nil
}

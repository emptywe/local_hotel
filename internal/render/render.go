package render

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"time"

	"github.com/emptywe/local_hotel/config"
	"github.com/emptywe/local_hotel/model"
	"github.com/justinas/nosurf"
	"go.uber.org/zap"
)

var pathToTemplates = "./html-templates"

type Renderer struct {
	app       *config.AppConfig
	functions template.FuncMap
}

func NewRenderer(a *config.AppConfig) *Renderer {
	return &Renderer{
		app: a,
		functions: template.FuncMap{
			"simpleDate": SimpleDate,
			"formatDate": FormatDate,
			"iterate":    Iterate,
			"add":        Add,
		},
	}
}

func (r *Renderer) AddDefaultData(data *model.TemplateData, req *http.Request) {
	data.Flash = r.app.Session.PopString(req.Context(), "flash")
	data.Error = r.app.Session.PopString(req.Context(), "error")
	data.Warning = r.app.Session.PopString(req.Context(), "warning")
	data.CSRFToken = nosurf.Token(req)
	if r.app.Session.Exists(req.Context(), "user_id") {
		data.IsAuthenticated = true
	}
}

// Add - sum two numbers
func Add(a, b int) int {
	return a + b
}

// Iterate - returns slice of int from 0 to count
func Iterate(count int) (items []int) {
	for i := 0; i < count; i++ {
		items = append(items, i)
	}
	return
}

// SimpleDate - returns time in YYYY-MM-DD format
func SimpleDate(t time.Time) string {
	return t.Format("2006-01-02")
}

// FormatDate - returns time in YYYY-MM-DD format
func FormatDate(t time.Time, f string) string {
	return t.Format(f)
}

// Template renders template using html-templates/template
func (r *Renderer) Template(w http.ResponseWriter, req *http.Request, tmpl string, data *model.TemplateData) error {
	var tc map[string]*template.Template
	if r.app.UseCache {
		tc = r.app.TemplateCache
	} else {
		tc, _ = r.CreateTemplateCache()
	}

	t, ok := tc[tmpl]
	if !ok {
		return errors.New("can't get template from cache")
	}

	buf := new(bytes.Buffer)
	r.AddDefaultData(data, req)
	if err := t.Execute(buf, data); err != nil {
		return err
	}

	if _, err := buf.WriteTo(w); err != nil {
		return err
	}
	return nil
}

// CreateTemplateCache creates a template cache as a map
func (r *Renderer) CreateTemplateCache() (map[string]*template.Template, error) {
	tplCache := make(map[string]*template.Template)
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.gohtml", pathToTemplates))
	if err != nil {
		return tplCache, err
	}

	for _, page := range pages {

		name := filepath.Base(page)
		zap.S().Debugf("Load template: %s", page)
		ts, err := template.New(name).Funcs(r.functions).ParseFiles(page)
		if err != nil {
			return tplCache, err
		}

		match, err := filepath.Glob(fmt.Sprintf("%s/*.layout.gohtml", pathToTemplates))
		if err != nil {
			return tplCache, err
		}

		if len(match) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.gohtml", pathToTemplates))
			if err != nil {
				return tplCache, err
			}
		}

		tplCache[name] = ts
	}

	return tplCache, nil
}

//func RenderTemplateNew(w http.ResponseWriter, tmpl string) {
//	fl, err := pongo2.NewLocalFileSystemLoader("/html-templates/")
//	if err != nil{
//		fmt.P
//	}
//}

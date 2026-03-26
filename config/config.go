package config

import (
	"github.com/alexedwards/scs/v2"
	"github.com/prometheus/client_golang/prometheus"
	"gitlab.com/rad13/mylittlebooking/model"
	"html/template"
)

// AppConfig holds the application config
type AppConfig struct {
	UseCache      bool                          // UseCache enables to use templates cache
	TemplateCache map[string]*template.Template // TemplateCache - parsed templates repository
	InProduction  bool                          // InProduction defines production or developing mode
	Session       *scs.SessionManager           // Session define session settings
	MailChan      chan model.MailData           // MailChan chan for sending email messages
	Stats         *prometheus.CounterVec        // Stats page visitors counter
}

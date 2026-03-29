package model

import "github.com/emptywe/local_hotel/internal/forms"

// TemplateData holds data sent from handlers to templates
type TemplateData struct {
	StringMap       map[string]string      // StringMap - for string variables passing through template
	IntMap          map[string]int         // IntMap - for integer variables passing through template
	floatMap        map[string]float32     // floatMap - for float variables passing through template
	Data            map[string]interface{} // Data - for random variables passing through template
	CSRFToken       string                 // CSRFToken - for CSRF security token
	Flash           string                 // Flash - for quick messages
	Warning         string                 // Warning - for users warning
	Error           string                 // Error - for unexpected errors
	Form            *forms.Form            // Form - form data from template
	IsAuthenticated bool                   // IsAuthenticated - shows is user authenticated
}

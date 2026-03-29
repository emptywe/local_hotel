package helpers

import (
	"encoding/json"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/emptywe/local_hotel/config"
	"github.com/emptywe/local_hotel/model"
)

var app *config.AppConfig

// NewHelpers sets upa app config for helpers
func NewHelpers(a *config.AppConfig) {
	app = a
}

// ClientError - simple client side error
func ClientError(w http.ResponseWriter, statusCode int) {
	http.Error(w, http.StatusText(statusCode), statusCode)
}

// ServerError - simple server side error
func ServerError(w http.ResponseWriter, err error) {
	//trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// JsonErrorResponse - json formatted error response
func JsonErrorResponse(w http.ResponseWriter, msg string, statusCode int) {
	resp := model.JsonResponse{
		Ok:      false,
		Message: msg,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(resp)
}

// IsAuthenticated - returns true if authenticated by user_id
func IsAuthenticated(r *http.Request, ses *scs.SessionManager) bool {
	return ses.Exists(r.Context(), "user_id")
}

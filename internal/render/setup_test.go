package render

import (
	"encoding/gob"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/emptywe/local_hotel/config"
	"github.com/emptywe/local_hotel/model"
)

type mockWriter struct {
}

func (ww *mockWriter) Write(data []byte) (int, error) {
	return len(data), nil
}

func (ww *mockWriter) Header() http.Header {
	return http.Header{}
}

func (ww *mockWriter) WriteHeader(statusCode int) {

}

var (
	mockApp     config.AppConfig
	mockSession *scs.SessionManager
)

func TestMain(m *testing.M) {
	gob.Register(model.Reservation{})
	// for production change to true
	mockApp.InProduction = false
	setupSession(&mockApp)
	os.Exit(m.Run())
}

func setupSession(app *config.AppConfig) {
	mockSession = scs.New()
	mockSession.Lifetime = time.Hour
	mockSession.Cookie.Persist = true
	mockSession.Cookie.SameSite = http.SameSiteLaxMode
	mockSession.Cookie.Secure = app.InProduction
	mockApp.Session = mockSession
}

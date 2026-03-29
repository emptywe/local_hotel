package router

import (
	"testing"

	"github.com/emptywe/local_hotel/config"
	"github.com/gorilla/mux"
)

var mokApp config.AppConfig

func TestInitRoutes(t *testing.T) {

	router := InitRoutes(&mokApp)
	switch v := router.(type) {
	case *mux.Router:
	default:
		t.Errorf("%T is not *mux.Router type", v)
	}
}

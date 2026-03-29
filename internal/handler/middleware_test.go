package handler

import (
	"net/http"
	"testing"

	"github.com/emptywe/local_hotel/config"
)

func TestNoSurf(t *testing.T) {
	var (
		th testHandler
	)

	NewHandler(&mockApp, nil, nil)
	hd := Handle.NoSurf(&th)

	switch v := hd.(type) {
	case http.Handler:
	default:
		t.Errorf("%T type is not http.Handler", v)
	}
}

func TestSessionLoad(t *testing.T) {
	var (
		mokApp config.AppConfig
		th     testHandler
	)

	NewHandler(&mokApp, nil, nil)
	hd := Handle.SessionLoad(&th)

	switch v := hd.(type) {
	case http.Handler:
	default:
		t.Errorf("%T type is not http.Handler", v)
	}
}

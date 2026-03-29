package main

import (
	"testing"

	"github.com/emptywe/local_hotel/config"
)

func TestSetupRenderer(t *testing.T) {
	var mockApp config.AppConfig
	_, err := setupRenderer(&mockApp)
	if err != nil {
		t.Errorf("cannot setup renderer %v", err)
	}
}

func TestSetupSession(t *testing.T) {
	var mockApp config.AppConfig
	if setupSession(&mockApp); mockApp.Session == nil {
		t.Error("cannot setup session")
	}

}

func TestRun(t *testing.T) {
	var mockApp config.AppConfig
	_, err := setup(&mockApp, nil)
	if err != nil {
		t.Errorf("failed Run(): %v", err)
	}
}
